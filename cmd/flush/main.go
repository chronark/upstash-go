package main

import (
	"github.com/chronark/upstash-go"
)

func main() {
	u, err := upstash.New(upstash.Options{})
	if err != nil {
		panic(err)
	}

	err = u.FlushAll()
	if err != nil {
		panic(err)
	}
}
