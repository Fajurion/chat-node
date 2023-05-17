package requests

import (
	"chat-node/database"
	"chat-node/database/conversations"
	"chat-node/util"
)

func LoadConversationDetails(id string) ([]string, []string, error) {

	// Get conversation
	var conversation conversations.Conversation
	if err := database.DBConn.Where(&conversations.Conversation{ID: id}).Take(&conversation).Error; err != nil {
		return nil, nil, err
	}

	// Get members and nodes
	var members []conversations.Member
	if err := database.DBConn.Where(&conversations.Member{Conversation: id}).Preload("Status").Find(&members).Error; err != nil {
		return nil, nil, err
	}

	// Turn into arrays
	var accounts []string
	var nodes []string
	for _, member := range members {
		accounts = append(accounts, member.Status.ID)
		nodes = append(nodes, util.Node64(member.Status.Node))
	}

	return accounts, nodes, nil
}
