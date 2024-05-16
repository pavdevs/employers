package requests

type UpdateEmployer struct {
	Name            string `json:"name,omitempty"`
	Surname         string `json:"surname,omitempty"`
	Phone           string `json:"phone,omitempty"`
	CompanyID       int    `json:"company_id,omitempty"`
	PassportType    string `json:"passport_type,omitempty"`
	PassportNumber  string `json:"passport_number,omitempty"`
	DepartmentName  string `json:"department_name,omitempty"`
	DepartmentPhone string `json:"department_phone,omitempty"`
}
