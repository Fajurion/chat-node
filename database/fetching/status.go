package fetching

type Status struct {
	ID string `gorm:"primaryKey"` // Account ID

	Type    uint   `gorm:"not null"` // 0 = Offline, 1 = Online, 2 = Away, 3 = Do Not Disturb
	Status  string `gorm:"not null"`
	Updated int64  `gorm:"autoUpdateTime:milli"`
	Node    int64  `gorm:"not null"`
}

const StatusOffline = 0
const StatusOnline = 1
const StatusAway = 2
const StatusDoNotDisturb = 3
