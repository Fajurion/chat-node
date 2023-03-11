package fetching

type Action struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account   uint   `json:"account" gorm:"not null"`
	Action    string `json:"action" gorm:"not null"`
	Target    string `json:"target" gorm:"not null"`
	CreatedAt int64  `json:"created" gorm:"autoUpdateTime:milli"`
}
