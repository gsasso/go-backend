package ticker

import (
	"fmt"
	"sync"
	"time"
)

type Summary struct {
	totalUnits       int
	totalReached     int
	messagePerSecond int
	mu               sync.Mutex
}

type SummaryResponse struct {
	TotalUnits   int64
	TotalReached int64
}

var summaryResult Summary

type SummaryService struct{}

type SummaryInt interface {
	IncreaseTotalUnits()
	ResetMessagePerSecond()
	IncreaseTotalReached()
	GetSummary() (SummaryResponse, error)
	Tick()
}

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

func (s *SummaryService) GetSummary() (SummaryResponse, error) {
	summaryResult.mu.Lock()
	defer summaryResult.mu.Unlock()
	return SummaryResponse{
		TotalUnits:   int64(summaryResult.totalUnits),
		TotalReached: int64(summaryResult.totalReached),
	}, nil
}
