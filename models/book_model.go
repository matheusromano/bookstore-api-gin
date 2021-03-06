package models

type Book struct {
	//Id     primitive.ObjectID `json:"id,omitempty"`
	Title  string  `json:"title,omitempty" validate:"required"`
	Author string  `json:"author,omitempty" validate:"required"`
	Price  float64 `json:"price,omitempty" validate:"required"`
}
