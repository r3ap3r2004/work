package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gocraft/work/webui"
	"github.com/gomodule/redigo/redis"
)

var (
	redisPassword  = flag.String("password", "", "redis password")
	redisHostPort  = flag.String("redis", ":6379", "redis hostport")
	redisDatabase  = flag.String("database", "0", "redis database")
	redisNamespace = flag.String("ns", "work", "redis namespace")
	webHostPort    = flag.String("listen", ":5040", "hostport to listen for HTTP JSON API")
)

func main() {
	flag.Parse()

	fmt.Println("Starting workwebui:")
	fmt.Println("redis = ", *redisHostPort)
	fmt.Println("database = ", *redisDatabase)
	fmt.Println("namespace = ", *redisNamespace)
	fmt.Println("listen = ", *webHostPort)
	if len(*redisPassword) > 0 {
		fmt.Println("accessing with password")
	}

	database, err := strconv.Atoi(*redisDatabase)
	if err != nil {
		fmt.Printf("Error: %v is not a valid database value", *redisDatabase)
		return
	}

	pool := newPool(*redisHostPort, database, redisPassword)

	server := webui.NewServer(*redisNamespace, pool, *webHostPort)
	server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-c

	server.Stop()

	fmt.Println("\nQuitting...")
}

func newPool(addr string, database int, password *string) *redis.Pool {
	return &redis.Pool{
		MaxActive:   3,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr, redis.DialDatabase(database), redis.DialPassword(*password), redis.DialUseTLS(true), redis.DialTLSSkipVerify(true))
		},
		Wait: true,
	}
}
