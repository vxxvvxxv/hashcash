package db

import (
	"sync"
	"testing"

	"github.com/vxxvvxxv/hashcash/internal/logger"
)

func TestRandID(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}

	checkRandID := func(t *testing.T, min, max int) {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			got := randID()
			if got < min || got > max {
				t.Errorf("ErrorLevel: got %d, want >= %d && <= %d", got, min, max)
			}
		}
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go checkRandID(t, minID, maxID)
	}

	wg.Wait()
}

func TestDB_GetRandomDataFromDB(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}

	s, err := NewDBService(logger.NewLogger(logger.DebugLevel))
	if err != nil {
		t.Fatal(err)
	}

	if err = s.FillTestData(); err != nil {
		t.Fatal(err)
	}

	getRandData := func(t *testing.T) {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			got := s.GetRandomDataFromDB()
			if len(got) == 0 {
				t.Errorf("ErrorLevel: returned empty: %s", got)
			}
		}
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go getRandData(t)
	}

	wg.Wait()
}
