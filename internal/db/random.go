package db

import "math/rand"

const (
	minID = 1
	maxID = 15
)

func randID() int {
	return minID + rand.Intn(maxID-minID)
}
