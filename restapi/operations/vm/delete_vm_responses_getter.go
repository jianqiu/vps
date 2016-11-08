package vm

import (
	"github.com/jianqiu/vps/models"
)

func (o *DeleteVMNoContent) GetPayload() string {
	return o.Payload
}

func (o *DeleteVMNotFound) GetStatusCode() int {
	return 404
}

func (o *DeleteVMDefault) GetStatusCode() int {
	return o._statusCode
}

func (o *DeleteVMDefault) GetPayload() *models.Error{
	return o.Payload
}