package conversations

type Conversation struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Data      string    `json:"data" gorm:"not null"`
	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime:milli"`
	Creator   int64     `json:"creator" gorm:"not null"`
	Members   []Member  `json:"-" gorm:"foreignKey:Conversation"`
	Messages  []Message `json:"-" gorm:"foreignKey:Conversation"`
}
