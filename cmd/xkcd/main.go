package main

import (
	"log"
	"net/http"
	http2 "yadro/internal/adapter/handler/http"
	"yadro/internal/core/config"
	"yadro/internal/core/update"
)

func main() {
	var c config.Conf
	c.GetConf("configs/config.yaml")
	mux := http.NewServeMux()
	mux.HandleFunc("/update", http2.Update)
	mux.HandleFunc("/pics", http2.Pics)
	mux.HandleFunc("/login", http2.LoginHandler)
	mux.HandleFunc("/comics", http2.Comics)
	mux.HandleFunc("/xkcd-server/login", http2.Login)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	go func() {
		log.Printf("Starting the server on http://%s", c.Port)
		err := http.ListenAndServe(c.Port, mux)
		log.Fatal(err)
	}()
	update.UpdateEveryDay(c.Port)
}
