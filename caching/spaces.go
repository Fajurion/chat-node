package caching

import (
	"chat-node/util"
	"log"
	"os"
	"strconv"

	"github.com/dgraph-io/ristretto"
)

// ! Always use cost 1
var spacesCache *ristretto.Cache // Account ID -> Space Info
var spaceApp uint

func setupCallsCache() {
	var err error

	// TODO: Check if values really are enough
	spacesCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,     // number of objects expected (100,000).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // something from the docs
	})
	if err != nil {
		panic(err)
	}

	app, err := strconv.Atoi(os.Getenv("SPACES_APP"))
	if err != nil {
		panic(err)
	}
	spaceApp = uint(app)
}

type SpaceInfo struct {
	Account   string
	Connected bool
	Data      string // Encrypted data
}

// Join a space
func JoinSpace(accId string, space string, data string) error {
	return nil
}

// Create a space
func CreateSpace(accId string, cluster uint) (util.AppToken, bool) {

	_, ok := spacesCache.Get(accId)
	if ok {
		return util.AppToken{}, false
	}

	if os.Getenv("SPACES_APP") == "" {
		log.Println("Spaces is currently disabled. Please set SPACES_APP in your .env file to enable it.")
		return util.AppToken{}, false
	}

	// Get new space
	token, err := util.ConnectToApp(accId, accId, spaceApp, cluster) // Use accId as roomId so it's unique
	if err != nil {
		log.Println("Error while connecting to Spaces:", err)
		return util.AppToken{}, false
	}
	spacesCache.Set(accId, SpaceInfo{
		Account:   accId,
		Connected: false,
		Data:      "",
	}, 1)

	return token, true
}
