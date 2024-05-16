package responses

import apidto "employer.dev/internal/api/dto"

type UpdateEmployer struct {
	Employer apidto.EmployerDTO `json:"employer"`
}
