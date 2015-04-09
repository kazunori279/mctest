package main

import (
	"fmt"
	"os"
	"bytes"
	"time"
    "crypto/rand"
    "encoding/base64"
    "github.com/bradfitz/gomemcache/memcache"
    "github.com/montanaflynn/stats"
)

func main() {
	mc := memcache.New("localhost:11211")
	n := 10000
	timeToSet := make([]float64, n)
	timeToGet := make([]float64, n)
	for i := 0; i < n; i++ {
		timeToSet[i], timeToGet[i] = measureSetAndGetTime(mc)
	}
	fmt.Printf("Set min: %v, max: %v, 90%%: %v, 99%%: %v \n", 
		stats.Min(timeToSet), stats.Max(timeToSet),
		stats.Percentile(timeToSet, 90), stats.Percentile(timeToSet, 99))
}

func measureSetAndGetTime(mc *memcache.Client) (float64, float64) {

	// create key and value
	k := randStr(10)
	v := []byte(randStr(1024))

	// test setting
	start := time.Now().Nanosecond()
	mc.Set(&memcache.Item{Key: k, Value: v})
	timeToSet := float64((time.Now().Nanosecond() - start)) / 1000.0

	// test getting
	start = time.Now().Nanosecond()
	it, err := mc.Get(k)
        if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	timeToGet := float64((time.Now().Nanosecond() - start)) / 1000.0

	// validate the value
	if !bytes.Equal(it.Value, v) {
		fmt.Printf("Wrong value: %v\n", it.Value)
		os.Exit(1)
	}
	return timeToSet, timeToGet
}

func randStr(l int) string {
	v := make([]byte, l)
	rand.Read(v)
	return base64.URLEncoding.EncodeToString(v)
}
