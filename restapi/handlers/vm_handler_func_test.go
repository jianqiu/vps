package handlers_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/jianqiu/vps/restapi/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jianqiu/vps/restapi/handlers"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/jianqiu/vps/restapi/handlers/fake_controllers"
	"github.com/jianqiu/vps/models"
	"github.com/jianqiu/vps/restapi/operations/vm"
	"github.com/go-openapi/runtime/middleware"
)

var _ = Describe("VmHandlerFunc", func() {
	var (
		logger     *lagertest.TestLogger
		controller *fake_controllers.FakeVirtualGuestController

		responseResponder middleware.Responder

		handler *handlers.VMHandler
		exitCh  chan struct{}

		requestBody interface{}

		request *http.Request
	)

	Describe("ListVM", func() {
		var (
			vm1          models.VM
			vm2          models.VM
			params	     vm.ListVMParams
		)

		BeforeEach(func() {
			vm1 = models.VM{Cid: 1234567}
			vm2 = models.VM{Cid: 1234568}
			params = vm.NewListVMParams()
		})

		JustBeforeEach(func() {
			responseResponder = handler.ListVM(params)


		})

		Context("when reading tasks from controller succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VMResponse{&vm1, &vm2}
				controller.AllVirtualGuestsReturns(vms, nil)
			})

			It("returns a list of task", func() {
				responseResponder.WriteResponse()
			})

			It("calls the controller with no filter", func() {
				Expect(controller.TasksCallCount()).To(Equal(1))
				_, actualDomain, actualCellId := controller.TasksArgsForCall(0)
				Expect(actualDomain).To(Equal(domain))
				Expect(actualCellId).To(Equal(cellId))
			})

			Context("and filtering by domain", func() {
				BeforeEach(func() {
					domain = "domain-1"
				})

				It("calls the controller with a domain filter", func() {
					Expect(controller.TasksCallCount()).To(Equal(1))
					_, actualDomain, actualCellId := controller.TasksArgsForCall(0)
					Expect(actualDomain).To(Equal(domain))
					Expect(actualCellId).To(Equal(cellId))
				})
			})

			Context("and filtering by cell id", func() {
				BeforeEach(func() {
					cellId = "cell-id"
				})

				It("calls the controller with a cell filter", func() {
					Expect(controller.TasksCallCount()).To(Equal(1))
					_, actualDomain, actualCellId := controller.TasksArgsForCall(0)
					Expect(actualDomain).To(Equal(domain))
					Expect(actualCellId).To(Equal(cellId))
				})
			})
		})

	})
})
