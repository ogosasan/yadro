package main

import (
	"flag"
	"yadro/internal/comics"
	"yadro/internal/config"
)

func main() {
	var numComics int
	var screen bool
	flag.IntVar(&numComics, "n", 1234, "a int var")
	flag.BoolVar(&screen, "o", false, "a bool var")
	flag.Parse()
	var c config.Conf
	c.GetConf()
	baseURL := c.Url + "/%d/info.0.json"
	comicsMap := comics.GoToSite(numComics, baseURL)
	comics.WriteFile(screen, c.Bd, comicsMap)

}
