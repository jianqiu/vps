package vm

import (
	"github.com/jianqiu/vps/models"
)

func (o *ListVMOK) GetPayload() *models.VmsResponse {
	return o.Payload
}

func (o *ListVMNotFound) GetStatusCode() int {
	return 404
}

func (o *ListVMDefault) GetStatusCode() int{
	return o._statusCode
}

func (o *ListVMDefault) GetPayload() *models.Error {
	return o.Payload
}


