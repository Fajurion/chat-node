package conversations

type ConversationToken struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Token string `json:"token" gorm:"not null"` // Long token required to subscribe to the conversation
	Rank  uint   `json:"rank" gorm:"not null"`
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
