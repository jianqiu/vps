package vm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/jianqiu/vps/models"
)

// GetVMByCidHandlerFunc turns a function with the right signature into a get Vm by cid handler
type GetVMByCidHandlerFunc func(GetVMByCidParams, *models.User) middleware.Responder

// Handle executing the request and returning a response
func (fn GetVMByCidHandlerFunc) Handle(params GetVMByCidParams, principal *models.User) middleware.Responder {
	return fn(params, principal)
}

// GetVMByCidHandler interface for that can handle valid get Vm by cid params
type GetVMByCidHandler interface {
	Handle(GetVMByCidParams, *models.User) middleware.Responder
}

// NewGetVMByCid creates a new http.Handler for the get Vm by cid operation
func NewGetVMByCid(ctx *middleware.Context, handler GetVMByCidHandler) *GetVMByCid {
	return &GetVMByCid{Context: ctx, Handler: handler}
}

/*GetVMByCid swagger:route GET /vms/{cid} vm getVmByCid

Find vm by ID

Returns a vm when ID < 10.  ID > 10 or nonintegers will simulate API error conditions

*/
type GetVMByCid struct {
	Context *middleware.Context
	Handler GetVMByCidHandler
}

func (o *GetVMByCid) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetVMByCidParams()

	uprinc, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	var principal *models.User
	if uprinc != nil {
		principal = uprinc.(*models.User) // this is really a models.User, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
