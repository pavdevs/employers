package apidto

type EmployerDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Phone      string `json:"phone"`
	CompanyID  int    `json:"company_id"`
	Passport   PassportDTO
	Department DepartmentDTO
}
