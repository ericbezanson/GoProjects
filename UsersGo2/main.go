package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal((err))
	}

	if err := store.Init(); err != nil {
		log.Fatal((err))
	}
	fmt.Println("%+v\n", store)

	server := newAPIServer(":3000", store)
	fmt.Println("hello world")
	server.Run()
}
