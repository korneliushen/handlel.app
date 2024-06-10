package model

// structs for kategorier
type Categories struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Name          string   `json:"name"`
	Store         string   `json:"store"`
	SubCategories []string `json:"subcategories"`
}