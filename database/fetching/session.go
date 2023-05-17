package fetching

type Session struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Account   string `json:"account"`                  // Account ID
	Node      int64  `json:"node"`                     // Node ID
	LastFetch int64  `json:"fetch" gorm:"type:bigint"` // Last time the session was used
}
