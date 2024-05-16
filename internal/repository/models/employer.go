package repositorymodels

type Employer struct {
	ID              int
	Name            string
	Surname         string
	Phone           string
	CompanyID       int
	PassportType    string
	PassportNumber  string
	DepartmentName  string
	DepartmentPhone string
}

type UpdateEmployer struct {
	ID              int
	Name            *string
	Surname         *string
	Phone           *string
	CompanyID       *int
	PassportType    *string
	PassportNumber  *string
	DepartmentName  *string
	DepartmentPhone *string
}

func NewUpdateEmployer() *UpdateEmployer {
	return &UpdateEmployer{}
}
