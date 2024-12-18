package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hongminhcbg/velocity-rule/config"
	"github.com/hongminhcbg/velocity-rule/src/router"
	"github.com/hongminhcbg/velocity-rule/src/service"
	"github.com/hongminhcbg/velocity-rule/src/store"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	cfg    *config.Config
	logger logr.Logger
)

func initDb() *gorm.DB {
	return nil
	db, err := gorm.Open(mysql.Open(cfg.MySqlUrl))
	if err != nil {
		panic(err)
	}

	dbx, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = dbx.Ping()
	if err != nil {
		panic("ping db error" + err.Error())
	}

	dbx.SetConnMaxIdleTime(5 * time.Minute)
	dbx.SetMaxIdleConns(25)
	dbx.SetMaxOpenConns(25)
	if cfg.Env != "prod" {
		db = db.Debug()
	}

	return db
}

func initRedis() *redis.Client {
	rdOtps, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(rdOtps)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic("ping redis error " + err.Error())
	}

	return redisClient
}

func main() {
	var err error
	cfg, err = config.Init()
	if err != nil {
		panic(err)
	}

	if cfg.Env != "prod" {
		b, _ := json.MarshalIndent(cfg, "", "\t")
		fmt.Println(string(b))
	}

	logger = cfg.InitLog()
	db := initDb()
	redisClient := initRedis()
	userStore := store.NewUseStore(db)
	svc := service.NewService(cfg, userStore, redisClient, logger)
	engine := gin.Default()
	router.InitGin(engine, svc)
	go func() {
		engine.Run()
	}()

	// handler graceful shutdown
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	for {
		select {
		case <-osSig:
			{
				fmt.Println("shutting down")
				time.Sleep(25 * time.Second)
				dbx, _ := db.DB()
				dbx.Close()
				redisClient.Close()
				os.Exit(0)
			}
		}
	}
}
