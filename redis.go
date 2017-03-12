package main

import "os"

func getRedisURL() string {
	return os.Getenv("REDIS_URL") // i.e.) redis://user:secret@localhost:6379/0?foo=bar&qux=baz
}
