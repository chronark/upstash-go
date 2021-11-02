package upstash

import (
	"fmt"
	"os"
	"strconv"

	"github.com/chronark/upstash-go/client"
)

type Upstash struct {
	client client.Client
}

type Options struct {

	// The Upstash endpoint you want to use
	// Defaults to `UPSTASH_REDIS_EDGE_URL` and falls back to `UPSTASH_REDIS_REST_URL`
	// environment variables.
	Endpoint string

	// Requests to the Upstash API must provide an API token.
	Token string
}

func New(options Options) (Upstash, error) {
	if options.Endpoint == "" {
		options.Endpoint = os.Getenv("UPSTASH_REDIS_REST_URL")
	}
	if options.Token == "" {
		options.Token = os.Getenv("UPSTASH_REDIS_REST_TOKEN")
	}

	return Upstash{
		client: client.New(options.Endpoint, options.Token, nil),
	}, nil
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

// If key already exists and is a string, this command appends the value at
// the end of the string. If key does not exist it is created and set as an
// empty string, so APPEND will be similar to SET in this special case.
//
// Return the length of the string after the append operation.
//
// https://redis.io/commands/append
func (u *Upstash) Append(key string, value string) (int, error) {

	res, err := u.client.Call(client.Request{
		Body: []string{"append", key, value},
	})
	return int(res.(float64)), err
}

// Decrements the number stored at key by one. If the key does not exist, it is
// set to 0 before performing the operation. An error is returned if the key
// contains a value of the wrong type or contains a string that can not be
// represented as integer. This operation is limited to 64 bit signed integers.
//
// See INCR for extra information on increment/decrement operations.
//
// Returns  the value of key after the decrement
//
// https://redis.io/commands/decr
func (u *Upstash) Decr(key string) (int, error) {

	res, err := u.client.Call(client.Request{
		Body: []string{"decr", key},
	})
	return int(res.(float64)), err
}

// Decrements the number stored at key by decrement. If the key does not
// exist, it is set to 0 before performing the operation. An error is
// returned if the key contains a value of the wrong type or contains a
// string that can not be represented as integer. This operation is limited
// to 64 bit signed integers.
//
// See INCR for extra information on increment/decrement operations.
//
// Returns the value of key after the decrement
//
// https://redis.io/commands/decrby
func (u *Upstash) DecrBy(key string, decrement int) (int, error) {

	res, err := u.client.Call(client.Request{
		Body: []string{"decrby", key, fmt.Sprintf("%d", decrement)},
	})
	return int(res.(float64)), err
}

// Get the value of key. If the key does not exist the special value nil is
// returned. An error is returned if the value stored at key is not a
// string, because GET only handles string values.
//
// Returns the value of key, or empty string when key does not exist.
//
// https://redis.io/commands/get
func (u *Upstash) Get(key string) (string, error) {

	res, err := u.client.Call(client.Request{
		Body: []string{"get", key},
	})
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.(string), nil
}

// Get the value of key and delete the key. This command is similar to GET,
// except for the fact that it also deletes the key on success (if and only
// if the key's value type is a string).
//
//Returns the value of key, empty string when key does not exist, or an
// error if the key's value type isn't a string.
//
// https://redis.io/commands/getdel
func (u *Upstash) GetDel(key string) (string, error) {
	res, err := u.client.Call(client.Request{
		Body: []string{"getdel", key},
	})
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.(string), nil
}

// Get the value of key and optionally set its expiration. GETEX is similar
//  to GET, but is a write command with additional options.
//
// Returns the value of `key`, or empty string when `key` does not exist
//
// https://redis.io/commands/getex
func (u *Upstash) GetEX(key string, options GetEXOptions) (string, error) {
	body := []string{"getdel", key}
	if options.EX != 0 {
		body = append(body, "ex", fmt.Sprintf("%d", options.EX))
	} else if options.EXAT != 0 {
		body = append(body, "exat", fmt.Sprintf("%d", options.EXAT))

	} else if options.PX != 0 {
		body = append(body, "px", fmt.Sprintf("%d", options.PX))

	} else if options.PXAT != 0 {
		body = append(body, "pxat", fmt.Sprintf("%d", options.PXAT))

	} else if options.PERSIST {
		body = append(body, "persist")
	}

	res, err := u.client.Call(client.Request{
		Body: body,
	})
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}

	return res.(string), nil
}

// Returns the substring of the string value stored at key, determined by
// the offsets start and end (both are inclusive). Negative offsets can be
// used in order to provide an offset starting from the end of the string.
// So -1 means the last character, -2 the penultimate and so forth.
//
// The function handles out of range requests by limiting the resulting
// range to the actual length of the string.
//
// Returns the a part of value of `key`, or empty string when `key` does
// not exist
//
// https://redis.io/commands/getrange
func (u *Upstash) GetRange(key string, start int, end int) (string, error) {

	res, err := u.client.Call(client.Request{
		Body: []string{"getrange", key, fmt.Sprintf("%d", start), fmt.Sprintf("%d", end)},
	})
	if err != nil {
		return "", err
	}

	return res.(string), nil
}

