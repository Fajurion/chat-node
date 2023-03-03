package fetching

type Status struct {
	ID int64 `gorm:"primaryKey"` // Account ID

	Status  string `gorm:"not null"`
	Updated int64  `gorm:"autoUpdateTime:milli"`
	Node    int64  `gorm:"not null"`
}
