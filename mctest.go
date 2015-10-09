package main

import (
	"bytes"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/montanaflynn/stats"
	"time"
	"encoding/binary"
	"strconv"
	"crypto/rand"
	"runtime"
)

func main() {
	mc := memcache.New("169.254.10.0:11211")
 	mc.Timeout = 20 * 1000 * time.Millisecond // 20 sec				
	initGRs := runtime.NumGoroutine()
        maxGRs := 8 
	for {
		time.Sleep(1000)
		for runtime.NumGoroutine() - initGRs < maxGRs {
			go measure(1000, mc)
		}
	}
}

func measure(n int, mc *memcache.Client) {
	timeToSet := make([]float64, n)
	timeToGet := make([]float64, n)
	for i := 0; i < n; i++ {
		timeToSet[i], timeToGet[i] = measureSetAndGetTime(mc)
		if i > 0 && timeToSet[i] < 0 {
			timeToSet[i] = timeToSet[i - 1] 
		}
		if i > 0 && timeToGet[i] < 0 {
			timeToGet[i] = timeToGet[i - 1]
		}
	}
	min, _ := stats.Min(timeToSet)
	max, _ := stats.Max(timeToSet)
	med, _ := stats.Median(timeToSet)
	p95, _ := stats.Percentile(timeToSet, 95)
	p99, _ := stats.Percentile(timeToSet, 99)
	fmt.Printf("Set min: %.2f, max: %.2f, median: %.2f, 95%%: %.2f, 99%%: %.2f \n",
		min, max, med, p95, p99)

	min, _ = stats.Min(timeToGet)
	max, _ = stats.Max(timeToGet)
	med, _ = stats.Median(timeToGet)
	p95, _ = stats.Percentile(timeToGet, 95)
	p99, _ = stats.Percentile(timeToGet, 99)
	fmt.Printf("Get min: %.2f, max: %.2f, median: %.2f, 95%%: %.2f, 99%%: %.2f \n",
		min, max, med, p95, p99)
}

func measureSetAndGetTime(mc *memcache.Client) (float64, float64) {

	// create key and value
	k := rand12Chars(1) // 12 bytes
	v := []byte(rand12Chars(85)) // 1020 bytes

	// test setting
	start := time.Now().Nanosecond()
	mc.Set(&memcache.Item{Key: k, Value: v})
	timeToSet := float64((time.Now().Nanosecond() - start)) / 1000.0

	// test getting
	start = time.Now().Nanosecond()
	it, err := mc.Get(k)
	if err != nil {
		fmt.Printf("%v\n", err)
	 	return -1, -1
	}
	timeToGet := float64((time.Now().Nanosecond() - start)) / 1000.0

	// validate the value
	if !bytes.Equal(it.Value, v) {
		fmt.Printf("Wrong value: %v\n", it.Value)
	 	return -1, -1
	}
	return timeToSet, timeToGet
}

func rand12Chars(l int) string {
	var buf bytes.Buffer	
	for i := 0; i < l; i++ {
		var n uint64
		binary.Read(rand.Reader, binary.LittleEndian, &n)
		buf.WriteString(strconv.FormatUint(n, 36))
	}
	return buf.String() 
}
