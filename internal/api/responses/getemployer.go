package responses

import apidto "employer.dev/internal/api/dto"

type GetEmployers struct {
	Items []apidto.EmployerDTO `json:"employers"`
}
