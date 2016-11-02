package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/jianqiu/vps/restapi/operations"
	"github.com/jianqiu/vps/restapi/operations/user"
	"github.com/jianqiu/vps/restapi/operations/vm"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../swagger.json

func configureFlags(api *operations.SoftLayerVMPoolAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SoftLayerVMPoolAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.VMAddVMHandler = vm.AddVMHandlerFunc(func(params vm.AddVMParams) middleware.Responder {
		return middleware.NotImplemented("operation vm.AddVM has not yet been implemented")
	})
	api.UserCreateUserHandler = user.CreateUserHandlerFunc(func(params user.CreateUserParams) middleware.Responder {
		return middleware.NotImplemented("operation user.CreateUser has not yet been implemented")
	})
	api.UserCreateUsersWithArrayInputHandler = user.CreateUsersWithArrayInputHandlerFunc(func(params user.CreateUsersWithArrayInputParams) middleware.Responder {
		return middleware.NotImplemented("operation user.CreateUsersWithArrayInput has not yet been implemented")
	})
	api.UserCreateUsersWithListInputHandler = user.CreateUsersWithListInputHandlerFunc(func(params user.CreateUsersWithListInputParams) middleware.Responder {
		return middleware.NotImplemented("operation user.CreateUsersWithListInput has not yet been implemented")
	})
	api.UserDeleteUserHandler = user.DeleteUserHandlerFunc(func(params user.DeleteUserParams) middleware.Responder {
		return middleware.NotImplemented("operation user.DeleteUser has not yet been implemented")
	})
	api.VMDeleteVMHandler = vm.DeleteVMHandlerFunc(func(params vm.DeleteVMParams) middleware.Responder {
		return middleware.NotImplemented("operation vm.DeleteVM has not yet been implemented")
	})
	api.VMFindVmsByDeploymentHandler = vm.FindVmsByDeploymentHandlerFunc(func(params vm.FindVmsByDeploymentParams) middleware.Responder {
		return middleware.NotImplemented("operation vm.FindVmsByDeployment has not yet been implemented")
	})
	api.VMFindVmsByStatesHandler = vm.FindVmsByStatesHandlerFunc(func(params vm.FindVmsByStatesParams) middleware.Responder {
		return middleware.NotImplemented("operation vm.FindVmsByStates has not yet been implemented")
	})
	api.UserGetUserByNameHandler = user.GetUserByNameHandlerFunc(func(params user.GetUserByNameParams) middleware.Responder {
		return middleware.NotImplemented("operation user.GetUserByName has not yet been implemented")
	})
	api.VMGetVMByCidHandler = vm.GetVMByCidHandlerFunc(func(params vm.GetVMByCidParams) middleware.Responder {
		return middleware.NotImplemented("operation vm.GetVMByCid has not yet been implemented")
	})
	api.VMListVMHandler = vm.ListVMHandlerFunc(func(params vm.ListVMParams) middleware.Responder {
		return middleware.NotImplemented("operation vm.ListVM has not yet been implemented")
	})
	api.UserLoginUserHandler = user.LoginUserHandlerFunc(func(params user.LoginUserParams) middleware.Responder {
		return middleware.NotImplemented("operation user.LoginUser has not yet been implemented")
	})
	api.UserLogoutUserHandler = user.LogoutUserHandlerFunc(func(params user.LogoutUserParams) middleware.Responder {
		return middleware.NotImplemented("operation user.LogoutUser has not yet been implemented")
	})
	api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(func(params user.UpdateUserParams) middleware.Responder {
		return middleware.NotImplemented("operation user.UpdateUser has not yet been implemented")
	})
	api.VMUpdateVMHandler = vm.UpdateVMHandlerFunc(func(params vm.UpdateVMParams) middleware.Responder {
		return middleware.NotImplemented("operation vm.UpdateVM has not yet been implemented")
	})
	api.VMUpdateVMWithStateHandler = vm.UpdateVMWithStateHandlerFunc(func(params vm.UpdateVMWithStateParams) middleware.Responder {
		return middleware.NotImplemented("operation vm.UpdateVMWithState has not yet been implemented")
	})

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
