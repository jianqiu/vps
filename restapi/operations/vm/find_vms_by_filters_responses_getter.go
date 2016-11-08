package vm

import (
	"github.com/jianqiu/vps/models"
)

func (o *FindVmsByFiltersOK) GetPayload() *models.VmsResponse {
	return o.Payload
}

func (o *FindVmsByFiltersNotFound) GetStatusCode() int {
	return 404
}

func (o *FindVmsByFiltersDefault) GetStatusCode() int {
	return o._statusCode
}

func (o *FindVmsByFiltersDefault) GetPayload() *models.Error {
	return o.Payload
}
