package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/korovkin/limiter"
	_ "github.com/korovkin/limiter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
			fmt.Println("You are not an admin")
			return
		}
	}
	log.Println("Update...")
	var c config.Conf
	c.GetConf("configs/config.yaml")
	var fileExist bool
	if _, err := os.Stat("database.json"); err == nil {
		fileExist = true
	}
	resp := make(map[string]string)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	baseURL := c.Url + "/%d/info.0.json"
	numComics := comics2.GetNumComics(baseURL)
	comicsMap, indexMap, count := comics2.GoToSite(numComics, baseURL, signalChan, fileExist, c.Goroutines)
	db_path = c.Dsn
	climit = c.CLimit
	rlimit = c.RLimit
	tokentime = c.TokenTime
	repository.Head(db_path, comicsMap, indexMap)
	resp["total comics"] = strconv.Itoa(numComics)
	resp["new comics"] = strconv.Itoa(numComics - count)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp["message"] = "Status OK"
	jsonResp, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
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
	fmt.Println(rlimit)
	d, _ := r.Cookie("token")
	claims := &Claims{}
	tknStr := d.Value
	jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	ctx := context.Background()
	if claims.Username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !comics2.TokenBucketRateLimit(ctx, redisClient, claims.Username, refillWindow, rlimit) {
		http.Error(w, "Превышен лимит запросов", http.StatusTooManyRequests)
		return
	}
	limit.Execute(func() {
		log.Println("Search...")
		comicsMap, indexMap := repository.FetchRecords(db_path)
		line := r.URL.Query().Get("search")
		ans := comics2.IndexSearch(indexMap, comicsMap, line)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Status Created"
		jsonResp, err := json.MarshalIndent(ans, "", "\t")
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	repository.Down(db_path)
	log.Println("Tables have been deleted.")
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
	fmt.Println(tokentime)
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, _ := ioutil.ReadFile("users.json")
	var users []Credentials
	er := json.Unmarshal(file, &users)
	if er != nil {
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
