package v1

import (
	v1 "github.com/make-bin/server-tpl/pkg/api/dto/v1"
	"github.com/make-bin/server-tpl/pkg/domain/model"
)

// ToVariableModel 将DTO转换为模型
func ToVariableModel(req interface{}) *model.Variable {
	switch r := req.(type) {
	case *v1.CreateVariableRequest:
		return &model.Variable{
			ApplicationID: r.ApplicationID,
			Key:           r.Key,
			Value:         r.Value,
			Description:   r.Description,
			Type:          r.Type,
			IsSecret:      r.IsSecret,
		}
	case *v1.UpdateVariableRequest:
		return &model.Variable{
			ApplicationID: r.ApplicationID,
			Key:           r.Key,
			Value:         r.Value,
			Description:   r.Description,
			Type:          r.Type,
			IsSecret:      r.IsSecret,
		}
	default:
		return &model.Variable{}
	}
}

// ToVariableResponse 将模型转换为响应DTO
func ToVariableResponse(variable *model.Variable) *v1.VariableResponse {
	if variable == nil {
		return nil
	}
	return &v1.VariableResponse{
		ID:            variable.ID,
		ApplicationID: variable.ApplicationID,
		Key:           variable.Key,
		Value:         variable.Value,
		Description:   variable.Description,
		Type:          variable.Type,
		IsSecret:      variable.IsSecret,
		CreatedAt:     variable.CreatedAt,
		UpdatedAt:     variable.UpdatedAt,
	}
}

// ToVariableListResponse 将模型列表转换为响应DTO
func ToVariableListResponse(variables []*model.Variable) *v1.VariableListResponse {
	responses := make([]v1.VariableResponse, 0, len(variables))
	for _, variable := range variables {
		if response := ToVariableResponse(variable); response != nil {
			responses = append(responses, *response)
		}
	}
	return &v1.VariableListResponse{
		Variables: responses,
		Total:     int64(len(responses)),
	}
}
