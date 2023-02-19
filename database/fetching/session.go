package fetching

type Session struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Token string `json:"session" gorm:"not null"`
	Fetch uint   `json:"fetch" gorm:"not null"` // Last time the session was used
}
