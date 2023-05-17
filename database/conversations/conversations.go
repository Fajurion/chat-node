package conversations

type Conversation struct {
	ID string `json:"id" gorm:"primaryKey"`

	Type      string    `json:"type" gorm:"not null"` // chat or space
	Data      string    `json:"data" gorm:"not null"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime:milli"`
	Creator   string    `json:"creator" gorm:"not null"`
	Members   []Member  `json:"-" gorm:"foreignKey:Conversation"`
	Messages  []Message `json:"-" gorm:"foreignKey:Conversation"`
}
