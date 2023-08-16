package fetching

type Subscription struct {
	ID   string `gorm:"primaryKey"` // Subscription token
	Node int64  `gorm:"not null"`
}
