package ticker

import (
	"fmt"
	"sync"
	"time"
)

// TODO Minor: Better to keep clean package and have separate package for models like in this case - MetricSummary for example

type Summary struct {
	// TODO Minor: Int is Int32 and have range: -2147483648 through 2147483647.
	// TODO Minor: Int is Int32 by design is it planned to go below `0`?
	// TODO Minor: Int is Int32 why not to use: uint64 that have range from 0 through 18446744073709551615 ?
	totalUnits       int
	totalReached     int
	messagePerSecond int
	mu               sync.Mutex
}

// TODO Major: This naming will mislead, since it has postfix response - in this application domain you have controllers where responses are constructed.
type SummaryResponse struct {
	TotalUnits   int64
	TotalReached int64
}

// TODO Major: Global VAR - try to avoid using global variables for several reasons
// TODO Major: Global VAR - 1. Code Readability, shared and used across multiple files
// TODO Major: Global VAR - 2. Data Consistency, multiple pieces of code may read from and write
// TODO Major: Global VAR - 3. Testing, global variables make it difficult to write unit tests
// TODO Major: Global VAR - 4. Encapsulation, all parts of a program direct access to a piece of data - not OOP principle
// TODO Major: Global VAR - 5. Not a thread safe

var summaryResult Summary

type SummaryService struct{}

// TODO Major: I'm already not understanding what `SummaryInt` is meant to do by name.
type SummaryInt interface {
	// TODO Minor: By method names it's looks like there is more than one purpose of service, it's collecting something, then it trying to be a timer and aggregator.
	IncreaseTotalUnits()
	ResetMessagePerSecond()
	IncreaseTotalReached()
	GetSummary() (SummaryResponse, error)
	Tick()
}

// TODO Major: Tick SRP Broken: Breaking single responsibility principle here is why:
// TODO Major: Tick SRP Broken: 1. Creating ticker locally
// TODO Major: Tick SRP Broken: 2. Manipulating `summaryResult`
// TODO Major: Tick SRP Broken: 3. Printing output that is related not to Tick
// TODO Major: Tick SRP Broken: Now what if, you need to have 2 ways of outputting `summaryResult` to STDOUT and calling another API?

func (s *SummaryService) Tick() {

	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go func() {
	loop:
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				// TODO: Minor: There could be a issue with 'summaryResult' if it's not properly initialized and has not required fields

				// TODO: Critical: Code not protected from Race Conditions
				// TODO: Critical: Concurrently occur if IncreaseTotalUnits(), IncreaseTotalReached(), ResetMessagePerSecond(), or GetSummary() called
				// TODO: Critical: Where `summaryResult` used and Race Conditions can happen
				fmt.Println("Messages per second at ", t, "are: ", summaryResult.messagePerSecond)
				if summaryResult.totalUnits != 0 && summaryResult.messagePerSecond == 0 {
					fmt.Println("Total units reached: ", summaryResult.totalReached)
					fmt.Println("Total units processed: ", summaryResult.totalUnits)
					ticker.Stop()
					done <- true
					break loop
				}

				s.ResetMessagePerSecond()
			}
		}
	}()

}

// TODO Major: After reading whole `SummaryService` as a minimum, it will be better to move `summaryResult` from global to structure so that instance can control it but no one else.

func (s *SummaryService) IncreaseTotalUnits() {
	summaryResult.mu.Lock()
	defer summaryResult.mu.Unlock()
	summaryResult.totalUnits++
	summaryResult.messagePerSecond++
}

func (s *SummaryService) IncreaseTotalReached() {
	summaryResult.mu.Lock()
	defer summaryResult.mu.Unlock()
	summaryResult.totalReached++
	summaryResult.messagePerSecond++
}

func (s *SummaryService) ResetMessagePerSecond() {
	summaryResult.mu.Lock()
	defer summaryResult.mu.Unlock()
	summaryResult.messagePerSecond = 0
}

// TODO Minor: Where is point to return error if it's always == NIL ?

func (s *SummaryService) GetSummary() (SummaryResponse, error) {
	summaryResult.mu.Lock()
	defer summaryResult.mu.Unlock()
	return SummaryResponse{
		TotalUnits:   int64(summaryResult.totalUnits),
		TotalReached: int64(summaryResult.totalReached),
	}, nil
}
