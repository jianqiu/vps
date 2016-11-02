package vm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteVMHandlerFunc turns a function with the right signature into a delete Vm handler
type DeleteVMHandlerFunc func(DeleteVMParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteVMHandlerFunc) Handle(params DeleteVMParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteVMHandler interface for that can handle valid delete Vm params
type DeleteVMHandler interface {
	Handle(DeleteVMParams, interface{}) middleware.Responder
}

// NewDeleteVM creates a new http.Handler for the delete Vm operation
func NewDeleteVM(ctx *middleware.Context, handler DeleteVMHandler) *DeleteVM {
	return &DeleteVM{Context: ctx, Handler: handler}
}

/*DeleteVM swagger:route DELETE /vms/{cid} vm deleteVm

Deletes a vm from pool

*/
type DeleteVM struct {
	Context *middleware.Context
	Handler DeleteVMHandler
}

func (o *DeleteVM) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewDeleteVMParams()

	uprinc, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
