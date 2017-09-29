package search

import (
	"log"
	"sync"
)

// A map of registered matchers for searching
// This variable is located outside the scope of any function and so is considered a package-level variable.
var matchers = make(map[string]Matcher)

// Run performs the search logic
func Run(searchTerm string) {
	// Retrieve the list of feeds to search through
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	// Create a unbuffered channel to receive match results.
	results := make(chan *Result)

	// The next two lines of code are used later to prevent the program from terminating
	// before all the search processing is complete
	// In Go, once the main function returns, the program terminates.
	// Any goroutines that were launched and are still running at this time will also be terminated by the Go runtime
	// Setup a wait group so we can process all the feeds.
	var waitGroup sync.WaitGroup

	// Set the number of goroutines we need to wait for while they process the individual feeds
	waitGroup.Add(len(feeds))

	// Launch a goroutine for each feed to find the results.
	for _, feed := range feeds {
		// Retrieve a matcher for the search
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}

		// Launch the goroutine to perform the search.
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done()
		}(matcher, feed)
	}

	// Launch a goroutine to monitor when all the work is done.
	go func() {
		// Wait for everything to be processed
		waitGroup.Wait()

		// Close the channel to signal to the Display function that we can exit the program.
		close(results)
	}()

	// Start displaying results as they are available and return after the final result is displayed
	Display(results)
}
