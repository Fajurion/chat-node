package conversation

type Friend struct {
	ID   int64 `json:"id"`
	Node int64 `json:"node"`
}

// User ID -> Friends
var Friends map[int64][]int64 = make(map[int64][]int64)

func SetFriends(id int64, friends []int64) {
	Friends[id] = friends
}

func GetFriends(id int64) []int64 {
	return Friends[id]
}
