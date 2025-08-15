package v1

import (
	v1 "github.com/make-bin/server-tpl/pkg/api/dto/v1"
	"github.com/make-bin/server-tpl/pkg/domain/model"
)

// ToApplicationModel 将DTO转换为模型
func ToApplicationModel(req interface{}) *model.Application {
	switch r := req.(type) {
	case *v1.CreateApplicationRequest:
		return &model.Application{
			Name:        r.Name,
			Description: r.Description,
			Version:     r.Version,
			CreatedBy:   r.CreatedBy,
		}
	case *v1.UpdateApplicationRequest:
		return &model.Application{
			Name:        r.Name,
			Description: r.Description,
			Version:     r.Version,
			Status:      r.Status,
		}
	default:
		return &model.Application{}
	}
}

// ToApplicationResponse 将模型转换为响应DTO
func ToApplicationResponse(app *model.Application) *v1.ApplicationResponse {
	if app == nil {
		return nil
	}
	return &v1.ApplicationResponse{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		Version:     app.Version,
		Status:      app.Status,
		CreatedBy:   app.CreatedBy,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}
}

// ToApplicationListResponse 将模型列表转换为响应DTO
func ToApplicationListResponse(apps []*model.Application) *v1.ApplicationListResponse {
	responses := make([]v1.ApplicationResponse, 0, len(apps))
	for _, app := range apps {
		if response := ToApplicationResponse(app); response != nil {
			responses = append(responses, *response)
		}
	}
	return &v1.ApplicationListResponse{
		Applications: responses,
		Total:        int64(len(responses)),
	}
}
