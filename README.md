# Upstash Redis Go

An HTTP/REST based Redis client built on top of [Upstash REST API](https://docs.upstash.com/features/restapi).

Inspired by [The official typescript client](https://github.com/upstash/upstash-redis)

See [the list of APIs](https://docs.upstash.com/features/restapi#rest---redis-api-compatibility) supported.

[![codecov](https://codecov.io/gh/chronark/upstash-go/branch/main/graph/badge.svg?token=BCNI6L3TRT)](https://codecov.io/gh/chronark/upstash-go)

## Quick Start

Error handling has been omitted for better readability.

```go
package main

import (
	"fmt"
	"github.com/chronark/upstash-go"
)

func main() {


    // Get your url and token from https://console.upstash.com/redis/{id}
    // Or leave empty to load from environment variables
    options := upstash.Options{
        Url: "", // env: UPSTASH_REDIS_REST_URL
        Token:    "", // env: UPSTASH_REDIS_REST_TOKEN
    }

    u, _ := upstash.New(options)

    u.Set("foo", "bar")

    value, _ := u.Get("foo")

    fmt.Println(value)
    // -> "bar"

}

```
