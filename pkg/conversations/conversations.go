package conversations

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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Assistant   string `json:"assistant"`
	Description string `json:"description"`
}

// ThreadsMap is a map of Project names to their respective Threads.
type ThreadsMap map[string][]Thread

// ThreadStore encapsulates behavior for managing and persisting conversations.
type ThreadStore struct {
	location   string
	threadsMap ThreadsMap
	mu         sync.Mutex
}

// NewThreadStore creates a new ThreadStore instance using the provided configuration.
func NewThreadStore(cfg *config.Config) *ThreadStore {
	if cfg.ThreadsLocation == "" {
		log.Error().Msg("Threads location not configured")
		return nil
	}
	ts := &ThreadStore{
		threadsMap: make(ThreadsMap),
		location:   cfg.ThreadsLocation,
	}
	if err := ts.load(); err != nil {
		log.Error().
			Err(err).
			Str("location", ts.location).
			Msg("Error loading conversations")
		return nil
	}
	log.Info().
		Str("location", ts.location).
		Int("conversations", len(ts.threadsMap)).
		Msg("Loaded thread store from disk")
	return ts
}

// AddThread adds a new thread to the thread map and persists the map to storage.
func (ts *ThreadStore) AddThread(projectName string, thread Thread) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	// TODO: do we need to check first if the key exists in the map?
	ts.threadsMap[projectName] = append(ts.threadsMap[projectName], thread)
	return ts.save()
}

func (ts *ThreadStore) GetAllThreads(projectName string) []Thread {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.threadsMap[projectName]
}

// GetThread retrieves a specific thread from a project by its ID.
// Returns the thread and true if found, or an empty thread and false if not found.
func (ts *ThreadStore) GetThread(projectName string, threadID string) (Thread, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	threads := ts.threadsMap[projectName]
	for _, thread := range threads {
		if thread.ID == threadID {
			return thread, true
		}
	}
	return Thread{}, false
}

// RemoveThread removes a specific thread from a project.
// Returns true if the thread was found and removed, false otherwise.
func (ts *ThreadStore) RemoveThread(projectName string, threadID string) (bool, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	threads := ts.threadsMap[projectName]
	for i, thread := range threads {
		if thread.ID == threadID {
			// Remove the thread by slicing
			ts.threadsMap[projectName] = append(threads[:i], threads[i+1:]...)
			return true, ts.save()
		}
	}
	return false, nil
}

// load retrieves the conversations data from the disk.
func (ts *ThreadStore) load() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	file, err := os.Open(ts.location)
	if err != nil {
		if os.IsNotExist(err) {
			ts.threadsMap = make(ThreadsMap)
			return nil
		}
		log.Error().Err(err).
			Str("location", ts.location).
			Msg("Error opening conversations file")
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).
				Str("location", ts.location).
				Msg("Error occurred while closing conversations file")
		}
	}()
	return json.NewDecoder(file).Decode(&ts.threadsMap)
}

// save persists the current state of the conversations map to the disk.
func (ts *ThreadStore) save() error {
	file, err := os.Create(ts.location)
	if err != nil {
		log.Error().Err(err).
			Str("location", ts.location).
			Msg("Error while saving conversations")
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).
				Str("location", ts.location).
				Msg("Error while closing saved conversations file")
		}
	}()
	return json.NewEncoder(file).Encode(ts.threadsMap)
}
