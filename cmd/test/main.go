package main

import (
	"fmt"
	"github.com/chronark/upstash-go"
)

func main() {
	// Get your url and token from https://console.upstash.com/redis/{id}
	// Or leave empty to load from environment variables
	options := upstash.Options{
		Url:   "", // env: UPSTASH_REDIS_REST_URL
		Token: "", // env: UPSTASH_REDIS_REST_TOKEN
	}

	u, err := upstash.New(options)
	if err != nil {
		panic(err)
	}
	err = u.Set("foo", "bar")
	if err != nil {
		panic(err)
	}
	value, err := u.Get("foo")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	// -> "bar"
}
