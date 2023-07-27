package conversations

type ConversationToken struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Conversation string `json:"conversation" gorm:"not null"` // Conversation id
	PubToken     string `json:"pub_token" gorm:"not null"`    // Short token required to publish in the conversation
	SubToken     string `json:"token" gorm:"not null"`        // Long token required to subscribe to the conversation
	Data         string `json:"payload" gorm:"not null"`      // Encrypted data about the user (account id, username, etc.)
	Rank         uint   `json:"rank" gorm:"not null"`
}

// * Permissions
const MinRankManageMembers = RankModerator
const MinRankChangeConversationDetails = RankModerator
const MinRankManageModerators = RankAdmin
const MinRankDeleteConversation = RankAdmin

// * Ranks
const RankUser = 0
const RankModerator = 1 // Can remove/add users
const RankAdmin = 2     // Manages moderators and can delete the conversation
