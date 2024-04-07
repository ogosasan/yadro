package main

import (
	"flag"
	"fmt"
	"yadro/internal/comics"
	"yadro/internal/config"
)

func main() {
	var num int
	var screen bool
	flag.IntVar(&num, "n", -1, "a int var")
	flag.BoolVar(&screen, "o", false, "a bool var")
	flag.Parse()
	var c config.Conf
	c.GetConf()
	baseURL := c.Url + "/%d/info.0.json"
	numComics, err := comics.GetNumComics(c.Url + "/info.0.json")
	if err != nil {
		fmt.Println("Failed to get number of comics:", err)
		return
	}
	comicsMap := comics.GoToSite(numComics, baseURL)
	comics.WriteFile(screen, num, c.Bd, comicsMap)
}
