package app

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	server "github.com/tumbleweedd/shortener"
	"github.com/tumbleweedd/shortener/internal/handlers"
	"github.com/tumbleweedd/shortener/internal/repositories"
	"github.com/tumbleweedd/shortener/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	repo := repositories.NewRepository(client)
	service := services.NewService(repo)
	handler := handlers.NewHandler(service)

	srv := new(server.Server)

	go func() {
		if err := srv.Run("8000", handler.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Println("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("App Shutting Down")

	if err := srv.ShutDown(context.Background()); err != nil {
		logrus.Errorf("error ocured on server shutting down: %s", err.Error())
	}
}
