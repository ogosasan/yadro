package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"yadro/internal/comics"
	"yadro/internal/config"
)

func main() {
	var confPath, str string
	flag.StringVar(&confPath, "c", "default", "a string var")
	flag.StringVar(&str, "s", "default", "a string var")
	index := flag.Bool("i", false, "a bool flag")
	flag.Parse()
	var c config.Conf
	c.GetConf(confPath)
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
	}
}
