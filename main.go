package main

import (
	"net/http"

	"GoMailer/app"
	"GoMailer/conf"
	"GoMailer/handler"
	_ "GoMailer/handler/mail"
	_ "GoMailer/handler/shortcut"
	"GoMailer/log"
)

func init() {
	http.Handle("/api/", handler.Router)
}

func main() {
	port := "8080"
	if s := conf.App().Port; s != "" {
		port = s
	}

	host := ""
	if app.IsDevAppServer() {
		host = "127.0.0.1"
	}

	addr := host + ":" + port
	log.Infof("server start at %s with env: %s", addr, conf.Env())
	if err := http.ListenAndServe(addr, app.RootHandler()); err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
