package main

import (
	"fmt"
	"math/rand"
	"os"
	"tasksvc/config"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // it's important for random string generation
	var cfg config.StartupConfig
	cfg.ReadConfig()
	if err := os.Setenv("IN", "Asia/Kolkata"); err != nil {
		fmt.Println("Failed to set the timezone")
		os.Exit(1)
	}

	app, err := InitDependency(&cfg)
	if err != nil {
		panic(err)
	}

	err = app.Start()
	if err != nil {
		panic(err)
	}

	fmt.Println("Server started on:", time.Now().Format(time.RFC1123))
}
