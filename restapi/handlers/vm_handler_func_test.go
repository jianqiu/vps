package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jianqiu/vps/restapi/handlers"
	"github.com/jianqiu/vps/restapi/handlers/fake_controllers"
	"github.com/jianqiu/vps/models"
	"github.com/jianqiu/vps/restapi/operations/vm"
	"github.com/go-openapi/runtime/middleware"
	"code.cloudfoundry.org/lager/lagertest"
)

var _ = Describe("VmHandlerFunc", func() {
	var (
		logger     *lagertest.TestLogger
		controller *fake_controllers.FakeVirtualGuestController

		responseResponder middleware.Responder

		handler *handlers.VMHandler
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		controller = &fake_controllers.FakeVirtualGuestController{}
		handler = handlers.NewVmHandler(logger, controller)
	})

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

		Context("when reading virtual guests from controller succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				controller.AllVirtualGuestsReturns(vms, nil)
			})

			It("returns a list of virtual guests", func() {
				Expect(controller.AllVirtualGuestsCallCount()).To(Equal(1))
				listVMOK, ok:=responseResponder.(*vm.ListVMOK)
				Expect(ok).To(BeTrue())
				Expect(listVMOK.GetPayload().Vms).To(Equal(vms))
			})
		})

		Context("when the controller returns no virtual guest", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{}
				controller.AllVirtualGuestsReturns(vms, nil)
			})

			It("returns 404 status code", func() {
				listVmNotFound, ok:=responseResponder.(*vm.ListVMNotFound)
				Expect(ok).To(BeTrue())
				Expect(listVmNotFound.GetStatusCode()).To(Equal(404))
			})
		})

		Context("when the controller errors out", func() {
			BeforeEach(func() {
				controller.AllVirtualGuestsReturns(nil, models.ErrUnknownError)
			})

			It("provides relevant error information", func() {
				listVmDefault, ok:=responseResponder.(*vm.ListVMDefault)
				Expect(ok).To(BeTrue())
				Expect(listVmDefault.GetStatusCode()).To(Equal(500))
				Expect(listVmDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})

	Describe("GetVMByCid", func() {
		var (
			params vm.GetVMByCidParams
		)

		BeforeEach(func() {
			params = vm.NewGetVMByCidParams()
			params.Cid = 1234567
		})

		JustBeforeEach(func() {
			responseResponder = handler.GetVMByCid(params)
		})

		Context("when reading a virtual guest from the controller succeeds", func() {
			var vm1 *models.VM

			BeforeEach(func() {
				vm1 = &models.VM{Cid: params.Cid}
				controller.VirtualGuestByCidReturns(vm1, nil)
			})

			It("fetches virtual guest by cid", func() {
				Expect(controller.VirtualGuestByCidCallCount()).To(Equal(1))
				_, actualCid := controller.VirtualGuestByCidArgsForCall(0)
				Expect(actualCid).To(Equal(params.Cid))
			})

			It("returns the virtual guest", func() {
				getVmByCidOK, ok:=responseResponder.(*vm.GetVMByCidOK)
				Expect(ok).To(BeTrue())
				Expect(getVmByCidOK.GetPayload().VM).To(Equal(vm1))
			})
		})

		Context("when the controller returns no virtual guest", func() {
			BeforeEach(func() {
				controller.VirtualGuestByCidReturns(nil, nil)
			})

			It("returns 404 status code", func() {
				getVmByCidNotFound, ok:=responseResponder.(*vm.GetVMByCidNotFound)
				Expect(ok).To(BeTrue())
				Expect(getVmByCidNotFound.GetStatusCode()).To(Equal(404))
			})
		})

		Context("when the controller errors out", func() {
			BeforeEach(func() {
				controller.VirtualGuestByCidReturns(nil, models.ErrUnknownError)
			})

			It("provides relevant error information", func() {
				getVmByCidDefault, ok:=responseResponder.(*vm.GetVMByCidDefault)
				Expect(ok).To(BeTrue())
				Expect(getVmByCidDefault.GetStatusCode()).To(Equal(500))
				Expect(getVmByCidDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})

	Describe("DeleteVM", func() {
		Context("when the delete request is normal", func() {
			var (
				params vm.DeleteVMParams
			)

			BeforeEach(func() {
				params = vm.NewDeleteVMParams()
				params.Cid = 1234567
			})

			JustBeforeEach(func() {
				responseResponder = handler.DeleteVM(params)
			})

			Context("when deleting the virtual guest succeeds", func() {
				It("returns no error", func() {
					Expect(controller.DeleteVMCallCount()).To(Equal(1))
					_, cid := controller.DeleteVMArgsForCall(0)
					Expect(cid).To(Equal(params.Cid))

					deleteVmNoContent, ok:=responseResponder.(*vm.DeleteVMNoContent)
					Expect(ok).To(BeTrue())
					Expect(deleteVmNoContent.GetPayload()).To(Equal("vm removed"))
				})
			})

			Context("when the controller returns an error", func() {
				BeforeEach(func() {
					controller.DeleteVMReturns(models.ErrUnknownError)
				})

				It("provides relevant error information", func() {
					deleteVmDefault, ok:=responseResponder.(*vm.DeleteVMDefault)
					Expect(ok).To(BeTrue())
					Expect(deleteVmDefault.GetStatusCode()).To(Equal(500))
					Expect(deleteVmDefault.GetPayload()).To(Equal(models.ErrUnknownError))
				})
			})
		})
	})

	Describe("AddVM", func() {
		var (
			vm1  *models.VM
			params vm.AddVMParams
		)

		BeforeEach(func() {
			vm1 = &models.VM {
				Cid: 1234567,
				CPU: 4,
				MemoryMb: 1024,
				IP: "10.0.0.1",
				Hostname: "fake.test.com",
				PrivateVlan: 1234567,
				PublicVlan: 1234568,
			}
			params = vm.NewAddVMParams()
			params.Body = vm1
		})

		JustBeforeEach(func() {
			responseResponder = handler.AddVM(params)
		})

		Context("when the virtual guest is added successful", func() {
			It("added into pool", func() {
				Expect(controller.CreateVMCallCount()).To(Equal(1))
				_, actualVm := controller.CreateVMArgsForCall(0)
				Expect(actualVm).To(Equal(vm1))

				addVmOk, ok:=responseResponder.(*vm.AddVMOK)
				Expect(ok).To(BeTrue())
				Expect(addVmOk.GetPayload()).To(Equal("added successfully"))
			})
		})

		Context("when adding virtual guest fails", func() {
			BeforeEach(func() {
				controller.CreateVMReturns(models.ErrUnknownError)
			})

			It("responds with an error", func() {
				addVmDefault, ok:=responseResponder.(*vm.AddVMDefault)
				Expect(ok).To(BeTrue())
				Expect(addVmDefault.GetStatusCode()).To(Equal(500))
				Expect(addVmDefault.GetPayload()).To(Equal(models.ErrUnknownError))
			})
		})
	})
})
