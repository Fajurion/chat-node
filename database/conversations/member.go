package conversations

import "chat-node/database/fetching"

type Member struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Conversation string `json:"conversation" gorm:"not null"`

	// 1 - member, 2 - admin, 3 - owner
	Role    uint   `json:"role" gorm:"not null"`
	Account string `json:"account" gorm:"not null"`

	// Relationships
	Status fetching.Status `json:"-" gorm:"foreignKey:Account"`
}

const RoleMember = 1
const RoleAdmin = 2
const RoleOwner = 3
