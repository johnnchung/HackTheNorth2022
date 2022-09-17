package main

import (
	"log"
	// add this

	handler "github.com/johnnchung/HackTheNorth2022/handler"
	"github.com/joho/godotenv"
)

func main() {
	repo := &handler.Repo{}

	godotenv.Load()
	if err := repo.HandlerInit(); err != nil {
		log.Fatal(err)
		return
	}
	if err := repo.Run(); err != nil {
		log.Fatal(err)
	}
}
