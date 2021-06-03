# Go Party Pool

![](./docs/gppl.png)

Go Party Pool helps write workers pool. It's as easy as starting a pool party (if you have pool).

```go
func hashPassword(i interface{}) interface{} {
	h := sha256.New()
	h.Write([]byte(i.(string)))
	return string(hex.EncodeToString(h.Sum(nil)))
}

func main(){
	patryPool := pool.NewPool(12, hashPassword)
	patryPool.Run(context.Background())

	println(patryPool.Do("admin123#").(string))
}
```

### Example with best practices:

```go
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

func ComputeHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return string(hex.EncodeToString(h.Sum(nil)))
}

const WorkersNumber = 6

var (
	onceFastComputeHash sync.Once
	doFastComputeHash   func(in string) string
)

func FastComputeHash(password string) string {
	onceFastComputeHash.Do(func() {
		patryPool := pool.NewPool(WorkersNumber, func(i interface{}) interface{} {
			return ComputeHash(i.(string))
		})
		patryPool.Run(context.Background())
		doFastComputeHash = func(in string) string {
			return patryPool.Do(in).(string)
		}
	})
	return doFastComputeHash(password)
}

func main() {
	println(FastComputeHash("admin123"))

	wg := &sync.WaitGroup{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 60; i++ {
		wg.Add(1)
		go func() {
			randomData := make([]byte, 10)
			rand.Read(randomData)
			checksum := FastComputeHash(string(randomData))
			println(checksum)
			wg.Done()
		}()
	}
	wg.Wait()
}

```