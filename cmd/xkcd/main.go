package main

import (
	"log"
	"net/http"
	"time"
	http2 "yadro/internal/adapter/handler/http"
	"yadro/internal/core/config"
)

func main() {
	var c config.Conf
	c.GetConf("configs/config.yaml")
	mux := http.NewServeMux()
	mux.HandleFunc("/update", http2.Update)
	mux.HandleFunc("/pics", http2.Pics)
	go func() {
		log.Printf("Starting the server on http://127.0.0.1%s", c.Port)
		err := http.ListenAndServe(c.Port, mux)
		log.Fatal(err)
	}()

	timer := time.NewTicker(24 * time.Hour)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			resp, err := http.Get("http://127.0.0.1" + c.Port + "/update")
			if err != nil {
				log.Fatalf("Request execution error: %v", err)
			}
			defer resp.Body.Close()

			log.Println("The request was completed successfully")
		}
	}

}
