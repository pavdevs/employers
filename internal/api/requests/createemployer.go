package requests

import apidto "employer.dev/internal/api/dto"

type CreateEmployer struct {
	Name       string               `json:"name" binding:"required"`
	Surname    string               `json:"surname" binding:"required"`
	Phone      string               `json:"phone" binding:"required"`
	CompanyID  int                  `json:"company_id" binding:"required"`
	Passport   apidto.PassportDTO   `json:"passport" binding:"required"`
	Department apidto.DepartmentDTO `json:"department" binding:"required"`
}
