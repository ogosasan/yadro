package comics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type GistRequest struct {
	Tscript string `json:"transcript"`
	Img     string `json:"img"`
}

func GoToSite(numComics int, baseURL string) map[int]GistRequest {
	comicsMap := make(map[int]GistRequest)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 1; i <= numComics; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			url := fmt.Sprintf(baseURL, i)

			response, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			if response.StatusCode != http.StatusOK {
				fmt.Printf("%s not found", url)
				return
			}

			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			var xkcd GistRequest
			err = json.Unmarshal(body, &xkcd)
			if err != nil {
				fmt.Println(err)
				return
			}

			xkcd.Tscript = normalization(xkcd.Tscript)

			mu.Lock()
			comicsMap[i] = xkcd
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return comicsMap
}

func WriteFile(screen bool, file string, comicsMap map[int]GistRequest) {
	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	if screen {
		bytes, _ := json.MarshalIndent(comicsMap, "", "\t")
		fmt.Println(string(bytes))
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(comicsMap)
	if err != nil {
		fmt.Println(err)
		return
	}
}
