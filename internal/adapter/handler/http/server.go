package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"yadro/internal/comics"
	"yadro/internal/config"
)

var comicsMap = map[int]comics.Write{}
var indexMap = map[string][]int{}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/update" {
		http.NotFound(w, r)
		return
	}
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
	numComics := comics.GetNumComics(baseURL)
	var count int
	comicsMap, indexMap, count = comics.GoToSite(numComics, baseURL, signalChan, fileExist, c.Goroutines)
	<-signalChan
	resp["total comics"] = strconv.Itoa(numComics)
	resp["new comics"] = strconv.Itoa(numComics - count)
	comics.WriteFile(c.Bd, comicsMap, indexMap)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp["message"] = "Status OK"
	jsonResp, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func Pics(w http.ResponseWriter, r *http.Request) {
	line := r.URL.Query().Get("search")
	ans := comics.IndexSearch(indexMap, comicsMap, line)
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
