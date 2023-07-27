package conversations

type Conversation struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Data string `json:"data" gorm:"not null"` // Encrypted with the conversation key
}
