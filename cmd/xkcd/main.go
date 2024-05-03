package main

import (
	"log"
	"net/http"
	"time"
	"yadro/internal/config"
)

func main() {
	/*var str string
	flag.StringVar(&str, "s", "default", "a string var")
	index := flag.Bool("i", false, "a bool flag")
	flag.Parse()
	var c config.Conf
	c.GetConf("configs/config.yaml")
	var fileExist bool
	if _, err := os.Stat("database.json"); err == nil {
		fileExist = true
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	baseURL := c.Url + "/%d/info.0.json"
	numComics := comics.GetNumComics(baseURL)
	comicsMap, indexMap := comics.GoToSite(numComics, baseURL, signalChan, fileExist, c.Goroutines)
	<-signalChan
	comics.WriteFile(c.Bd, comicsMap, indexMap)
	if *index {
		comics.IndexSearch(indexMap, comicsMap, str)
	} else {
		comics.Search(comicsMap, str)
	}*/
	var c config.Conf
	c.GetConf("configs/config.yaml")
	mux := http.NewServeMux()
	mux.HandleFunc("/update", home)
	mux.HandleFunc("/pics", showSnippet)
	//mux.HandleFunc("/update", createSnippet)
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
