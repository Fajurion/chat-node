package fetching

type Session struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	Account   int64  `json:"account"`                  // Account ID
	Node      int64  `json:"node"`                     // Node ID
	LastFetch int64  `json:"fetch" gorm:"type:bigint"` // Last time the session was used
}
