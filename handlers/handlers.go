package handlers

import (
	"BlockMe/config"
	"BlockMe/cron"
	"BlockMe/middleware"
	"BlockMe/to"
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"time"
)

func Home(w http.ResponseWriter, r *http.Request) {
	page, err := template.ParseFiles("templates/index.html")
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
	theirIpStr := middleware.IpFromContext(r.Context()).String()
	if theirIpStr == "" {
		fmt.Println("IP not found in context")
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	// Only block the address if it is not in the ignore list
	if !slices.Contains(config.Env.IgnoreList, theirIpStr) {
		sqlStatement := `INSERT INTO ip_blocklist (ip, timestamp) VALUES ($1, $2)`
		_, err := to.DB.Exec(sqlStatement, theirIpStr, time.Now().UTC())
		if err != nil {
			fmt.Printf("Error executing insert statement: %s\n", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		fmt.Printf("Ignoring block request for %s - IP is included in the IGNORE_LIST\n", theirIpStr)
	}

	page, err := template.ParseFiles("templates/blocked.html")
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

func ResetThem(w http.ResponseWriter, r *http.Request) {
	theirIpStr := middleware.IpFromContext(r.Context()).String()
	if theirIpStr == "" {
		fmt.Println("IP not found in context")
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	if !config.Env.ResetMe {
		fmt.Println("Reset endpoint is not enabled")
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	requestResetKey := r.URL.Query().Get("resetKey")
	fmt.Printf("requestResetKey: %s resetKey %s\n", requestResetKey, cron.ResetKey)
	if requestResetKey == *cron.ResetKey {
		fmt.Printf("Valid reset key, removing %s from the blockList\n", theirIpStr)

		sqlStatement := `DELETE FROM ip_blocklist WHERE ip = ?`
		_, err := to.DB.Exec(sqlStatement, theirIpStr)
		if err != nil {
			fmt.Printf("Error executing delete statement: %s\n", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
