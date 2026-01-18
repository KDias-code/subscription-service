package main

import (
	"log"
	"subscription-service/configs"
	_ "subscription-service/docs"
	"subscription-service/internal/app"
)

const (
	cgfPath = "manifests/config.yaml"
)

// @title Subscriptions Service API
// @version 1.0
// @description API для управления подписками
// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	cfg, err := configs.Load(cgfPath)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Start(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	return
}
