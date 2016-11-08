package vm

import (
	"github.com/jianqiu/vps/models"
)

func (o *FindVmsByStatesOK) GetPayload() *models.VmsResponse{
	return o.Payload
}

func (o *FindVmsByStatesNotFound) GetStatusCode() int {
	return 404
}

func (o *FindVmsByStatesDefault) GetStatusCode() int {
	return o._statusCode
}

func (o *FindVmsByStatesDefault) GetPayload() *models.Error {
	return o.Payload
}