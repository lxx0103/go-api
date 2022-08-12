package organization

type NewOrganizationCreated struct {
	OrganizationID string `json:"organization_id"`
	Owner          string `json:"owner"`
	OwnerEmail     string `json:"owner_email"`
	Password       string `json:"password"`
}
