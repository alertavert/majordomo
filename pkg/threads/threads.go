package threads

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/rs/zerolog/log"
)

// Thread models simply the ID and name of the Thread.
// The ID is used to retrieve the Thread from the OpenAI API; while the name is
// used to display the Thread in the UI.
type Thread struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Assistant string `json:"assistant"`
	Description string `json:"description"`
}

// ThreadsMap is a map of Project names to their respective Threads.
// TODO: Will need to eventually store the Threads (per Project) in persistent storage.
type ThreadsMap map[string][]Thread

// ThreadStore encapsulates behavior for managing and persisting threads.
type ThreadStore struct {
	location  string
	threadsMap ThreadsMap
	mu         sync.Mutex
}

// NewThreadStore creates a new ThreadStore instance using the provided configuration.
func NewThreadStore(cfg *config.Config) *ThreadStore {
	ts := &ThreadStore{
		threadsMap: make(ThreadsMap),
		location:     cfg.ThreadsLocation,
	}

	if err := ts.load(); err != nil {
		log.Error().
			Err(err).
			Str("location", ts.location).
			Msg("Error loading threads")
		return nil
	}
	return ts
}

// AddThread adds a new thread to the thread map and persists the map to storage.
func (ts *ThreadStore) AddThread(projectName string, thread Thread) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.threadsMap[projectName] = append(ts.threadsMap[projectName], thread)
	return ts.save()
}

func (ts *ThreadStore) GetThreads(projectName string) []Thread {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.threadsMap[projectName]
}

// load retrieves the threads data from the disk.
func (ts *ThreadStore) load() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	file, err := os.Open(ts.location)
	if err != nil {
		if os.IsNotExist(err) {
			ts.threadsMap = make(ThreadsMap)
			return nil
		}
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).
				Str("location", ts.location).
				Msg("Error opening threads file")
		}
	}()
	return json.NewDecoder(file).Decode(&ts.threadsMap)
}

// save persists the current state of the threads map to the disk.
func (ts *ThreadStore) save() error {
	file, err := os.Create(ts.location)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).
				Str("location", ts.location).
				Msg("Error saving threads file")
		}
	}()
	return json.NewEncoder(file).Encode(ts.threadsMap)
}