// Atomically sets key to value and returns the old value stored at key.
// Returns an error when key exists but does not hold a string value. Any
// previous time to live associated with the key is discarded on successful
// SET operation.
//
// Returns the old value stored at key, or empty string when key did not exist.
//
// https://redis.io/commands/getset
func (u *Upstash) GetSet(key string, value string) (string, error) {
	res, err := u.client.Call(client.Request{
		Body: []string{"getset", key, value},
	})
	if err != nil {
		return "", err
	}

	return res.(string), nil
}

// Increments the number stored at key by one. If the key does not exist,
// it is set to 0 before performing the operation. An error is returned if
// the key contains a value of the wrong type or contains a string that can
// not be represented as integer. This operation is limited to 64 bit
// signed integers.
//
// Note: this is a string operation because Redis does not have a dedicated
// integer type. The string stored at the key is interpreted as a base-10
// 64 bit signed integer to execute the operation.
//
// Redis stores integers in their integer representation, so for string
// values that actually hold an integer, there is no overhead for storing
//  the string representation of the integer.
//
// Returns the value of key after the increment
//
// https://redis.io/commands/incr
func (u *Upstash) Incr(key string) (int, error) {

	res, err := u.client.Call(client.Request{
		Body: []string{"incr", key},
	})
	return int(res.(float64)), err
}

// Increments the number stored at key by increment. If the key does not
// exist, it is set to 0 before performing the operation. An error is
// returned if the key contains a value of the wrong type or contains a
// string that can not be represented as integer. This operation is limited
// to 64 bit signed integers.
//
// See INCR for extra information on increment/decrement operations.
//
// Returns the value of key after the increment
//
// https://redis.io/commands/incrby
func (u *Upstash) IncrBy(key string, increment int) (int, error) {
	res, err := u.client.Call(client.Request{
		Body: []string{"incrby", key, fmt.Sprintf("%d", increment)},
	})
	return int(res.(float64)), err

}

// Increment the string representing a floating point number stored at key
// by the specified increment. By using a negative increment value, the
// result is that the value stored at the key is decremented (by the obvious
// properties of addition). If the key does not exist, it is set to 0 before
// performing the operation. An error is returned if one of the following
// conditions occur:
//
// The key contains a value of the wrong type (not a string). The current key
// content or the specified increment are not parsable as a double precision
// floating point number. If the command is successful the new incremented value
// is stored as the new value of the key (replacing the old one), and returned
// to the caller as a string.
//
// Both the value already contained in the string key and the increment argument
// can be optionally provided in exponential notation, however the value
// computed after the increment is stored consistently in the same format, that
// is, an integer number followed (if needed) by a dot, and a variable number of
// digits representing the decimal part of the number. Trailing zeroes are
// always removed.
//
// The precision of the output is fixed at 17 digits after the decimal point
// regardless of the actual internal precision of the computation.
//
// Returns the value of key after the increment.
//
//https://redis.io/commands/incrbyfloat
func (u *Upstash) IncrByFloat(key string, increment float64) (float64, error) {
	res, err := u.client.Call(client.Request{
		Body: []string{"incrbyfloat", key, fmt.Sprintf("%f", increment)},
	})
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(res.(string), 10)
	if err != nil {
		return 0, err
	}
	return f, nil

}

// Returns the values of all specified keys. For every key that does not
// hold a string value or does not exist, the special value nil is returned.
// Because of this, the operation never fails.
//
// Returns a list of values at the specified keys.
//
// https://redis.io/commands/mget
func (u *Upstash) MGet(keys []string) ([]string, error) {
	res, err := u.client.Call(client.Request{
		Body: append([]string{"mget"}, keys...),
	})

	values := make([]string, len(keys))
	for i, value := range res.([]interface{}) {
		values[i] = fmt.Sprint(value)
	}

	return values, err
}

// Sets the given keys to their respective values. MSET replaces existing
// values with new values, just as regular SET. See MSETNX if you don't want
// to overwrite existing values.
//
// MSET is atomic, so all given keys are set at once. It is not possible for
// clients to see that some of the keys were updated while others are unchanged.
//
//Returns nil, MSET can't fail.
//
// https://redis.io/commands/mset
func (u *Upstash) MSet(kvPairs []KV) error {
	body := []string{"mset"}
	for _, kv := range kvPairs {
		body = append(body, kv.Key, kv.Value)
	}

	_, err := u.client.Call(client.Request{
		Body: body,
	})
	return err
}

// Sets the given keys to their respective values. MSETNX will not perform
//  any operation at all even if just a single key already exists.
//
// Because of this semantic MSETNX can be used in order to set different
// keys representing different fields of an unique logic object in a way
// that ensures that either all the fields or none at all are set.
//
// MSETNX is atomic, so all given keys are set at once. It is not possible
// for clients to see that some of the keys were updated while others are
// unchanged.
//
// Returns:
// 1 if the all the keys were set.
// 0 if no key was set (at least one key already existed
//
// https://redis.io/commands/msetnx
func (u *Upstash) MSetNX(kvPairs []KV) (int, error) {
	body := []string{"msetnx"}
	for _, kv := range kvPairs {
		body = append(body, kv.Key, kv.Value)
	}

	res, err := u.client.Call(client.Request{
		Body: body,
	})
	return int(res.(float64)), err
}

