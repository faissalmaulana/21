package dto

type PostProject struct {
	Name string `json:"name" validate:"required,max=255"`
}
