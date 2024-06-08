package entity

type OnlineEntity struct {
	Host   string   `json:"host"`
	Title  string   `json:"title"`
	Type   string   `json:"type"`
	Images []string `json:"images"`
}