// PSETEX works exactly like SETEX with the sole difference that the expire
// time is specified in milliseconds instead of seconds.
func (u *Upstash) PSetEX(key string, milliseconds int, value string) error {
	_, err := u.client.Call(client.Request{
		Body: []string{"psetex", key, fmt.Sprintf("%d", milliseconds), value},
	})
	return err
}

// Set key to hold the string value. If key already holds a value, it is
// overwritten, regardless of its type. Any previous time to live associated
// with the key is discarded on successful SET operation.
//
// https://redis.io/commands/set
func (u *Upstash) Set(key string, value string) error {
	_, err := u.client.Call(client.Request{
		Body: []string{"set", key, value},
	})
	return err
}

// Same as Set but with additional options
//
// https://redis.io/commands/set
func (u *Upstash) SetWithOptions(key string, value string, options SetOptions) error {
	body := []string{"set", key}
	if options.EX != 0 {
		body = append(body, "ex", fmt.Sprintf("%d", options.EX))
	} else if options.EXAT != 0 {
		body = append(body, "exat", fmt.Sprintf("%d", options.EXAT))

	} else if options.KEEPTTL {
		body = append(body, "keepttl")

	} else if options.PX != 0 {
		body = append(body, "px", fmt.Sprintf("%d", options.PX))
	} else if options.PXAT != 0 {
		body = append(body, "pxat", fmt.Sprintf("%d", options.PXAT))
	}
	if options.NX {
		body = append(body, "nx")
	} else if options.XX {
		body = append(body, "xx")
	}

	if options.GET {
		body = append(body, "get")

	}
	_, err := u.client.Call(client.Request{
		Body: body,
	})
	return err

}

// Set key to hold the string value and set key to timeout after a given
// number of seconds. This command is equivalent to executing the following
// commands:
//      SET mykey value
//      EXPIRE mykey seconds
//
// SETEX is atomic, and can be reproduced by using the previous two commands
// inside an MULTI / EXEC block. It is provided as a faster alternative to
// the given sequence of operations, because this operation is very common
// when Redis is used as a cache.
//
// An error is returned when seconds is invalid.
//
// https://redis.io/commands/setex
func (u *Upstash) SetEX(key string, seconds int, value string) error {

	_, err := u.client.Call(client.Request{
		Body: []string{"setex", key, fmt.Sprintf("%d", seconds), value},
	})
	return err

}

// Set key to hold string value if key does not exist. In that case, it is
// equal to SET. When key already holds a value, no operation is performed.
// SETNX is short for "SET if Not eXists".
//
// Return:
// - 1 if the key was set
// - 0 if the key was not set
//
// https://redis.io/commands/setnx
func (u *Upstash) SetNX(key string, value string) (int, error) {

	res, err := u.client.Call(client.Request{
		Body: []string{"setnx", key, value},
	})
	return int(res.(float64)), err

}

// Overwrites part of the string stored at key, starting at the specified
// offset, for the entire length of value. If the offset is larger than the
// current length of the string at key, the string is padded with zero-bytes
// to make offset fit. Non-existing keys are considered as empty strings, so
// this command will make sure it holds a string large enough to be able to
// set value at offset.
//
// Note that the maximum offset that you can set is 229 -1 (536870911), as
// Redis Strings are limited to 512 megabytes. If you need to grow beyond
// this size, you can use multiple keys.
//
// Warning: When setting the last possible byte and the string value stored
// at key does not yet hold a string value, or holds a small string value,
// Redis needs to allocate all intermediate memory which can block the
// server for some time. On a 2010 MacBook Pro, setting byte number
// 536870911 (512MB allocation) takes ~300ms, setting byte number 134217728
// (128MB allocation) takes ~80ms, setting bit number 33554432 (32MB
// allocation) takes ~30ms and setting bit number 8388608 (8MB allocation)
// takes ~8ms. Note that once this first allocation is done, subsequent
// calls to SETRANGE for the same key will not have the allocation overhead.
//
// Returns the length of the string after it was modified by the command.
//
// https://redis.io/commands/setrange
func (u *Upstash) SetRange(key string, offset int, value string) error {

	_, err := u.client.Call(client.Request{
		Body: []string{"setrange", key, fmt.Sprintf("%d", offset), value},
	})
	return err

}

// Returns the length of the string value stored at key. An error is
// returned when key holds a non-string value.
//
// Returns the length of the string at key, or 0 when key does not exist.
//
// https://redis.io/commands/strlen
func (u *Upstash) StrLen(key string) (int, error) {
	res, err := u.client.Call(client.Request{
		Body: []string{"strlen", key},
	})

	return int(res.(float64)), err
}
