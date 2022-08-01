package organization

type NewOrganizationCreated struct {
	OrganizationID int64  `json:"organization_id"`
	Owner          string `json:"owner"`
	Password       string `json:"password"`
}
