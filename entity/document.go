package entity

// Document represents a document entity
type Document struct {
	ID            int      `json:"id" bson:"id" validate:"required"`
	Title         string   `json:"title" bson:"title"`
	OriginalTitle string   `json:"original_title" bson:"original_title"`
	Overview      string   `json:"overview" bson:"overview"`
	CreationDate  string   `json:"creation_date" bson:"creation_date"`
	Keywords      []string `json:"keywords" bson:"keywords"`
}
