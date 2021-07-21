package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var servers []string
var ctx = context.Background()
var rdb *redis.Client

func main() {
	setupLogger()
	setupRedis()
	configServers()
	startWorker()

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/data", handleData)

	err := http.ListenAndServe(":8585", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setupLogger() {
	n := fmt.Sprintf("./logs/%s.log", time.Now().Format("2006-01-02"))
	f, err := os.OpenFile(n, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)
	log.Println("START")
}

func setupRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
}

func configServers() {
	servers = strings.Split(os.Getenv("TARGETS"), ",")
	if len(servers) > 5 {
		log.Fatalln("Cannot handle more than 5 servers.")
	}

	log.Println(servers)
}

func handleData(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-type", "application/json")

	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		publishError(w, err)
	}

	values := make(map[string]map[string]int64)
	for _, k := range keys {
		kp := strings.Split(k, "=")

		cmd := rdb.Get(ctx, k)
		if cmd.Err() != nil {
			publishError(w, cmd.Err())
		}

		v, err := strconv.Atoi(cmd.Val())
		if err != nil {
			publishError(w, err)
		}

		if values[kp[0]] == nil {
			values[kp[0]] = make(map[string]int64)
		}

		values[kp[0]][kp[1]] = int64(v)
	}

	j, err := json.Marshal(values)
	if err != nil {
		publishError(w, err)
	}

	if _, err := fmt.Fprintf(w, string(j)); err != nil {
		log.Fatal(err)
	}
}

func publishError(w http.ResponseWriter, e error) {
	w.WriteHeader(500)

	if _, err := fmt.Fprint(w, "{\"error\":\"Internal Error\"}"); err != nil {
		log.Println(e)
		log.Fatal(err)
	}

	log.Fatal(e)
}

func startWorker() {
	go func() {
		for range time.Tick(15 * time.Second) {
			for _, s := range servers {
				go func(s string) {
					err := call(s)
					if err != nil {
						log.Println(err)
					}
				}(s)
			}
		}
	}()
}

func call(server string) error {
	pinger, err := ping.NewPinger(server)
	if err != nil {
		return err
	}

	pinger.SetPrivileged(true)
	pinger.Count = 3
	pinger.Timeout = 10 * time.Second
	if err = pinger.Run(); err != nil {
		return err
	}

	stats := pinger.Statistics()
	if err = persist(server, int64(stats.AvgRtt)); err != nil {
		return err
	}

	return nil
}

func persist(server string, rtt int64) error {
	key := fmt.Sprintf("%s=%s", server, time.Now().Format("2006/01/02 15:04"))
	if err := rdb.Set(ctx, key, rtt, 24*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}
