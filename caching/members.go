package caching

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"time"

	"github.com/dgraph-io/ristretto"
)

// ! Always use cost 1
var membersCache *ristretto.Cache // Conversation ID -> Members
const MemberTTL = time.Hour * 1   // 1 hour

func setupMembersCache() {
	var err error

	// TODO: Check if values really are enough
	membersCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6, // number of objects expected (1,000,000).
		MaxCost:     1e6, // maximum cost of cache (1,000,000).
		BufferItems: 64,  // something from the docs
	})
	if err != nil {
		panic(err)
	}
}

// Does database requests and stuff
func LoadMembers(id string) ([]string, error) {

	// Check cache
	if value, found := membersCache.Get(id); found {
		return value.([]string), nil
	}

	var members []string
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Select("id").Where("conversation = ?", id).Find(&members).Error; err != nil {
		return []string{}, err
	}

	// Add to cache
	membersCache.SetWithTTL(id, members, 1, MemberTTL)

	return members, nil
}
