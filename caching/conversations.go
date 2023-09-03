package caching

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
)

// ! Always use cost 1
var conversationsCache *ristretto.Cache // Conversation token ID -> Conversation Token
const ConversationTTL = time.Hour * 1   // 1 hour

// Errors
var ErrInvalidToken = errors.New("invalid")

func setupConversationsCache() {
	var err error

	// TODO: Check if values really are enough
	conversationsCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6, // number of objects expected (1,000,000).
		MaxCost:     1e6, // maximum cost of cache (1,000,000).
		BufferItems: 64,  // something from the docs
	})
	if err != nil {
		panic(err)
	}
}

// This does database requests and stuff
func ValidateToken(id string, token string) (conversations.ConversationToken, error) {

	// Check cache
	if value, found := conversationsCache.Get(id); found {

		// Check if token is valid
		if value.(conversations.ConversationToken).Token != token {
			return conversations.ConversationToken{}, ErrInvalidToken
		}

		return value.(conversations.ConversationToken), nil
	}

	var conversationToken conversations.ConversationToken
	if err := database.DBConn.Where("id = ?", id).Take(&conversationToken).Error; err != nil {
		return conversations.ConversationToken{}, err
	}

	// Add to cache
	conversationsCache.SetWithTTL(id, conversationToken, 1, ConversationTTL)

	if conversationToken.Token != token {
		return conversations.ConversationToken{}, ErrInvalidToken
	}

	return conversationToken, nil
}
