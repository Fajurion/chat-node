package fetching

type Session struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	LastFetch int64  `json:"fetch" gorm:"type:bigint"` // Last time the session was used
}
