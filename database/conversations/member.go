package conversations

type Member struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Conversation uint `json:"conversation" gorm:"not null"`
	Role         uint `json:"role" gorm:"not null"`
	Account      uint `json:"account" gorm:"not null"`
}
