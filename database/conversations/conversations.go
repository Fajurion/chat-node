package conversations

type Conversation struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Account uint `json:"account" gorm:"not null"`
	Channel uint `json:"channel" gorm:"not null"`

	Members []Member `json:"members" gorm:"foreignKey:Conversation"`
}
