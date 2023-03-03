package requests

import (
	"chat-node/database"
	"chat-node/database/conversations"
)

func LoadConversationDetails(id uint) ([]int64, []int64, error) {

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
	var accounts []int64
	var nodes []int64
	for _, member := range members {
		accounts = append(accounts, member.Status.ID)
		nodes = append(nodes, member.Status.Node)
	}

	return accounts, nodes, nil
}
