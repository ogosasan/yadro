package update

import (
	"log"
	"net/http"
	"time"
)

func UpdateEveryDay(port string) {
	time.Sleep(1 * time.Second)
	url := "http://" + port + "/update?auto='true'"
	up, err := http.Get(url)
	if err != nil {
		log.Fatalf("Request execution error: %v", err)
	}
	defer up.Body.Close()
	timer := time.NewTicker(24 * time.Hour)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			resp, err := http.Get(url)
			if err != nil {
				log.Fatalf("Request execution error: %v", err)
			}
			func() {
				defer resp.Body.Close()
				log.Println("The request was completed successfully")
				// Дополнительная обработка ответа, если необходимо
			}()
		}
	}
}
