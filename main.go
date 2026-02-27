package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var blockList []string

func main() {
	fmt.Println("Setting up...")
	router := mux.NewRouter()
	err := godotenv.Load()

	router.Use(LogRequestMiddleware, IpMiddleware, BlockCheckMiddleware)

	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/blockme", BlockThem).Methods("POST")

	IAmALittleBitchStr := os.Getenv("I_AM_A_LITTLE_BITCH")
	IAmALittleBitch, err := strconv.ParseBool(IAmALittleBitchStr)
	if err != nil {
		IAmALittleBitch = false
	}

	if IAmALittleBitch {
		router.HandleFunc("/reset", ResetThem).Methods("GET")
	}

	port := 8080
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Starting server on port %d\n", port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("The server didn't start... %s\n", err.Error())
	}
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
	resetIpIndex := slices.Index(blockList, r.RemoteAddr)

	if resetIpIndex != -1 {
		blockList = slices.Delete(blockList, resetIpIndex, resetIpIndex+1)
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
		if strings.Contains(r.URL.Path, "reset") {
			next.ServeHTTP(w, r)
		}
		if slices.Contains(blockList, r.RemoteAddr) {
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
