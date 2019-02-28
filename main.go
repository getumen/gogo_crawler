package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/getumen/gogo_crawler/config"
	"github.com/getumen/gogo_crawler/di"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	configPath := flag.String("conf", "example_config.toml", "config file path")
	flag.Parse()
	file, err := os.Open(*configPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	var conf config.Config
	if _, err := toml.Decode(string(configBytes), &conf); err != nil {
		log.Fatalln(err)
	}

	db := NewMySQLConnection(&conf)
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	redisConn := NewRedisConnection(&conf)
	defer func() {
		err := redisConn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	crawler, err := di.InitializeCrawler(&conf, db, redisConn)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Initialized")

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	var signalChan = make(chan os.Signal, 1)
	go func() {
		cnt := 0
		signal.Notify(signalChan,
			syscall.SIGINT)
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGINT:
				log.Println("SIGINT")
				cnt += 1
				if cnt > 1 {
					os.Exit(0)
				} else {
					cancel()
				}
			default:
				log.Fatalln("Unknown signal.")
			}
		}
	}()

	crawler.Start(ctx, &conf)
}

func NewMySQLConnection(config *config.Config) *gorm.DB {
	c := config.MySQL
	password := os.Getenv(c.PasswordEnvKey)
	connStr := fmt.Sprintf(
		"%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		password,
		c.Connection,
		c.Host,
		c.Port,
		c.DatabaseName,
	)
	conn, err := gorm.Open("mysql", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func NewRedisConnection(conf *config.Config) *redis.Pool {

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				conf.Redis.Connection,
				fmt.Sprintf(
					"%s:%d",
					conf.Redis.Host,
					conf.Redis.Port),
			)
		},
	}
}
