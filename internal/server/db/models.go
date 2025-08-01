// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type AccountType string

const (
	AccountTypeOauth       AccountType = "oauth"
	AccountTypeEmail       AccountType = "email"
	AccountTypeCredentials AccountType = "credentials"
)

func (e *AccountType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccountType(s)
	case string:
		*e = AccountType(s)
	default:
		return fmt.Errorf("unsupported scan type for AccountType: %T", src)
	}
	return nil
}

type NullAccountType struct {
	AccountType AccountType
	Valid       bool // Valid is true if AccountType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAccountType) Scan(value interface{}) error {
	if value == nil {
		ns.AccountType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AccountType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAccountType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AccountType), nil
}

type ApplicationStatus string

const (
	ApplicationStatusPending  ApplicationStatus = "pending"
	ApplicationStatusRejected ApplicationStatus = "rejected"
	ApplicationStatusAccepted ApplicationStatus = "accepted"
	ApplicationStatusTraining ApplicationStatus = "training"
	ApplicationStatusHired    ApplicationStatus = "hired"
)

func (e *ApplicationStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ApplicationStatus(s)
	case string:
		*e = ApplicationStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ApplicationStatus: %T", src)
	}
	return nil
}

type NullApplicationStatus struct {
	ApplicationStatus ApplicationStatus
	Valid             bool // Valid is true if ApplicationStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullApplicationStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ApplicationStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ApplicationStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullApplicationStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ApplicationStatus), nil
}

type JobStatus string

const (
	JobStatusUpcoming  JobStatus = "upcoming"
	JobStatusOpen      JobStatus = "open"
	JobStatusClosed    JobStatus = "closed"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusCompleted JobStatus = "completed"
)

func (e *JobStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = JobStatus(s)
	case string:
		*e = JobStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for JobStatus: %T", src)
	}
	return nil
}

type NullJobStatus struct {
	JobStatus JobStatus
	Valid     bool // Valid is true if JobStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullJobStatus) Scan(value interface{}) error {
	if value == nil {
		ns.JobStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.JobStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullJobStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.JobStatus), nil
}

type JobType string

const (
	JobTypeFullTime   JobType = "fullTime"
	JobTypePartTime   JobType = "partTime"
	JobTypeContract   JobType = "contract"
	JobTypeInternship JobType = "internship"
	JobTypeFreelance  JobType = "freelance"
	JobTypeTemporary  JobType = "temporary"
	JobTypeVolunteer  JobType = "volunteer"
	JobTypeRemote     JobType = "remote"
	JobTypeOnSite     JobType = "onSite"
	JobTypeHybrid     JobType = "hybrid"
)

func (e *JobType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = JobType(s)
	case string:
		*e = JobType(s)
	default:
		return fmt.Errorf("unsupported scan type for JobType: %T", src)
	}
	return nil
}

type NullJobType struct {
	JobType JobType
	Valid   bool // Valid is true if JobType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullJobType) Scan(value interface{}) error {
	if value == nil {
		ns.JobType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.JobType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullJobType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.JobType), nil
}

type UploadType string

const (
	UploadTypeResume         UploadType = "resume"
	UploadTypeJobDescription UploadType = "jobDescription"
	UploadTypeRequirements   UploadType = "requirements"
)

func (e *UploadType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UploadType(s)
	case string:
		*e = UploadType(s)
	default:
		return fmt.Errorf("unsupported scan type for UploadType: %T", src)
	}
	return nil
}

type NullUploadType struct {
	UploadType UploadType
	Valid      bool // Valid is true if UploadType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUploadType) Scan(value interface{}) error {
	if value == nil {
		ns.UploadType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UploadType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUploadType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UploadType), nil
}

type Account struct {
	ID                pgtype.UUID
	UserID            pgtype.UUID
	Type              AccountType
	Provider          string
	ProviderAccountID string
	AccessToken       pgtype.Text
	RefreshToken      pgtype.Text
	ExpiresAt         pgtype.Int4
	TokenType         pgtype.Text
	IDToken           pgtype.Text
	SessionState      pgtype.Text
	Scope             pgtype.Text
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
}

type Application struct {
	ID           pgtype.UUID
	Status       ApplicationStatus
	JobOpeningID pgtype.UUID
	UserID       pgtype.UUID
	ParsedResume pgtype.Text
	LayoutScore  float64
	ContentScore float64
	SkillGap     pgtype.Text
	ResumeID     pgtype.UUID
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
}

type Assessment struct {
	ID             pgtype.UUID
	Title          string
	Description    pgtype.Text
	Questions      []byte
	LearningPlanID pgtype.UUID
	CreatedAt      pgtype.Timestamp
	UpdatedAt      pgtype.Timestamp
}

type JobOpening struct {
	ID                 pgtype.UUID
	Title              string
	Company            string
	Location           pgtype.Text
	Type               JobType
	Description        string
	Contact            string
	Address            pgtype.Text
	Status             JobStatus
	Deadline           pgtype.Timestamp
	StartDate          pgtype.Timestamp
	EndDate            pgtype.Timestamp
	RequirementsFileID pgtype.UUID
	ParsedRequirements pgtype.Text
	TrainingID         pgtype.UUID
	RecruiterID        pgtype.UUID
	CreatedAt          pgtype.Timestamp
	UpdatedAt          pgtype.Timestamp
}

type LearningPlan struct {
	ID            pgtype.UUID
	PlanDetails   []byte
	ApplicationID pgtype.UUID
	TrainingID    pgtype.UUID
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
}

type Recruiter struct {
	ID           pgtype.UUID
	Name         string
	Position     string
	Organization string
	Phone        pgtype.Text
	Address      pgtype.Text
	UserID       pgtype.UUID
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
}

type Session struct {
	ID           pgtype.UUID
	SessionToken string
	RefreshToken string
	ExpiresAt    pgtype.Timestamp
	UserID       pgtype.UUID
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
}

type Training struct {
	ID        pgtype.UUID
	Topics    string
	StartDate pgtype.Timestamp
	EndDate   pgtype.Timestamp
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type Upload struct {
	ID         pgtype.UUID
	UploadType UploadType
	Name       string
	FileType   string
	Url        string
	UserID     pgtype.UUID
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

type User struct {
	ID            pgtype.UUID
	Name          pgtype.Text
	Email         string
	Image         string
	EmailVerified pgtype.Timestamp
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
}

type VerificationToken struct {
	Identifier string
	Token      string
	ExpiresAt  pgtype.Timestamp
}
