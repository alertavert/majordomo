// Author: M. Massenzio (marco@alertavert.com), 5/3/25

package conversations_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alertavert/gpt4-go/pkg/config"
	"github.com/alertavert/gpt4-go/pkg/conversations"
)

var _ = Describe("ThreadStore", func() {
	var (
		tempDir     string
		threadStore *conversations.ThreadStore
		testConfig  *config.Config
		projectName string
		testThread  conversations.Thread
	)

	BeforeEach(func() {
		// Create a temporary directory for testing
		var err error
		tempDir, err = os.MkdirTemp("", "threadstore_test")
		Expect(err).NotTo(HaveOccurred())

		// Set up a test configuration
		testConfig = &config.Config{
			ThreadsLocation: filepath.Join(tempDir, "threads.json"),
		}

		// Initialize the ThreadStore
		threadStore = conversations.NewThreadStore(testConfig)
		Expect(threadStore).NotTo(BeNil())

		// Test data
		projectName = "TestProject"
		testThread = conversations.Thread{
			ID:          "123",
			Name:        "Test Thread",
			Assistant:   "Test Assistant",
			Description: "A test thread description",
		}
	})

	AfterEach(func() {
		// Clean up the temporary directory
		Expect(os.RemoveAll(tempDir)).NotTo(HaveOccurred())
	})

	Describe("AddThread", func() {
		It("should add a thread and persist it to disk", func() {
			err := threadStore.AddThread(projectName, testThread)
			Expect(err).NotTo(HaveOccurred())

			// Verify the thread was added
			threads := threadStore.GetAllThreads(projectName)
			Expect(threads).To(HaveLen(1))
			Expect(threads[0]).To(Equal(testThread))
		})
	})

	Describe("GetAllThreads", func() {
		It("should return an empty list if no threads exist for a project", func() {
			threads := threadStore.GetAllThreads("NonExistentProject")
			Expect(threads).To(BeEmpty())
		})
	})

	Describe("Persistence", func() {
		It("should save one thread to disk", func() {
			// Add a thread and reload the store
			err := threadStore.AddThread(projectName, testThread)
			Expect(err).NotTo(HaveOccurred())

			newStore := conversations.NewThreadStore(testConfig)
			Expect(newStore).NotTo(BeNil())

			// Verify the thread was loaded
			threads := newStore.GetAllThreads(projectName)
			Expect(threads).To(HaveLen(1))
			Expect(threads[0]).To(Equal(testThread))
		})
		It("should persist multiple threads to disk", func() {
			threads := []conversations.Thread{
				{
					ID:          "thread-1",
					Name:        "First Thread",
					Assistant:   "Test Assistant",
					Description: "First test thread",
				},
				{
					ID:          "thread-2",
					Name:        "Second Thread",
					Assistant:   "Test Assistant",
					Description: "Second test thread",
				},
				{
					ID:          "thread-3",
					Name:        "Third Thread",
					Assistant:   "Test Assistant",
					Description: "Third test thread",
				},
			}

			// Add all threads to the store
			for _, thread := range threads {
				err := threadStore.AddThread(projectName, thread)
				Expect(err).NotTo(HaveOccurred())
			}

			// Create a new store instance and load from disk
			newStore := conversations.NewThreadStore(testConfig)
			Expect(newStore).NotTo(BeNil())

			// Verify all threads were loaded
			loadedThreads := newStore.GetAllThreads(projectName)
			Expect(loadedThreads).To(HaveLen(3))
			Expect(loadedThreads).To(ConsistOf(threads))
		})
	})
	Describe("GetThread", func() {
		It("should return the correct thread when it exists", func() {
			err := threadStore.AddThread(projectName, testThread)
			Expect(err).NotTo(HaveOccurred())

			thread, found := threadStore.GetThread(projectName, testThread.ID)
			Expect(found).To(BeTrue())
			Expect(thread).To(Equal(testThread))
		})

		It("should return false when thread doesn't exist", func() {
			_, found := threadStore.GetThread(projectName, "nonexistent-id")
			Expect(found).To(BeFalse())
		})
	})

	Describe("RemoveThread", func() {
		It("should remove an existing thread", func() {
			err := threadStore.AddThread(projectName, testThread)
			Expect(err).NotTo(HaveOccurred())

			removed, err := threadStore.RemoveThread(projectName, testThread.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(removed).To(BeTrue())

			threads := threadStore.GetAllThreads(projectName)
			Expect(threads).To(BeEmpty())
		})

		It("should return false when removing non-existent thread", func() {
			removed, err := threadStore.RemoveThread(projectName, "nonexistent-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(removed).To(BeFalse())
		})
	})
})
