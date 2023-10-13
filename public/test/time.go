package main

import (
	"log"
	"time"
)

func main() {
	parse, err := time.Parse("2006-01-02", "2023-10-13")
	if err != nil {
		log.Println(err)
	}
	log.Println(parse.Unix())
}
