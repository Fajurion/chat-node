package fetching

type Action struct {
	ID string `json:"-" gorm:"primaryKey"`

	Account   string `json:"-" gorm:"not null"`
	Action    string `json:"action" gorm:"not null"`
	Target    string `json:"target" gorm:"not null"`
	CreatedAt int64  `json:"created" gorm:"autoUpdateTime:milli"`
}
