package caching

import (
	"chat-node/util"
	"log"
	"os"
	"strconv"

	"github.com/dgraph-io/ristretto"
)

// TODO: Move this into a shared redis instance
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
	Account      string
	ConnectionID string
	Domain       string
}

// Check if account is in a space
func IsInSpace(accId string) bool {
	_, ok := spacesCache.Get(accId)
	return ok
}

// Leave a space
func LeaveSpace(accId string) bool {

	obj, ok := spacesCache.Get(accId)
	if !ok {
		return false
	}
	space := obj.(SpaceInfo)

	// Disconnect from space
	body, err := util.PostRaw(util.Protocol+space.Domain+"/leave", map[string]interface{}{
		"conn": space.ConnectionID,
	})
	if err != nil {
		log.Println("Error while leaving space:", err)
		return false
	}
	if !body["success"].(bool) {
		log.Println("Error while leaving space:", body["error"])
		return false
	}

	// Actually leave the space (this took 10 minutes to add because I'm stupid)
	spacesCache.Del(accId)

	return true
}

// Join a space
func JoinSpace(accId string, space string, cluster uint) (util.AppToken, bool) {

	_, ok := spacesCache.Get(accId)
	if ok {
		return util.AppToken{}, false
	}

	connId := generateConnectionID()
	token, err := util.ConnectToApp(connId, space, spaceApp, cluster) // Use accId as roomId so it's unique
	if err != nil {
		log.Println("Error while connecting to Spaces:", err)
		return util.AppToken{}, false
	}
	spacesCache.Set(accId, SpaceInfo{
		Account:      accId,
		ConnectionID: connId,
		Domain:       token.Domain,
	}, 1)

	return token, true
}

// Create a space
func CreateSpace(accId string, cluster uint) (string, util.AppToken, bool) {

	/*
		_, ok := spacesCache.Get(accId)
		if ok {
			return "", util.AppToken{}, false
		}
	*/

	if os.Getenv("SPACES_APP") == "" {
		log.Println("Spaces is currently disabled. Please set SPACES_APP in your .env file to enable it.")
		return "", util.AppToken{}, false
	}

	// Get new space
	connId := generateConnectionID()
	roomId := util.GenerateToken(16)
	token, err := util.ConnectToApp(connId, roomId, spaceApp, cluster) // Use accId as roomId so it's unique
	if err != nil {
		log.Println("Error while connecting to Spaces:", err)
		return "", util.AppToken{}, false
	}
	spacesCache.Set(accId, SpaceInfo{
		Account:      accId,
		ConnectionID: connId,
		Domain:       token.Domain,
	}, 1)

	return roomId, token, true
}

func generateConnectionID() string {
	return "sn-" + util.GenerateToken(16)
}
