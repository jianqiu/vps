package vm

import (
	"github.com/jianqiu/vps/models"
)

func (o *UpdateVMWithStateOK) GetPayload() string {
	return o.Payload
}

func (o *UpdateVMWithStateNotFound) GetStatusCode() int {
	return 404
}

func (o *UpdateVMWithStateDefault) GetStatusCode() int {
	return o._statusCode
}

func (o *UpdateVMWithStateDefault) GetPayload() *models.Error{
	return o.Payload
}