package conversations

type Conversation struct {
	ID uint `json:"id" gorm:"primaryKey"`

	CreatedAt int64     `json:"created_at" gorm:"autoCreateTime:milli"`
	Creator   int64     `json:"creator" gorm:"not null"`
	Members   []Member  `json:"members" gorm:"foreignKey:Conversation"`
	Messages  []Message `json:"messages" gorm:"foreignKey:Conversation"`
}
