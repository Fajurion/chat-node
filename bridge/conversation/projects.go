package conversation

type Project struct {
	ID      int64           `json:"id"`
	Owner   int64           `json:"owner"`
	Data    string          `json:"name"`
	Members map[int64]int64 `json:"members"` // key = user id, value = node id
}

// Project ID -> Project
var Projects map[int64]Project = make(map[int64]Project)

func SetProject(id int64, project Project) {
	Projects[id] = project
}

func GetProject(id int64) (Project, error) {

	// Check if project exists
	if _, ok := Projects[id]; !ok {

		if err := FetchProject(id); err != nil {
			return Project{}, err
		}
	}

	return Projects[id], nil
}

func FetchProject(id int64) error {
	return nil
}
