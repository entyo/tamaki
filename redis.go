package main

import (
	"log"
	"os"
)

func getRedisURL() string {
	if os.Getenv("REDIS_URL") == "" {
		log.Fatalln("REDIS_URL is not set!")
	}
	return os.Getenv("REDIS_URL") // i.e.) redis://user:secret@localhost:6379/0?foo=bar&qux=baz
}
