package handlers

import (
	"github.com/jianqiu/vps/models"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jianqiu/vps/restapi/operations/vm"

	"code.cloudfoundry.org/lager"
)
//go:generate counterfeiter -o fake_controllers/fake_virtualguest_controller.go . VirtualGuestController
type VirtualGuestController interface {
	AllVirtualGuests(logger lager.Logger) ([]*models.VM, error)
	VirtualGuests(logger lager.Logger, publicVlan, privateVlan, cpu, memory_mb int32, state models.State) ([]*models.VM, error)
	OrderVirtualGuest(logger lager.Logger, vmFilter *models.VMFilter) (*models.VM, error)
	VirtualGuestsByDeployments(logger lager.Logger, names []string) ([]*models.VM, error)
	VirtualGuestsByStates(logger lager.Logger, states []string) ([]*models.VM, error)
	CreateVM(logger lager.Logger, vm *models.VM) error
	DeleteVM(logger lager.Logger, cid int32) error
	UpdateVM(logger lager.Logger, vm *models.VM) error
	UpdateVMWithState(logger lager.Logger, cid int32, updateData *models.State) error
	VirtualGuestByCid(logger lager.Logger, cid int32) (*models.VM, error)
}

func NewVmHandler(
logger lager.Logger,
controller VirtualGuestController,
) *VMHandler {
	return &VMHandler{
		logger: logger,
		controller: controller,
	}
}

type VMHandler struct {
	logger lager.Logger
	controller VirtualGuestController
}

func (h *VMHandler) AddVM (params vm.AddVMParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("add-vm")

	request := params.Body

	err = h.controller.CreateVM(h.logger, request)
	if err != nil {
		unExpectedResponse := vm.NewAddVMDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}
	return vm.NewAddVMOK().WithPayload("added successfully")
}

func (h *VMHandler) OrderVmByFilter(params vm.OrderVMByFilterParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("order-vm-by-filter")

	response := &models.VMResponse{}
	request := params.Body

	response.VM, err = h.controller.OrderVirtualGuest(h.logger, request)
	if err != nil {
		unExpectedResponse := vm.NewOrderVMByFilterDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}

	return vm.NewOrderVMByFilterOK().WithPayload(response)
}

func (h *VMHandler) UpdateVM (params vm.UpdateVMParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("update-vm")

	request := params.Body

	err = h.controller.UpdateVM(h.logger, request)
	if err != nil {
		unExpectedResponse := vm.NewUpdateVMDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}
	return vm.NewUpdateVMOK().WithPayload("updated successfully")
}

func (h *VMHandler) DeleteVM(params vm.DeleteVMParams)  middleware.Responder {
	var err error
	h.logger = h.logger.Session("delete-vm")

	err = h.controller.DeleteVM(h.logger, params.Cid)
	if err != nil {
		unExpectedResponse := vm.NewDeleteVMDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}

	return vm.NewDeleteVMNoContent().WithPayload("vm removed")
}


func (h *VMHandler) GetVMByCid(params vm.GetVMByCidParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("get-vm-by-cid")

	response := &models.VMResponse{}

	response.VM, err = h.controller.VirtualGuestByCid(h.logger,params.Cid)
	if err != nil {
		unExpectedResponse := vm.NewGetVMByCidDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}
	if response.VM == nil {
		getVMByCidNotFound := vm.NewGetVMByCidNotFound()
		return getVMByCidNotFound
	}

	return vm.NewGetVMByCidOK().WithPayload(response)
}

func (h *VMHandler) ListVM(params vm.ListVMParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("list-vms")

	response := &models.VmsResponse{}


	response.Vms, err = h.controller.AllVirtualGuests(h.logger)
	if err != nil {
		unExpectedResponse := vm.NewListVMDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}
	if len(response.Vms) == 0 {
		return vm.NewListVMNotFound()
	}

	return vm.NewListVMOK().WithPayload(response)
}

func (h *VMHandler) UpdateVMWithState(params vm.UpdateVMWithStateParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("update-vm-with-state")

	vmId := params.Cid
	if vmId == 0 {
		 return vm.NewUpdateVMWithStateNotFound()
	}

	updateData := params.Body
	err = h.controller.UpdateVMWithState(h.logger, vmId, &updateData.State)
	if err != nil {
		unExpectedResponse := vm.NewListVMDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}

	return vm.NewUpdateVMOK().WithPayload("updated successfully")
}

func (h *VMHandler) FindVmsByFilters(params vm.FindVmsByFiltersParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("find-vms-by-filter")

	response := &models.VmsResponse{}
	request := params.Body

	if request == nil {
		response.Vms, err = h.controller.AllVirtualGuests(h.logger)
	} else {
		response.Vms, err = h.controller.VirtualGuests(h.logger, request.PublicVlan, request.PrivateVlan, request.CPU, request.MemoryMb, request.State)
	}
	if err != nil {
		unExpectedResponse := vm.NewFindVmsByFiltersDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}

	return vm.NewFindVmsByFiltersOK().WithPayload(response)
}

func (h *VMHandler) FindVmsByDeployment(params vm.FindVmsByDeploymentParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("find-vms-by-deployment")

	response := &models.VmsResponse{}
	request := params.Deployment

	response.Vms, err = h.controller.VirtualGuestsByDeployments(h.logger, request)
	if err != nil {
		unExpectedResponse := vm.NewFindVmsByDeploymentDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}

	return vm.NewFindVmsByDeploymentOK().WithPayload(response)
}

func (h *VMHandler) FindVmsByStates(params vm.FindVmsByStatesParams) middleware.Responder {
	var err error
	h.logger = h.logger.Session("find-vms-by-state")

	response := &models.VmsResponse{}
	request := params.States

	response.Vms, err = h.controller.VirtualGuestsByStates(h.logger, request)
	if err != nil {
		unExpectedResponse := vm.NewFindVmsByStatesDefault(500)
		unExpectedResponse.SetPayload(models.ConvertError(err))
		return unExpectedResponse
	}

	return vm.NewFindVmsByStatesOK().WithPayload(response)
}