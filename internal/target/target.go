package target

type Target struct {
	BaseURL  string
	Content  []byte
	Priority float32
}

func NewTarget(baseURL string) *Target {
	return &Target{BaseURL: baseURL}
}
