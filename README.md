# Upstash Redis Go

An HTTP/REST based Redis client built on top of [Upstash REST API](https://docs.upstash.com/features/restapi).

Inspired by [The official typescript client](https://github.com/upstash/upstash-redis)

See [the list of APIs](https://docs.upstash.com/features/restapi#rest---redis-api-compatibility) supported.

[![codecov](https://codecov.io/gh/chronark/upstash-go/branch/main/graph/badge.svg?token=BCNI6L3TRT)](https://codecov.io/gh/chronark/upstash-go)
## Quick Start


```go
package main

import (
	"fmt"
	"github.com/chronark/upstash-go"
)

func main() {
	u, err := upstash.New(upstash.Options{
		// Get your url and token from https://console.upstash.com/redis/<uuid>
        // Or leave empty to load from environment variables
		Url: "", // env: UPSTASH_REDIS_REST_URL
		Token:    "", // env: UPSTASH_REDIS_REST_TOKEN
	})
	if err != nil {
		panic(err)
	}

	key := "foo"

	err = u.Set(key, "bar")
	if err != nil {
		panic(err)
	}

	value, err := u.Get(key)
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
	// -> "bar"

}

```
