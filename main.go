package main

import (
	"BlockMe/config"
	"BlockMe/cron"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var blockList []string
var resetKey string

func main() {
	fmt.Println("Setting up...")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	config.InitConfig()

	router := mux.NewRouter()

	resetKeyPtr := &resetKey
	cronConfig := cron.JobConfig{ResetKey: resetKeyPtr}
	c := cron.InitCronJobs(&cronConfig)

	router.Use(LogRequestMiddleware, IpMiddleware, BlockCheckMiddleware)

	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/blockme", BlockThem).Methods("POST")

	if config.Env.IAmALittleBitch {
		router.HandleFunc("/reset", ResetThem).Methods("GET")
	}

	port := 8080
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		fmt.Printf("\nStarting server on port %d\n", port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("The web server failed to start... %s\n", err.Error())
		}
	}()

	<-sigs
	fmt.Printf("\nReceived termination signal. Starting graceful shutdown.\n")

	fmt.Println("Stopping the cron jobs...")
	c.Stop()

	fmt.Println("Stopping web server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("The web server failed to shut down gracefully... %s\n", err.Error())
	}

	fmt.Println("\nServer gracefully stopped")
}

func Home(w http.ResponseWriter, r *http.Request) {
	page, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Printf("template parsing failed: %s\n", err.Error())
		http.Error(w, "something didn't work right", http.StatusInternalServerError)
		return
	}

	err = page.Execute(w, nil)
	if err != nil {
		fmt.Printf("template executing failed: %s\n", err.Error())
		http.Error(w, "something didn't work right", http.StatusInternalServerError)
		return
	}
}

func BlockThem(w http.ResponseWriter, r *http.Request) {
	blockList = append(blockList, r.RemoteAddr)
	w.WriteHeader(http.StatusCreated)
	// idc if they see the blocked message or not - Ignore errors
	_, _ = fmt.Fprintf(w, "Blocked!")
	fmt.Printf("Blocked %s\n", r.RemoteAddr)
	return
}

func ResetThem(w http.ResponseWriter, r *http.Request) {
	requestResetKey := r.URL.Query().Get("resetKey")
	fmt.Printf("requestResetKey: %s resetKey %s\n", requestResetKey, resetKey)
	if requestResetKey == resetKey {
		fmt.Printf("Valid reset key, removing %s from the blockList\n", r.RemoteAddr)
		resetIpIndex := slices.Index(blockList, r.RemoteAddr)

		if resetIpIndex != -1 {
			blockList = slices.Delete(blockList, resetIpIndex, resetIpIndex+1)
		}
	}
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func IpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse X-Forwarded-For header to net.IP
		theirIPStr := r.Header.Get("X-Forwarded-For")
		if theirIPStr == "" {
			theirIPStr = r.RemoteAddr
			errMsg := fmt.Sprintf("X-Forwarded-For header is empty, going with IP %s", r.RemoteAddr)
			fmt.Printf("%s\n", errMsg)
			next.ServeHTTP(w, r)
			return
		}

		// Verify XFF is an actual IP using net.ParseIP
		theirIP := net.ParseIP(theirIPStr)
		if theirIP == nil {
			http.Error(w, "You didn't pass a valid IP in the X-Forwarded-For header. Stop doing that.", http.StatusBadRequest)
			return
		}

		r.RemoteAddr = theirIP.String()
		next.ServeHTTP(w, r)
	})
}

func BlockCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "reset") {
			fmt.Printf("reset link, skipping block check\n")
		} else if slices.Contains(blockList, r.RemoteAddr) {
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte("You are blocked!"))
			if err != nil {
				fmt.Printf("main.BlockCheckMiddleware: could not write to response body: %s\n", err.Error())
			}

			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
