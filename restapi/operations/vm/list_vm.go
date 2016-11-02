package vm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ListVMHandlerFunc turns a function with the right signature into a list Vm handler
type ListVMHandlerFunc func(ListVMParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ListVMHandlerFunc) Handle(params ListVMParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ListVMHandler interface for that can handle valid list Vm params
type ListVMHandler interface {
	Handle(ListVMParams, interface{}) middleware.Responder
}

// NewListVM creates a new http.Handler for the list Vm operation
func NewListVM(ctx *middleware.Context, handler ListVMHandler) *ListVM {
	return &ListVM{Context: ctx, Handler: handler}
}

/*ListVM swagger:route GET /vms vm listVm

List vms of the pool

*/
type ListVM struct {
	Context *middleware.Context
	Handler ListVMHandler
}

func (o *ListVM) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewListVMParams()

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
