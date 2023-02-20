package fetching

type Session struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Token     string `json:"token" gorm:"not null"`
	LastFetch int64  `json:"fetch" gorm:"type:bigint"` // Last time the session was used
}
