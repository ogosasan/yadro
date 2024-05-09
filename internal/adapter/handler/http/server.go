package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"yadro/internal/adapter/repository"
	comics2 "yadro/internal/core/comics"
	"yadro/internal/core/config"
)

var db_path string

func Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Update...")
	var c config.Conf
	c.GetConf("configs/config.yaml")
	var fileExist bool
	if _, err := os.Stat("database.json"); err == nil {
		fileExist = true
	}
	resp := make(map[string]string)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	baseURL := c.Url + "/%d/info.0.json"
	numComics := comics2.GetNumComics(baseURL)
	comicsMap, indexMap, count := comics2.GoToSite(numComics, baseURL, signalChan, fileExist, c.Goroutines)
	<-signalChan
	db_path = c.Dsn
	repository.Head(db_path, comicsMap, indexMap)
	resp["total comics"] = strconv.Itoa(numComics)
	resp["new comics"] = strconv.Itoa(numComics - count)
	//comics2.WriteFile(c.Bd, comicsMap, indexMap)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp["message"] = "Status OK"
	jsonResp, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	log.Println("The update is finished.")
	return
}

func Pics(w http.ResponseWriter, r *http.Request) {
	log.Println("Search...")
	comicsMap, indexMap := repository.FetchRecords(db_path)
	line := r.URL.Query().Get("search")
	ans := comics2.IndexSearch(indexMap, comicsMap, line)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status Created"
	jsonResp, err := json.MarshalIndent(ans, "", "\t")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func Delete(w http.ResponseWriter, r *http.Request) {
	repository.Down(db_path)
	log.Println("Tables have been deleted.")
}
