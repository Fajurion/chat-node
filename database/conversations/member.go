package conversations

type Member struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Conversation uint `json:"conversation" gorm:"not null"`

	// 1 - member, 2 - admin, 3 - owner
	Role    uint `json:"role" gorm:"not null"`
	Account uint `json:"account" gorm:"not null"`
}

const RoleMember = 1
const RoleAdmin = 2
const RoleOwner = 3
