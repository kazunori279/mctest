package main

import (
	"fmt"
	"os"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	mc := memcache.New("localhost:11211")
	mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})
	it, err := mc.Get("foo")
        if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", it)
}
