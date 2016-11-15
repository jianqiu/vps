package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"

	"github.com/jianqiu/vps/restapi/operations"
	"github.com/jianqiu/vps/restapi/operations/vm"
	"github.com/jianqiu/vps/restapi/handlers"
	"github.com/jianqiu/vps/db"
	"github.com/jianqiu/vps/controllers"

	"code.cloudfoundry.org/lager"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../swagger.json

func configureFlags(api *operations.SoftLayerVMPoolAPI) {
}

func configureAPI(api *operations.SoftLayerVMPoolAPI,
logger lager.Logger,
db db.DB,
) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	vmController := controllers.NewVirtualGuestController(db)
	vmHandler := handlers.NewVmHandler(logger,vmController)

	api.VMAddVMHandler = vm.AddVMHandlerFunc(vmHandler.AddVM)
	api.VMDeleteVMHandler = vm.DeleteVMHandlerFunc(vmHandler.DeleteVM)
	api.VMGetVMByCidHandler = vm.GetVMByCidHandlerFunc(vmHandler.GetVMByCid)
	api.VMListVMHandler = vm.ListVMHandlerFunc(vmHandler.ListVM)
	api.VMUpdateVMWithStateHandler = vm.UpdateVMWithStateHandlerFunc(vmHandler.UpdateVMWithState)
	api.VMFindVmsByFiltersHandler = vm.FindVmsByFiltersHandlerFunc(vmHandler.FindVmsByFilters)
	api.VMFindVmsByDeploymentHandler = vm.FindVmsByDeploymentHandlerFunc(vmHandler.FindVmsByDeployment)
	api.VMFindVmsByStatesHandler = vm.FindVmsByStatesHandlerFunc(vmHandler.FindVmsByStates)
	api.VMUpdateVMHandler = vm.UpdateVMHandlerFunc(vmHandler.UpdateVM)
	api.VMOrderVMByFilterHandler = vm.OrderVMByFilterHandlerFunc(vmHandler.OrderVmByFilter)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
