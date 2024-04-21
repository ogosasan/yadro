package нфвкщ

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"yadro/internal/comics"
	"yadro/internal/config"
)

var (
	c         config.Conf
	comicsMap map[int]comics.Write
	indexMap  map[string][]int
	str       string
)

func setup() {
	flag.StringVar(&str, "s", "default", "a string var")
	flag.Parse()
	c.GetConf("configs/config.yaml")

	var fileExist bool
	if _, err := os.Stat("database.json"); err == nil {
		fileExist = true
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	baseURL := c.Url + "/%d/info.0.json"
	numComics := comics.GetNumComics(baseURL)
	comicsMap, indexMap = comics.GoToSite(numComics, baseURL, signalChan, fileExist, c.Goroutines)
	<-signalChan
}

func BenchmarkIndexSearch(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comics.IndexSearch(indexMap, comicsMap, str)
	}
}

func BenchmarkSearch(b *testing.B) {
	setup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		comics.Search(comicsMap, str)
	}
}
