package comics

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"testing"
)

func TestGetNumComics(t *testing.T) {
	url := "https://xkcd.com/%d/info.0.json"
	response, _ := http.Get("https://xkcd.com/info.0.json")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)
	var info InfoResponse
	err := json.NewDecoder(response.Body).Decode(&info)
	if err != nil {
		return
	}
	expected := info.Num
	result := GetNumComics(url)
	if result != expected {
		t.Errorf("Incorrect result. Expect %d, got %d", expected, result)
	}
}

func TestGoToSite(t *testing.T) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	comicsMap, indexMap, count := GoToSite(1, "https://xkcd.com/%d/info.0.json", signalChan, false, 1)
	comicsMap_expected := make(map[int]Write)
	// Присвоение значения
	comicsMap_expected[1] = Write{
		Tscript: []string{"boy", "sit", "barrel", "float", "ocean", "wonder", "next", "drift", "distanc", "noth", "els", "seen", "alt", "part", "1"},
		Img:     "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg",
	}
	indexMap_expected := make(map[string][]string)
	indexMap_expected["alt"] = []string{"1"}
	indexMap_expected["barrel"] = []string{"1"}
	indexMap_expected["boy"] = []string{"1"}
	indexMap_expected["distanc"] = []string{"1"}
	indexMap_expected["drift"] = []string{"1"}
	indexMap_expected["els"] = []string{"1"}
	indexMap_expected["float"] = []string{"1"}
	indexMap_expected["next"] = []string{"1"}
	indexMap_expected["noth"] = []string{"1"}
	indexMap_expected["ocean"] = []string{"1"}
	indexMap_expected["seen"] = []string{"1"}
	indexMap_expected["sit"] = []string{"1"}
	indexMap_expected["wonder"] = []string{"1"}
	indexMap_expected["part"] = []string{"1"}
	indexMap_expected["1"] = []string{"1"}
	if !reflect.DeepEqual(comicsMap_expected, comicsMap) || !reflect.DeepEqual(indexMap_expected, indexMap) || count != 0 {
		t.Errorf("Incorrect result.")
	}
}
