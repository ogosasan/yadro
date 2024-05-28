package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/korovkin/limiter"
	_ "github.com/korovkin/limiter"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"yadro/internal/adapter/repository"
	comics2 "yadro/internal/core/comics"
	"yadro/internal/core/config"
)

var db_path string
var climit int
var rlimit int64
var tokentime int

func Update(w http.ResponseWriter, r *http.Request) {
	auto := r.URL.Query().Get("auto")
	if auto == "" {
		d, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := d.Value

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims.Username != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("You are not an admin")
			return
		}
	}
	log.Println("Update...")
	var c config.Conf
	c.GetConf("configs/config.yaml")
	/*var fileExist bool
	if _, err := os.Stat("database.json"); err == nil {
		fileExist = true
	}*/
	resp := make(map[string]string)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	baseURL := c.Url + "/%d/info.0.json"
	numComics := comics2.GetNumComics(baseURL)
	//comicsMap, indexMap, count := comics2.GoToSite(numComics, baseURL, signalChan, fileExist, c.Goroutines)
	db_path = c.Dsn
	climit = c.CLimit
	rlimit = c.RLimit
	tokentime = c.TokenTime
	//repository.Head(db_path, comicsMap, indexMap)
	resp["total comics"] = strconv.Itoa(numComics)
	//resp["new comics"] = strconv.Itoa(numComics - count)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp["message"] = "Status OK"
	jsonResp, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
	log.Println("The update is finished.")
	return
}

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	refillWindow int64 = 60
)

var limit = limiter.NewConcurrencyLimiter(climit)

func Pics(w http.ResponseWriter, r *http.Request) {
	line := r.URL.Query().Get("search")
	if !strings.HasPrefix(line, "tst") {
		d, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		tknStr := d.Value
		_, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil || claims.Username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.Background()
		if !comics2.TokenBucketRateLimit(ctx, redisClient, claims.Username, refillWindow, rlimit) {
			http.Error(w, "Превышен лимит запросов", http.StatusTooManyRequests)
			return
		}
	}
	log.Println("Search...")
	w.WriteHeader(http.StatusOK)
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	comicsMap, indexMap := repository.FetchRecords(db)
	ans := comics2.IndexSearch(indexMap, comicsMap, line)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.MarshalIndent(ans, "", "\t")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		return
	}
	return
}

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, _ := os.Open("users.json")
	content, _ := io.ReadAll(file)
	var users []Credentials
	er := json.Unmarshal(content, &users)
	if er != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userMap := make(map[string]string)
	for _, user := range users {
		userMap[user.Username] = user.Password
	}

	expectedPassword, ok := userMap[creds.Username]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Duration(tokentime) * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
