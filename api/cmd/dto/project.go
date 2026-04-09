package dto

type PostProject struct {
	Name string `json:"name" validate:"required,max=255"`
}

type UpdateProject struct {
	ID           string  `param:"id"`
	ToBeArchived *bool   `json:"to_be_archived" validate:"omitempty"`
	Name         *string `json:"name" validate:"omitempty,max=255"`
}
