package fetching

type Status struct {
	ID string `gorm:"primaryKey"` // Account ID

	Type    uint   `gorm:"not null"` // 0 = Online, 1 = Offline, 2 = Away, 3 = Do Not Disturb
	Status  string `gorm:"not null"`
	Updated int64  `gorm:"autoUpdateTime:milli"`
	Node    int64  `gorm:"not null"`
}

const StatusOnline = 0
const StatusOffline = 1
const StatusAway = 2
const StatusDoNotDisturb = 3
