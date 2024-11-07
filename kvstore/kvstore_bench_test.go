package kvstore

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkKVStore_Set_Parallel(b *testing.B) {
	store, err := NewKVStore("benchmark_test.db")
	if err != nil {
		b.Fatalf("Failed to create KVStore: %v", err)
	}
	defer os.Remove("benchmark_test.db")
	defer store.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			store.Set(fmt.Sprintf("key-%d", counter), strconv.Itoa(counter))
			counter++
		}
	})
}

func BenchmarkKVStore_Get_Parallel(b *testing.B) {
	store, err := NewKVStore("benchmark_test.db")
	if err != nil {
		b.Fatalf("Failed to create KVStore: %v", err)
	}
	defer os.Remove("benchmark_test.db")
	defer store.Close()

	// Przygotowanie danych
	for i := 0; i < 1000; i++ {
		store.Set(fmt.Sprintf("key-%d", i), strconv.Itoa(i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			store.Get(fmt.Sprintf("key-%d", counter%1000))
			counter++
		}
	})
}

func BenchmarkKVStore_Mixed_Parallel(b *testing.B) {
	store, err := NewKVStore("benchmark_test.db")
	if err != nil {
		b.Fatalf("Failed to create KVStore: %v", err)
	}
	defer os.Remove("benchmark_test.db")
	defer store.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := "key-" + strconv.Itoa(counter)
			switch counter % 3 {
			case 0:
				store.Set(key, strconv.Itoa(counter))
			case 1:
				store.Get(key)
			case 2:
				store.Delete(key)
			}
			counter++
		}
	})
}

func BenchmarkKVStore_Concurrent_HeavyLoad(b *testing.B) {
	store, err := NewKVStore("benchmark_test.db")
	if err != nil {
		b.Fatalf("Failed to create KVStore: %v", err)
	}
	defer os.Remove("benchmark_test.db")
	defer store.Close()

	var wg sync.WaitGroup
	workers := 100
	opsPerWorker := b.N / workers

	b.ResetTimer()

	// Writers
	for i := 0; i < workers/2; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				key := fmt.Sprintf("key-%d-%d", workerID, j)
				store.Set(key, strconv.Itoa(j))
			}
		}(i)
	}

	// Readers
	for i := workers / 2; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				key := fmt.Sprintf("key-%d-%d", workerID%50, j)
				store.Get(key)
			}
		}(i)
	}

	wg.Wait()
}
