package organization

type NewOrganizationCreated struct {
	OrganizationID string `json:"organization_id"`
	Owner          string `json:"owner"`
	Password       string `json:"password"`
}
