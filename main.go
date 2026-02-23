package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"slices"
	"time"

	"github.com/gorilla/mux"
)

var blockList []string

func main() {
	fmt.Println("Setting up...")
	router := mux.NewRouter()

	router.Use(LogRequestMiddleware, IpMiddleware, BlockCheckMiddleware)

	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/blockme", BlockThem).Methods("POST")

	port := 8080
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Starting server on port %d\n", port)
	err := server.ListenAndServe()
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

func IpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse X-Forwarded-For header to net.IP
		theirIPStr := r.Header.Get("X-Forwarded-For")
		if theirIPStr == "" {
			http.Error(w, "Your IP wasn't included in your request. Figure it out", http.StatusBadRequest)
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
		if slices.Contains(blockList, r.RemoteAddr) {
			w.WriteHeader(http.StatusForbidden)
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
