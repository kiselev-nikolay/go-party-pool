package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"sync"
	"time"

	pool "github.com/kiselev-nikolay/go-party-pool"
)

func main() {
	p := pool.NewPool(6, func(i interface{}) interface{} {
		h := sha256.New()
		h.Write([]byte(i.(string)))
		return string(hex.EncodeToString(h.Sum(nil)))
	})
	do := func(in string) string {
		return p.Do(in).(string)
	}
	p.Run(context.Background())

	checksum := do("krabs")
	print("krabs password hash = ")
	println(checksum)

	wg := &sync.WaitGroup{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 60; i++ {
		wg.Add(1)
		go func() {
			randomData := make([]byte, 10)
			rand.Read(randomData)
			checksum := do(string(randomData))
			println(checksum)
			wg.Done()
		}()
	}
	wg.Wait()
}
