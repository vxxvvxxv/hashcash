package db

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/vxxvvxxv/hashcash/internal/logger"
)

type dbService struct {
	db map[int]string
	mu sync.RWMutex

	log logger.Logger
}

func NewDBService(log logger.Logger) (DB, error) {
	s := &dbService{
		db:  make(map[int]string),
		log: log,
	}

	return s, nil
}

func (s *dbService) FillTestData() error {
	// Just for generating ids, should be set once
	rand.Seed(time.Now().UTC().UnixNano())

	dbTmp := make(map[string]string)

	// Read data from file
	// Don't need to mutex here, because it's just one time
	err := json.Unmarshal(dataJson, &dbTmp)
	if err != nil {
		s.log.Error("can't read data: %v", err)
		return err
	}

	// Create new map with int keys
	dbNewVersion := make(map[int]string)

	// Convert map[string]string => map[int]string
	for k, v := range dbTmp {
		id, errConv := strconv.Atoi(k)
		if errConv != nil {
			s.log.Error("can't convert key to int: %v", errConv)
			return errConv
		}
		dbNewVersion[id] = v
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Replace old map with new
	s.db = dbNewVersion

	return nil
}

func (s *dbService) GetDataFromDB(id int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.db[id]
	return v, ok
}

func (s *dbService) GetRandomDataFromDB() string {
	v, ok := s.GetDataFromDB(randID())
	if !ok {
		s.log.Error("can't get random data from db")
	}
	return v
}
