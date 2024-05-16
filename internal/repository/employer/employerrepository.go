package employerrepository

import (
	"database/sql"
	"employer.dev/internal/database"
	repositorymodels "employer.dev/internal/repository/models"
	"fmt"
	"github.com/sirupsen/logrus"
)

type EmployerRepository struct {
	database *database.Database
	logger   *logrus.Logger
}

func NewEmployerRepository(database *database.Database, logger *logrus.Logger) *EmployerRepository {
	return &EmployerRepository{
		database: database,
		logger:   logger,
	}
}

func (e *EmployerRepository) AddEmployer(emr repositorymodels.Employer) (int, error) {
	db, err := e.getPreparedDB()

	if err != nil {
		e.logger.Error(err)
		return -1, fmt.Errorf("can't connect to database: %w", err)
	}

	var empID int

	q := `INSERT INTO employers (name, surname, phone, company_id, passport_type, passport_number, department_name, department_phone) 
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err = db.QueryRow(
		q,
		emr.Name,
		emr.Surname,
		emr.Phone,
		emr.CompanyID,
		emr.PassportType,
		emr.PassportNumber,
		emr.DepartmentName,
		emr.DepartmentPhone,
	).Scan(&empID)

	if err != nil {
		e.logger.Error(err)
		return -1, fmt.Errorf("can't add employer: %w", err)
	}

	return empID, nil
}

func (e *EmployerRepository) DeleteEmployer(id int) error {
	db, err := e.getPreparedDB()

	if err != nil {
		e.logger.Error(err)
		return fmt.Errorf("can't connect to database: %w", err)
	}

	q := `DELETE FROM employers WHERE id = $1`
	_, err = db.Exec(q, id)

	if err != nil {
		e.logger.Error(err)
		return fmt.Errorf("can't delete employer: %w", err)
	}

	return nil
}

func (e *EmployerRepository) GetEmployeesByCompany(id int) ([]repositorymodels.Employer, error) {
	db, err := e.getPreparedDB()

	if err != nil {
		e.logger.Error(err)
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	q := `SELECT * FROM employers WHERE company_id = $1`

	rows, err := db.Query(q, id)

	if err != nil {
		e.logger.Error(err)
		return nil, fmt.Errorf("can't get employees: %w", err)
	}
	defer rows.Close()

	var emps []repositorymodels.Employer

	for rows.Next() {
		var emr repositorymodels.Employer

		err = rows.Scan(&emr.ID, &emr.Name, &emr.Surname, &emr.Phone, &emr.CompanyID, &emr.PassportType, &emr.PassportNumber, &emr.DepartmentName, &emr.DepartmentPhone)
		if err != nil {
			e.logger.Error(err)
			return nil, fmt.Errorf("can't get employees: %w", err)
		}

		emps = append(emps, emr)
	}

	return emps, nil
}

func (e *EmployerRepository) GetEmployeesByDepartment(name string) ([]repositorymodels.Employer, error) {
	db, err := e.getPreparedDB()

	if err != nil {
		e.logger.Error(err)
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	q := `SELECT * FROM employers WHERE department_name LIKE '%' || $1 || '%'`

	rows, err := db.Query(q, name)

	if err != nil {
		e.logger.Error(err)
		return nil, fmt.Errorf("can't get employees: %w", err)
	}
	defer rows.Close()

	var emps []repositorymodels.Employer

	for rows.Next() {
		var emr repositorymodels.Employer

		err = rows.Scan(&emr.ID, &emr.Name, &emr.Surname, &emr.Phone, &emr.CompanyID, &emr.PassportType, &emr.PassportNumber, &emr.DepartmentName, &emr.DepartmentPhone)
		if err != nil {
			e.logger.Error(err)
			return nil, fmt.Errorf("can't get employees: %w", err)
		}

		emps = append(emps, emr)
	}

	return emps, nil
}

func (e *EmployerRepository) UpdateEmployer(emr *repositorymodels.UpdateEmployer, id int) (*repositorymodels.Employer, error) {
	db, err := e.getPreparedDB()

	if err != nil {
		e.logger.Error(err)
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	q := `SELECT * FROM employers WHERE id = $1`

	var employer repositorymodels.Employer

	err = db.QueryRow(q, id).Scan(&employer.ID, &employer.Name, &employer.Surname, &employer.Phone, &employer.CompanyID, &employer.PassportType, &employer.PassportNumber, &employer.DepartmentName, &employer.DepartmentPhone)

	if err != nil {
		e.logger.Error(err)
		return nil, fmt.Errorf("can't get employer for update: %w", err)
	}

	if emr.Name != nil {
		employer.Name = *emr.Name
	}

	if emr.Surname != nil {
		employer.Surname = *emr.Surname
	}

	if emr.Phone != nil {
		employer.Phone = *emr.Phone
	}

	if emr.PassportType != nil {
		employer.PassportType = *emr.PassportType
	}

	if emr.PassportNumber != nil {
		employer.PassportNumber = *emr.PassportNumber
	}

	if emr.DepartmentName != nil {
		employer.DepartmentName = *emr.DepartmentName
	}

	if emr.DepartmentPhone != nil {
		employer.DepartmentPhone = *emr.DepartmentPhone
	}

	q = `UPDATE employers SET name = $1, surname = $2, phone = $3, passport_type = $4, passport_number = $5, department_name = $6, department_phone = $7 WHERE id = $8`
	_, err = db.Exec(q, employer.Name, employer.Surname, employer.Phone, employer.PassportType, employer.PassportNumber, employer.DepartmentName, employer.DepartmentPhone, employer.ID)

	if err != nil {
		e.logger.Error(err)
		return nil, fmt.Errorf("can't update employer: %w", err)
	}

	return &employer, nil
}

func (e *EmployerRepository) getPreparedDB() (*sql.DB, error) {
	db := e.database.GetDB()

	if err := db.Ping(); err != nil {
		e.logger.Error(err)
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return db, nil
}
