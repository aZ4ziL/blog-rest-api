package utils

type CategoryPayload struct {
	Title    string `form:"title" validate:"required"`
	Slug     string `form:"slug" validate:"required"`
	Approved bool   `form:"approved" validate:"required"`
}

type CategoryQuerySlug struct {
	Slug string `form:"slug"`
}
