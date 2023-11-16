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
		NumCounters: 1e7,     // number of objects expected (1,000,000).
		MaxCost:     1 << 30, // maximum cost of cache (1,000,000).
		BufferItems: 64,      // something from the docs
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

func ValidateTokens(tokens *[]conversations.SentConversationToken) ([]conversations.ConversationToken, error) {

	// Check cache
	foundTokens := []conversations.ConversationToken{}

	notFound := map[string]conversations.SentConversationToken{}
	notFoundIds := []string{}
	for _, token := range *tokens {
		if value, found := conversationsCache.Get(token.ID); found {
			if value.(conversations.ConversationToken).Token == token.Token {
				foundTokens = append(foundTokens, value.(conversations.ConversationToken))
			}
			continue
		} else {
			notFound[token.ID] = token
			notFoundIds = append(notFoundIds, token.ID)
		}
	}

	// Get tokens from database
	var conversationTokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", notFoundIds).Find(&conversationTokens).Error; err != nil {
		return nil, err
	}

	for _, token := range conversationTokens {
		conversationsCache.SetWithTTL(token.ID, token, 1, ConversationTTL)
		if token.Token == notFound[token.ID].Token {
			foundTokens = append(foundTokens, token)
		}
	}

	return foundTokens, nil
}

func UpdateToken(token conversations.ConversationToken) error {

	// Update cache
	conversationsCache.SetWithTTL(token.ID, token, 1, ConversationTTL)

	return nil
}
