package main

import (
	"net/http"
	"os"

	"GoMailer/app"
	"GoMailer/log"
)

func main() {
	port := "8080"
	if s := os.Getenv("PORT"); s != "" {
		port = s
	}

	host := ""
	if app.IsDevAppServer() {
		host = "127.0.0.1"
	}

	addr := host + ":" + port
	log.Infof("server start at %s", addr)
	if err := http.ListenAndServe(addr, app.Handler()); err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
