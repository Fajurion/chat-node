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

type StoredMember struct {
	Token string
	Node  int64
}

// Does database requests and stuff
func LoadMembers(id string) ([]StoredMember, error) {

	// Check cache
	if value, found := membersCache.Get(id); found {
		return value.([]StoredMember), nil
	}

	var members []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation = ?", id).Find(&members).Error; err != nil {
		return []StoredMember{}, err
	}

	storedMembers := make([]StoredMember, len(members))
	for i, member := range members {
		storedMembers[i] = StoredMember{
			Token: member.Token,
			Node:  member.Node,
		}
	}

	// Add to cache
	membersCache.SetWithTTL(id, storedMembers, 1, MemberTTL)

	return storedMembers, nil
}

func LoadMembersArray(ids []string) (map[string][]StoredMember, error) {

	// Check cache
	returnMap := make(map[string][]StoredMember, len(ids)) // Conversation ID -> Members
	notFound := []string{}

	for _, id := range ids {
		if value, found := membersCache.Get(id); found {
			returnMap[id] = value.([]StoredMember)
		} else {
			notFound = append(notFound, id)
		}
	}

	var tokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation IN ?", notFound).Find(&tokens).Error; err != nil {
		return nil, err
	}
	for _, token := range tokens {
		returnMap[token.Conversation] = append(returnMap[token.Conversation], StoredMember{
			Token: token.Token,
			Node:  token.Node,
		})
	}
	for key, memberTokens := range returnMap {
		membersCache.SetWithTTL(key, memberTokens, 1, MemberTTL)
	}

	return returnMap, nil
}
