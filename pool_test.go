package gopartypool_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	pool "github.com/kiselev-nikolay/go-party-pool"
)

func assert(t *testing.T, actual string, expect string) {
	if actual != expect {
		t.Errorf("expected `%s`, got `%s`", expect, actual)
	}
}

func TestPool(t *testing.T) {
	p := pool.NewPool(6, func(i interface{}) interface{} {
		<-time.After(time.Duration(rand.Intn(1000)) * time.Millisecond) // We all know that Squidward is just pretending to work
		return "oh yeah mr. " + i.(string)
	})
	p.Run(context.Background())
	do := func(in string) string {
		return p.Do(in).(string)
	}

	assert(t, "oh yeah mr. krabs", do("krabs"))

	wg := &sync.WaitGroup{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 60; i++ {
		wg.Add(1)
		go func() {
			randomData := make([]byte, 10)
			rand.Read(randomData)
			anchovyName := string(randomData)
			assert(t, "oh yeah mr. "+anchovyName, do(anchovyName))
			wg.Done()
		}()
	}
	wg.Wait()
}
