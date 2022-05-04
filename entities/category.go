package entities

type Category struct {
	ID    string `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required,min=4,max=20"`
	Image string `json:"image" binding:"url"`
}

func (c Category) String() string {
	return c.Name
}
