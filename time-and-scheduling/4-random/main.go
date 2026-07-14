package main

import (
	"fmt"
	"math/rand/v2"
)

/*
RANDOM NUMBERS (math/rand/v2)
=============================

Go 1.22+ ships math/rand/v2 — cleaner and automatically seeded, so the
top-level functions give different results on every run:

  rand.IntN(n)    -> int in [0, n)
  rand.Float64()  -> float in [0.0, 1.0)
  rand.Shuffle(n, swap) -> shuffle in place
  rand.Perm(n)    -> a random permutation of [0, n)

For REPRODUCIBLE sequences (tests, simulations) create your own generator with a
fixed seed via a source such as PCG:

  rng := rand.New(rand.NewPCG(seedA, seedB))

Same seed -> same sequence, every time.

NOTE: math/rand is NOT cryptographically secure. For tokens, passwords, or keys
use crypto/rand instead.
*/

func main() {
	// --- 1. Auto-seeded top-level helpers (vary each run) ---
	fmt.Println("five dice-ish rolls in [0,100):")
	for i := 0; i < 5; i++ {
		fmt.Printf("  %d\n", rand.IntN(100))
	}

	// --- 2. Shuffle a slice in place ---
	nums := []int{10, 20, 30, 40, 50}
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
	fmt.Printf("\nshuffled: %v\n", nums)

	// --- 3. A random permutation of indices 0..4 ---
	fmt.Printf("permutation: %v\n", rand.Perm(5))

	// --- 4. Seeded generators are reproducible ---
	fmt.Println("\ntwo generators with the SAME seed produce the SAME sequence:")
	rng1 := rand.New(rand.NewPCG(12345, 67890))
	rng2 := rand.New(rand.NewPCG(12345, 67890))
	for i := 0; i < 3; i++ {
		a, b := rng1.IntN(100), rng2.IntN(100)
		fmt.Printf("  rng1=%2d  rng2=%2d  (equal: %v)\n", a, b, a == b)
	}

	// --- 5. A different seed -> a different sequence ---
	rng3 := rand.New(rand.NewPCG(99, 100))
	fmt.Print("\ndifferent seed: ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d ", rng3.IntN(100))
	}
	fmt.Println()
}
