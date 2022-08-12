package organization

type OrganizationFilter struct {
	Name     string `form:"name" binding:"omitempty,max=64,min=1"`
	PageId   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=5,max=200"`
}

type OrganizationNew struct {
	Name     string `json:"name" binding:"required,min=1,max=64"`
	UserName string `json:"user_name" binding:"required,min=1,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required,min=8"`
}

type OrganizationID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
