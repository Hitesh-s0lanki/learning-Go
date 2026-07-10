package main

import (
	"fmt"
	"sync"
)

/*
sync.Once (and a note on atomic vs mutex vs once)
=================================================

sync.Once guarantees a piece of code runs EXACTLY ONCE, no matter how many
goroutines call it concurrently. The classic use is lazy, thread-safe
initialization (singletons, config loading, connection setup).

  var once sync.Once
  once.Do(func() { ... })   // the func runs on the first call only

All other callers of once.Do block until the first call's function returns,
then proceed without running it again.
*/

type Config struct {
	value string
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig lazily builds the Config the first time it is called, and safely
// returns the same instance to every caller afterwards.
func GetConfig() *Config {
	once.Do(func() {
		fmt.Println("  >> initializing config (this must print only ONCE)")
		instance = &Config{value: "loaded-from-somewhere"}
	})
	return instance
}

func main() {
	var wg sync.WaitGroup

	// Call GetConfig from 10 goroutines at the same time. The initializer
	// inside once.Do must run exactly once, and all goroutines must see the
	// same *Config pointer.
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c := GetConfig()
			fmt.Printf("goroutine %d got config value=%q (ptr=%p)\n", id, c.value, c)
		}(i)
	}

	wg.Wait()
	fmt.Println("done — the initializer ran only once even under concurrency")
}
