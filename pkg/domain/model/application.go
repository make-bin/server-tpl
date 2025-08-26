package model

// Application represents the application domain model
type Application struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

// TableName returns the table name for the Application model
func (a *Application) TableName() string {
	return "applications"
}

// ShortTableName returns abbreviated table name
func (a *Application) ShortTableName() string {
	return "app"
}

// Index returns indexable fields for the Application model
func (a *Application) Index() map[string]interface{} {
	index := a.BaseModel.Index()
	index["name"] = a.Name
	index["description"] = a.Description
	return index
}

// Validate performs business rule validation on the Application model
func (a *Application) Validate() error {
	if a.Name == "" {
		return ErrApplicationNameRequired
	}
	if len(a.Name) > 100 {
		return ErrApplicationNameTooLong
	}
	if len(a.Description) > 500 {
		return ErrApplicationDescriptionTooLong
	}
	return nil
}

// Domain errors for Application
var (
	ErrApplicationNameRequired       = NewDomainError("application name is required")
	ErrApplicationNameTooLong        = NewDomainError("application name too long")
	ErrApplicationDescriptionTooLong = NewDomainError("application description too long")
	ErrApplicationNotFound           = NewDomainError("application not found")
)

// DomainError represents domain-specific errors
type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new domain error
func NewDomainError(message string) *DomainError {
	return &DomainError{Message: message}
}
