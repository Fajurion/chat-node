package conversation

import (
	"chat-node/caching"
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/util"
	"log"

	"github.com/Fajurion/pipes"
	"github.com/Fajurion/pipes/adapter"
	"github.com/Fajurion/pipes/send"
	"github.com/Fajurion/pipesfiber/wshandler"
)

// Action: conv_sub
func subscribe(message wshandler.Message) {

	if message.ValidateForm("tokens", "status", "date") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	date := int64(message.Data["date"].(float64))
	conversationTokens, tokenIds, members, ok := PrepareConversationTokens(message)
	if !ok {
		return
	}

	// Update all node IDs
	if database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", tokenIds).Update("node", util.NodeTo64(pipes.CurrentNode.ID)).Error != nil {
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	statusMessage := message.Data["status"].(string)
	readDates := make(map[string]int64, len(conversationTokens))
	adapters := make([]string, len(conversationTokens))
	for _, token := range conversationTokens {

		// Register adapter for the subscription
		adapter.AdaptWS(adapter.Adapter{
			ID: "s-" + token.Token,
			// TODO: Fix this not disconnecting
			Receive: func(ctx *adapter.Context) error {
				client := *message.Client
				log.Println(ctx.Adapter.ID, token.Token, client.ID)
				err := client.SendEvent(*ctx.Event)
				if err != nil {
					log.Println("COULDN'T SEND:", err.Error())
				}
				return err
			},
		})
		log.Println("SUB", "s-"+token.Token)
		adapters = append(adapters, "s-"+token.Token)

		var memberIds []string
		var memberNodes []string
		log.Printf("%d", len(members[token.Conversation]))
		if len(members[token.Conversation]) == 2 {
			for _, member := range members[token.Conversation] {
				if member.Token != token.Token {
					memberIds = append(memberIds, "s-"+member.Token)
					memberNodes = append(memberNodes, util.Node64(member.Node))
				}
			}
		}
		log.Printf("Sending to %d members", len(memberIds))

		// Send the subscription event
		send.Pipe(send.ProtocolWS, pipes.Message{
			Channel: pipes.Conversation(memberIds, memberNodes),
			Event: pipes.Event{
				Name: "acc_st",
				Data: map[string]interface{}{
					"st": statusMessage,
					"d":  "",
				},
			},
		})

		readDates[token.Conversation] = token.LastRead
		AddConversationToken(TokenTask{"s-" + token.Token, token.Conversation, date})
	}

	// Insert adapters into cache (to be deleted when disconnecting)
	caching.InsertAdapters(message.Client.ID, adapters)

	wshandler.NormalResponse(message, map[string]interface{}{
		"success": true,
		"read":    readDates,
	})
}

func PrepareConversationTokens(message wshandler.Message) ([]conversations.ConversationToken, []string, map[string][]caching.StoredMember, bool) {

	tokensUnparsed := message.Data["tokens"].([]interface{})
	tokens := make([]conversations.SentConversationToken, len(tokensUnparsed))
	for i, token := range tokensUnparsed {
		unparsed := token.(map[string]interface{})
		tokens[i] = conversations.SentConversationToken{
			ID:    unparsed["id"].(string),
			Token: unparsed["token"].(string),
		}
	}

	if len(tokens) > 500 {
		wshandler.ErrorResponse(message, "invalid")
		return nil, nil, nil, false
	}

	conversationTokens, err := caching.ValidateTokens(&tokens)
	if err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return nil, nil, nil, false
	}

	tokenIds := make([]string, len(conversationTokens))
	conversationIds := make([]string, len(conversationTokens))
	for i, token := range conversationTokens {
		tokenIds[i] = token.ID
		conversationIds[i] = token.Conversation
	}

	members, err := caching.LoadMembersArray(conversationIds)
	if err != nil {
		wshandler.ErrorResponse(message, "server.error")
		return nil, nil, nil, false
	}

	for id, token := range members {
		log.Printf("%s %d", id, len(token))
	}

	return conversationTokens, tokenIds, members, true
}
