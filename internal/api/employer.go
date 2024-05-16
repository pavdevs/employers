package api

import (
	apidto "employer.dev/internal/api/dto"
	"employer.dev/internal/api/requests"
	"employer.dev/internal/api/responses"
	employerrepository "employer.dev/internal/repository/employer"
	repositorymodels "employer.dev/internal/repository/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type EmployerAPI struct {
	rep *employerrepository.EmployerRepository
}

func NewEmployerAPI(rep *employerrepository.EmployerRepository) *EmployerAPI {
	return &EmployerAPI{
		rep: rep,
	}
}

func (e *EmployerAPI) RegisterRoutes(r *gin.Engine) {
	employerGroup := r.Group("/employers")
	{
		employerGroup.POST("", e.addEmployer)
		employerGroup.DELETE("/:id", e.deleteEmployer)
		employerGroup.PATCH("/:id", e.updateEmployer)
		employerGroup.GET("/company/:companyId", e.getEmployersListForCompany)
		employerGroup.GET("/department/:name", e.getEmployersListForDepartment)
	}
}

// @Summary Create employer
// @Tags Employers
// @Description Данный роут создает сотрудника и возвращает его идентификатор в ответе
// @ID create_employer
// @Accept json
// @Produce json
// @Param {object} body requests.CreateEmployer true "Модель сотрудника"
// @Success 200 {object} responses.CreateEmployer "Сотрудник создан"
// @Router /employers [post]
func (e *EmployerAPI) addEmployer(c *gin.Context) {
	var req requests.CreateEmployer
	if err := c.BindJSON(&req); err != nil {

		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: err.Error(),
		})
		return
	}

	emp := repositorymodels.Employer{
		Name:            req.Name,
		Surname:         req.Surname,
		Phone:           req.Phone,
		CompanyID:       req.CompanyID,
		PassportType:    req.Passport.Type,
		PassportNumber:  req.Passport.Number,
		DepartmentName:  req.Department.Name,
		DepartmentPhone: req.Department.Phone,
	}

	id, err := e.rep.AddEmployer(emp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.CreateEmployer{
		ID: id,
	})
}

// @Summary Delete employer
// @Tags Employers
// @Description Данный удаляе сотрудника
// @ID delete_employer
// @Accept json
// @Produce json
// @Param id path int true "ID сотрудника"
// @Success 200 {object} responses.DeleteEmployer "Сотрудник удален"
// @Router /employers/{id} [delete]
func (e *EmployerAPI) deleteEmployer(c *gin.Context) {
	employerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employer ID"})
		return
	}

	if err = e.rep.DeleteEmployer(employerId); err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.DeleteEmployer{"Employer deleted"})
}

// @Summary Update employer by id
// @Tags Employers
// @Description Данный роут позволяет изменить данные сотрудника
// @ID update_employer_by_id
// @Accept json
// @Produce json
// @Param id path int true "Ид сторудника"
// @Param {object} body requests.UpdateEmployer true "Модель сотрудника"
// @Success 200 {object} responses.UpdateEmployer "Измененная мрдель сотрудника"
// @Router /employers/{id} [patch]
func (e *EmployerAPI) updateEmployer(c *gin.Context) {
	employerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employer ID"})
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emp := repositorymodels.NewUpdateEmployer()
	if name, ok := req["name"].(string); ok {
		emp.Name = &name
	} else {
		emp.Name = nil
	}
	if surname, ok := req["surname"].(string); ok {
		emp.Surname = &surname
	} else {
		emp.Surname = nil
	}
	if phone, ok := req["phone"].(string); ok {
		emp.Phone = &phone
	} else {
		emp.Phone = nil
	}
	if companyID, ok := req["companyId"].(int); ok {
		emp.CompanyID = &companyID
	} else {
		emp.CompanyID = nil
	}
	if passportType, ok := req["passportType"].(string); ok {
		emp.PassportType = &passportType
	} else {
		emp.PassportType = nil
	}
	if passportNumber, ok := req["passportNumber"].(string); ok {
		emp.PassportNumber = &passportNumber
	} else {
		emp.PassportNumber = nil
	}
	if departmentName, ok := req["departmentName"].(string); ok {
		emp.DepartmentName = &departmentName
	} else {
		emp.DepartmentName = nil
	}
	if departmentPhone, ok := req["departmentPhone"].(string); ok {
		emp.DepartmentPhone = &departmentPhone
	} else {
		emp.DepartmentPhone = nil
	}

	employer, err := e.rep.UpdateEmployer(emp, employerId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.UpdateEmployer{
		Employer: apidto.EmployerDTO{
			ID:         employer.ID,
			Name:       employer.Name,
			Surname:    employer.Surname,
			Phone:      employer.Phone,
			CompanyID:  employer.CompanyID,
			Passport:   apidto.PassportDTO{Type: employer.PassportType, Number: employer.PassportNumber},
			Department: apidto.DepartmentDTO{Name: employer.DepartmentName, Phone: employer.DepartmentPhone},
		},
	})
}

// @Summary Get employers by companyID
// @Tags Employers
// @Description Данный роут позволяет получить список сотрудников по ИД компании
// @ID get_employer_by_company_id
// @Accept json
// @Produce json
// @Param companyId path int true "ID компании"
// @Success 200 {object} responses.GetEmployers "Список сотрудников"
// @Router /employers/company/{companyId} [get]
func (e *EmployerAPI) getEmployersListForCompany(c *gin.Context) {
	companyId, err := strconv.Atoi(c.Param("companyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	employees, err := e.rep.GetEmployeesByCompany(companyId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: err.Error(),
		})
		return
	}

	var res []apidto.EmployerDTO

	for _, emp := range employees {
		res = append(res, apidto.EmployerDTO{
			ID:         emp.ID,
			Name:       emp.Name,
			Surname:    emp.Surname,
			Phone:      emp.Phone,
			CompanyID:  emp.CompanyID,
			Passport:   apidto.PassportDTO{emp.PassportType, emp.PassportNumber},
			Department: apidto.DepartmentDTO{emp.DepartmentName, emp.DepartmentPhone},
		})
	}

	c.JSON(http.StatusOK, responses.GetEmployers{
		Items: res,
	})
}

// @Summary Get employers by department name
// @Tags Employers
// @Description Данный роут позволяет получить список сотрудников по названию отдела
// @ID get_employer_by_department_name
// @Accept json
// @Produce json
// @Param name path string true "Название отдела"
// @Success 200 {object} responses.GetEmployers "Список сотрудников"
// @Router /employers/department/{name} [get]
func (e *EmployerAPI) getEmployersListForDepartment(c *gin.Context) {
	depName := c.Param("name")
	if depName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department name"})
		return
	}

	employees, err := e.rep.GetEmployeesByDepartment(depName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{
			Message: err.Error(),
		})
		return
	}

	var res []apidto.EmployerDTO

	for _, emp := range employees {
		res = append(res, apidto.EmployerDTO{
			ID:         emp.ID,
			Name:       emp.Name,
			Surname:    emp.Surname,
			Phone:      emp.Phone,
			CompanyID:  emp.CompanyID,
			Passport:   apidto.PassportDTO{emp.PassportType, emp.PassportNumber},
			Department: apidto.DepartmentDTO{emp.DepartmentName, emp.DepartmentPhone},
		})
	}

	c.JSON(http.StatusOK, responses.GetEmployers{
		Items: res,
	})
}
