package fetching

type Mailbox struct {
	ID    string `json:"id" gorm:"primaryKey"`  // Account ID
	Token string `json:"token" gorm:"not null"` // Mailbox token
}

type MailboxEntry struct {
	ID string `json:"id" gorm:"primaryKey"` // Entry ID

	Mailbox string `json:"mailbox" gorm:"not null"` // Mailbox ID
	Data    string `json:"data" gorm:"not null"`    // Encrypted data
}
