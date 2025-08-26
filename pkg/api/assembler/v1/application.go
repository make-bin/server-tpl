package v1

import (
	dto "github.com/make-bin/server-tpl/pkg/api/dto/v1"
	"github.com/make-bin/server-tpl/pkg/domain/model"
)

// ApplicationAssembler handles conversion between domain models and DTOs
type ApplicationAssembler struct{}

// NewApplicationAssembler creates a new ApplicationAssembler instance
func NewApplicationAssembler() *ApplicationAssembler {
	return &ApplicationAssembler{}
}

// ToModel converts ApplicationRequest DTO to domain model
func (a *ApplicationAssembler) ToModel(req *dto.ApplicationRequest) *model.Application {
	return &model.Application{
		Name:        req.Name,
		Description: req.Description,
	}
}

// ToResponse converts domain model to ApplicationResponse DTO
func (a *ApplicationAssembler) ToResponse(app *model.Application) *dto.ApplicationResponse {
	return &dto.ApplicationResponse{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}
}

// ToResponseList converts slice of domain models to ApplicationListResponse DTO
func (a *ApplicationAssembler) ToResponseList(apps []*model.Application, total int64, page, pageSize int) *dto.ApplicationListResponse {
	responses := make([]dto.ApplicationResponse, len(apps))
	for i, app := range apps {
		responses[i] = *a.ToResponse(app)
	}

	return &dto.ApplicationListResponse{
		Applications: responses,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
	}
}
