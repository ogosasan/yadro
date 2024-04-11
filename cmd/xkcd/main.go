package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"yadro/internal/comics"
	"yadro/internal/config"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "c", "default", "a string var")
	flag.Parse()
	var c config.Conf
	c.GetConf(confPath)
	var fileExist bool
	if _, err := os.Stat("database.json"); err == nil {
		fileExist = true
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})
	baseURL := c.Url + "/%d/info.0.json"
	numComics, err := comics.GetNumComics(c.Url + "/info.0.json")
	if err != nil {
		fmt.Println("Failed to get number of comics:", err)
		return
	}
	comicsMap := comics.GoToSite(numComics, baseURL, done, fileExist, c.Goroutines)
	select {
	case <-signalChan:
		fmt.Println("The signal is interrupted.")
	case <-done:
		fmt.Println("All comics fetched.")
	}
	comics.WriteFile(c.Bd, comicsMap)
}
