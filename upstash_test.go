package upstash_test

import (
	"testing"
	"time"

	"fmt"

	"github.com/chronark/upstash-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAppendSuccess(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	addition := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	length, err := u.Append(key, addition)
	require.NoError(t, err)
	require.Equal(t, 72, length)

	got, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s%s", value, addition), got)
}

func TestAppendFailure(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	length, err := u.Append(key, value)
	require.NoError(t, err)
	require.Equal(t, 36, length)

	got, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, value, got)
}

func TestDecr(t *testing.T) {
	key := uuid.NewString()
	value := "1"
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	after, err := u.Decr(key)
	require.NoError(t, err)

	require.Equal(t, 0, after)
}

func TestDecrBy(t *testing.T) {
	key := uuid.NewString()
	value := "5"
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	after, err := u.DecrBy(key, 4)
	require.NoError(t, err)

	require.Equal(t, 1, after)
}

func TestGet(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	got, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, value, got)
}

func TestGetRange(t *testing.T) {
	key := uuid.NewString()
	value := "abcde"
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	got, err := u.GetRange(key, 1, 3)
	require.NoError(t, err)
	require.Equal(t, "bcd", got)
}

func TestGetSet(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	value2 := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	got, err := u.GetSet(key, value2)
	require.NoError(t, err)
	require.Equal(t, value, got)

	got2, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, value2, got2)

}

func TestIncr(t *testing.T) {
	key := uuid.NewString()
	value := "1"
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	after, err := u.Incr(key)
	require.NoError(t, err)
	require.Equal(t, 2, after)
}

func TestIncrBy(t *testing.T) {
	key := uuid.NewString()
	value := "5"
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	after, err := u.IncrBy(key, 3)
	require.NoError(t, err)

	require.Equal(t, 8, after)
}
func TestIncrByFloat(t *testing.T) {
	key := uuid.NewString()
	value := "5"
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	after, err := u.IncrByFloat(key, 3.5)
	require.NoError(t, err)

	require.Equal(t, 8.5, after)
}

func TestMGet(t *testing.T) {
	key1 := uuid.NewString()
	key2 := uuid.NewString()
	value1 := uuid.NewString()
	value2 := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key1, value1)
	require.NoError(t, err)

	err = u.Set(key2, value2)
	require.NoError(t, err)

	got, err := u.MGet([]string{key1, key2})
	require.NoError(t, err)

	require.Equal(t, []string{value1, value2}, got)
}

func TestMSet(t *testing.T) {
	key1 := uuid.NewString()
	key2 := uuid.NewString()
	value1 := uuid.NewString()
	value2 := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.MSet([]upstash.KV{{Key: key1, Value: value1}, {Key: key2, Value: value2}})
	require.NoError(t, err)

	got1, err := u.Get(key1)
	require.NoError(t, err)

	require.Equal(t, value1, got1)

	got2, err := u.Get(key2)
	require.NoError(t, err)

	require.Equal(t, value2, got2)
}

func TestMSetNX(t *testing.T) {
	key1 := uuid.NewString()
	key2 := uuid.NewString()
	value1 := uuid.NewString()
	value2 := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	allSet, err := u.MSetNX([]upstash.KV{{Key: key1, Value: value1}, {Key: key2, Value: value2}})
	require.NoError(t, err)
	require.Equal(t, 1, allSet)

	got1, err := u.Get(key1)
	require.NoError(t, err)

	require.Equal(t, value1, got1)

	got2, err := u.Get(key2)
	require.NoError(t, err)

	require.Equal(t, value2, got2)
}
func TestPSetEX(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.PSetEX(key, 1000, value)
	require.NoError(t, err)

	got1, err := u.Get(key)
	require.NoError(t, err)

	require.Equal(t, value, got1)

	time.Sleep(2 * time.Second)

	got2, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, "", got2)
}

func TestSet(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	got, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, value, got)
}
func TestSetEX(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.SetEX(key, 1, value)
	require.NoError(t, err)

	got1, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, value, got1)

	time.Sleep(2 * time.Second)

	got2, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, "", got2)
}

func TestSetRange(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	overwrite := "HELLO"
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	err = u.SetRange(key, 4, overwrite)
	require.NoError(t, err)

	got, err := u.Get(key)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s%s%s", value[:4], overwrite, value[9:]), got)

}

func TestStrLen(t *testing.T) {
	key := uuid.NewString()
	value := uuid.NewString()
	u, _ := upstash.New(upstash.Options{})

	err := u.Set(key, value)
	require.NoError(t, err)

	res, err := u.StrLen(key)
	require.NoError(t, err)

	require.Equal(t, 36, res)

}
