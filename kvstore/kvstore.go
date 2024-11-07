package kvstore

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type DataEntry struct {
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type KVStore struct {
	data     sync.Map // klucz: string, wartość: DataEntry
	filePath string
	done     chan struct{}
}

func NewKVStore(filePath string) (*KVStore, error) {
	store := &KVStore{
		filePath: filePath,
		done:     make(chan struct{}),
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		if err := store.load(); err != nil {
			return nil, err
		}
	}

	go store.periodicSave()
	return store, nil
}

// Set zapisuje wartość z automatycznym timestampem
func (kv *KVStore) Set(key string, value string) {
	entry := DataEntry{
		Value:     value,
		Timestamp: time.Now(),
	}
	kv.data.Store(key, entry)
}

// SetWithTimestamp pozwala na ustawienie własnego timestampa
func (kv *KVStore) SetWithTimestamp(key string, value string, timestamp time.Time) {
	entry := DataEntry{
		Value:     value,
		Timestamp: timestamp,
	}
	kv.data.Store(key, entry)
}

// Get zwraca DataEntry dla danego klucza
func (kv *KVStore) Get(key string) (DataEntry, bool) {
	if value, ok := kv.data.Load(key); ok {
		return value.(DataEntry), true
	}
	return DataEntry{}, false
}

// GetValue zwraca tylko wartość (string) dla danego klucza
func (kv *KVStore) GetValue(key string) (string, bool) {
	if value, ok := kv.data.Load(key); ok {
		return value.(DataEntry).Value, true
	}
	return "", false
}

func (kv *KVStore) Delete(key string) {
	kv.data.Delete(key)
}

func (kv *KVStore) periodicSave() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			kv.save()
		case <-kv.done:
			return
		}
	}
}

func (kv *KVStore) load() error {
	file, err := os.ReadFile(kv.filePath)
	if err != nil {
		return err
	}

	var data map[string]DataEntry
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}

	for k, v := range data {
		kv.data.Store(k, v)
	}

	return nil
}

func (kv *KVStore) save() error {
	data := make(map[string]DataEntry)
	kv.data.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			data[k] = value.(DataEntry)
		}
		return true
	})

	file, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(kv.filePath, file, 0644)
}

func (kv *KVStore) Close() error {
	close(kv.done)
	return kv.save()
}

// Dodatkowe pomocnicze metody

// GetAll zwraca wszystkie wpisy
func (kv *KVStore) GetAll() map[string]DataEntry {
	result := make(map[string]DataEntry)
	kv.data.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			result[k] = value.(DataEntry)
		}
		return true
	})
	return result
}

// GetEntriesAfter zwraca wpisy po określonej dacie
func (kv *KVStore) GetEntriesAfter(timestamp time.Time) map[string]DataEntry {
	result := make(map[string]DataEntry)
	kv.data.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			entry := value.(DataEntry)
			if entry.Timestamp.After(timestamp) {
				result[k] = entry
			}
		}
		return true
	})
	return result
}
