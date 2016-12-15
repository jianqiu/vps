package vm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/jianqiu/vps/models"
)

// FindVmsByFiltersHandlerFunc turns a function with the right signature into a find vms by filters handler
type FindVmsByFiltersHandlerFunc func(FindVmsByFiltersParams, *models.User) middleware.Responder

// Handle executing the request and returning a response
func (fn FindVmsByFiltersHandlerFunc) Handle(params FindVmsByFiltersParams, principal *models.User) middleware.Responder {
	return fn(params, principal)
}

// FindVmsByFiltersHandler interface for that can handle valid find vms by filters params
type FindVmsByFiltersHandler interface {
	Handle(FindVmsByFiltersParams, *models.User) middleware.Responder
}

// NewFindVmsByFilters creates a new http.Handler for the find vms by filters operation
func NewFindVmsByFilters(ctx *middleware.Context, handler FindVmsByFiltersHandler) *FindVmsByFilters {
	return &FindVmsByFilters{Context: ctx, Handler: handler}
}

/*FindVmsByFilters swagger:route POST /vms/findByFilters vm findVmsByFilters

Finds Vms by filters (cpu, memory, private_vlan, public_vlan, state)

*/
type FindVmsByFilters struct {
	Context *middleware.Context
	Handler FindVmsByFiltersHandler
}

func (o *FindVmsByFilters) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewFindVmsByFiltersParams()

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
