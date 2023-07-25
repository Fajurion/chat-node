package conversations

type Conversation struct {
	ID string `json:"id" gorm:"primaryKey"`

	SubscriptionToken string `json:"token" gorm:"not null"`
	Data              string `json:"data" gorm:"not null"` // Encrypted with the conversation key
}
