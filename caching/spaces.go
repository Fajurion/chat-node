package caching

import (
	"github.com/dgraph-io/ristretto"
)

// ! Always use cost 1
var spacesCache *ristretto.Cache // Account ID -> Space Info

func setupCallsCache() {
	var err error

	// TODO: Check if values really are enough
	spacesCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e5, // number of objects expected (100,000).
		MaxCost:     1e5, // maximum cost of cache (100,000).
		BufferItems: 64,  // something from the docs
	})
	if err != nil {
		panic(err)
	}
}

type SpaceInfo struct {
	Account string
	Space   string // Space ID
	Start   int64  // Unix milli
}

// Join a space
func JoinSpace(accId string, space string) error {
	return nil
}

// Create a space
func CreateSpace(accId string) error {
	return nil
}
