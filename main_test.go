package main

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	comics2 "yadro/internal/core/comics"
	"yadro/internal/core/config"
)

var (
	c         config.Conf
	comicsMap map[int]comics2.Write
	indexMap  map[string][]int
	str       string
)

func setup() {
	c.GetConf("configs/config.yaml")

	var fileExist bool
	if _, err := os.Stat("database.json"); err == nil {
		fileExist = true
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	baseURL := c.Url + "/%d/info.0.json"
	numComics := comics2.GetNumComics(baseURL)
	comicsMap, indexMap, _ = comics2.GoToSite(numComics, baseURL, signalChan, fileExist, c.Goroutines)
	<-signalChan
}

func BenchmarkIndexSearch(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comics2.IndexSearch(indexMap, comicsMap, str)
	}
}

func BenchmarkSearch(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comics2.Search(comicsMap, str)
	}
}
