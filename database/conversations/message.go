package conversations

type Message struct {
	ID string `json:"id" gorm:"primaryKey"`

	Conversation uint   `json:"conversation" gorm:"not null"`
	Author       uint   `json:"author" gorm:"not null"`
	Creation     uint   `json:"creation" gorm:"not null"` // Unix timestamp (ms)
	Data         string `json:"data" gorm:"not null"`
}
