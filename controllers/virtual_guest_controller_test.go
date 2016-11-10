package controllers_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jianqiu/vps/db/dbfakes"
	"github.com/jianqiu/vps/controllers"
	"github.com/jianqiu/vps/models"

	"code.cloudfoundry.org/lager/lagertest"
)

var _ = Describe("VirtualGuestController", func() {
	var (
		logger                   *lagertest.TestLogger
		fakeVirtualGuestDB               *dbfakes.FakeVirtualGuestDB
		controller *controllers.VirtualGuestController
	)

	BeforeEach(func() {
		fakeVirtualGuestDB = new(dbfakes.FakeVirtualGuestDB)
		logger = lagertest.NewTestLogger("test")
		controller = controllers.NewVirtualGuestController(fakeVirtualGuestDB)
	})

	Describe("AllVirtualGuests", func() {
		var (
			vm1 models.VM
			vm2 models.VM
			actualVms    []*models.VM
			err            error
		)

		BeforeEach(func() {
			vm1 = models.VM{Cid: 1234567}
			vm2 = models.VM{Cid: 1234568}
		})

		JustBeforeEach(func() {
			actualVms, err = controller.AllVirtualGuests(logger)
		})

		Context("when reading all vms from DB succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				fakeVirtualGuestDB.VirtualGuestsReturns(vms, nil)
			})

			It("returns a list of vm", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualVms).To(Equal(vms))
			})

			It("calls the DB with no filter", func() {
				Expect(fakeVirtualGuestDB.VirtualGuestsCallCount()).To(Equal(1))
				_, filter := fakeVirtualGuestDB.VirtualGuestsArgsForCall(0)
				Expect(filter).To(Equal(models.VMFilter{}))
			})
		})

		Context("when the DB returns an error", func() {
			BeforeEach(func() {
				fakeVirtualGuestDB.VirtualGuestsReturns(nil, errors.New("kaboom"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("kaboom"))
			})
		})
	})

	Describe("VirtualGuests", func() {
		var (
			public_vlan, private_vlan, cpu, memory_mb int32
			state models.State
			vm1 models.VM
			vm2 models.VM
			actualVms    []*models.VM
			err            error
		)

		BeforeEach(func() {
			vm1 = models.VM{
				Cid: 1234567,
				CPU: 4,
				MemoryMb: 1024,
				PublicVlan: 1234567,
				PrivateVlan: 12345678,
				State: models.StateFree,
			}
			vm2 = models.VM{
				Cid: 1234568,
				CPU: 4,
				MemoryMb: 1024,
				PublicVlan: 1234567,
				PrivateVlan: 12345678,
				State: models.StateFree,
			}
			state = models.StateUnknown
		})

		JustBeforeEach(func() {
			actualVms, err = controller.VirtualGuests(logger, public_vlan, private_vlan, cpu, memory_mb, state)
		})

		Context("when reading tasks from DB succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				fakeVirtualGuestDB.VirtualGuestsReturns(vms, nil)
			})

			Context("and filtering by public_vlan, private_vlan, cpu, memory_mb", func() {
				BeforeEach(func() {
					public_vlan = 1234567
					private_vlan = 12345678
					cpu = 4
					memory_mb = 1024
				})

				It("calls the DB with a domain filter", func() {
					Expect(fakeVirtualGuestDB.VirtualGuestsCallCount()).To(Equal(1))
					_, filter := fakeVirtualGuestDB.VirtualGuestsArgsForCall(0)
					Expect(filter.PublicVlan).To(Equal(public_vlan))
					Expect(filter.PrivateVlan).To(Equal(private_vlan))
					Expect(filter.CPU).To(Equal(cpu))
					Expect(filter.MemoryMb).To(Equal(memory_mb))
				})
			})

			Context("and filtering by state", func() {
				BeforeEach(func() {
					state = models.StateFree
				})

				It("calls the DB with a state", func() {
					Expect(fakeVirtualGuestDB.VirtualGuestsCallCount()).To(Equal(1))
					_, filter := fakeVirtualGuestDB.VirtualGuestsArgsForCall(0)
					Expect(filter.State).To(Equal(state))
				})
			})
		})
	})

	Describe("VirtualGuestsByDeployments", func() {
		var (
			deployment_names = []string{"depoyment1","deployment2","deployment3"}
			vm1 models.VM
			vm2 models.VM
			actualVms    []*models.VM
			err            error
		)

		BeforeEach(func() {
			vm1 = models.VM{
				Cid: 1234567,
				CPU: 4,
				MemoryMb: 1024,
				PublicVlan: 1234567,
				PrivateVlan: 12345678,
				State: models.StateFree,
			}
			vm2 = models.VM{
				Cid: 1234568,
				CPU: 4,
				MemoryMb: 1024,
				PublicVlan: 1234567,
				PrivateVlan: 12345678,
				State: models.StateFree,
			}
		})

		JustBeforeEach(func() {
			actualVms, err = controller.VirtualGuestsByDeployments(logger, deployment_names)
		})

		Context("when filtering vms by deployment names from the DB succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				fakeVirtualGuestDB.VirtualGuestsByDeploymentsReturns(vms, nil)
			})

			It("fetches vms by deployment names", func() {
				Expect(fakeVirtualGuestDB.VirtualGuestsByDeploymentsCallCount()).To(Equal(1))
				_, actualDeploymentNames := fakeVirtualGuestDB.VirtualGuestsByDeploymentsArgsForCall(0)
				Expect(actualDeploymentNames).To(Equal(deployment_names))
			})

			It("returns the vms", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualVms).To(Equal(vms))
			})
		})

		Context("when the DB errors out", func() {
			BeforeEach(func() {
				fakeVirtualGuestDB.VirtualGuestsByDeploymentsReturns(nil, errors.New("kaboom"))
			})

			It("provides relevant error information", func() {
				Expect(err).To(MatchError("kaboom"))
			})
		})
	})

	Describe("VirtualGuestsByStates", func() {
		var (
			states = []string{"state1","state2","state3"}
			vm1 models.VM
			vm2 models.VM
			actualVms    []*models.VM
			err            error
		)

		BeforeEach(func() {
			vm1 = models.VM{
				Cid: 1234567,
				CPU: 4,
				MemoryMb: 1024,
				PublicVlan: 1234567,
				PrivateVlan: 12345678,
				State: models.StateFree,
			}
			vm2 = models.VM{
				Cid: 1234568,
				CPU: 4,
				MemoryMb: 1024,
				PublicVlan: 1234567,
				PrivateVlan: 12345678,
				State: models.StateFree,
			}
		})

		JustBeforeEach(func() {
			actualVms, err = controller.VirtualGuestsByStates(logger, states)
		})

		Context("when filtering vms by states from the DB succeeds", func() {
			var vms []*models.VM

			BeforeEach(func() {
				vms = []*models.VM{&vm1, &vm2}
				fakeVirtualGuestDB.VirtualGuestsByStatesReturns(vms, nil)
			})

			It("fetches vms by states", func() {
				Expect(fakeVirtualGuestDB.VirtualGuestsByStatesCallCount()).To(Equal(1))
				_, actualDeploymentNames := fakeVirtualGuestDB.VirtualGuestsByStatesArgsForCall(0)
				Expect(actualDeploymentNames).To(Equal(states))
			})

			It("returns the vms", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualVms).To(Equal(vms))
			})
		})

		Context("when the DB errors out", func() {
			BeforeEach(func() {
				fakeVirtualGuestDB.VirtualGuestsByStatesReturns(nil, errors.New("kaboom"))
			})

			It("provides relevant error information", func() {
				Expect(err).To(MatchError("kaboom"))
			})
		})
	})

	Describe("VirtualGuestByCid", func() {
		var (
			cid   int32
			actualVm *models.VM
			err        error
		)

		BeforeEach(func() {
			cid = 1234567
		})

		JustBeforeEach(func() {
			actualVm, err = controller.VirtualGuestByCid(logger, cid)
		})

		Context("when reading a vm from the DB succeeds", func() {
			var vm1 *models.VM

			BeforeEach(func() {
				vm1 = &models.VM{Cid: int32(cid)}
				fakeVirtualGuestDB.VirtualGuestByCIDReturns(vm1, nil)
			})

			It("fetches task by cid", func() {
				Expect(fakeVirtualGuestDB.VirtualGuestByCIDCallCount()).To(Equal(1))
				_, actualCid := fakeVirtualGuestDB.VirtualGuestByCIDArgsForCall(0)
				Expect(actualCid).To(Equal(cid))
			})

			It("returns the vm", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(actualVm).To(Equal(vm1))
			})
		})

		Context("when the DB errors out", func() {
			BeforeEach(func() {
				fakeVirtualGuestDB.VirtualGuestByCIDReturns(nil, errors.New("kaboom"))
			})

			It("provides relevant error information", func() {
				Expect(err).To(MatchError("kaboom"))
			})
		})
	})

	Describe("DeleteVM", func() {
		Context("when the delete request is normal", func() {
			var (
				cid int32
				err      error
			)

			BeforeEach(func() {
				cid = 1234567
			})

			JustBeforeEach(func() {
				err = controller.DeleteVM(logger, cid)
			})

			Context("when deleting the vm succeeds", func() {
				It("returns no error", func() {
					Expect(fakeVirtualGuestDB.DeleteVirtualGuestFromPoolCallCount()).To(Equal(1))
					_, actualCid := fakeVirtualGuestDB.DeleteVirtualGuestFromPoolArgsForCall(0)
					Expect(actualCid).To(Equal(cid))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when desiring the vm fails", func() {
				BeforeEach(func() {
					fakeVirtualGuestDB.DeleteVirtualGuestFromPoolReturns(errors.New("kaboom"))
				})

				It("responds with an error", func() {
					Expect(err).To(MatchError("kaboom"))
				})
			})
		})
	})

	Describe("CreateVM", func() {
		Context("when the createVM request is normal", func() {
			var (
				vmDefinition *models.VM
				err   error
			)

			BeforeEach(func() {
				vmDefinition = &models.VM{
					Cid: 1234567,
					CPU: 4,
					MemoryMb: 1024,
					PublicVlan: 1234567,
					PrivateVlan: 12345678,
					State: models.StateFree,
				}
			})

			JustBeforeEach(func() {
				err = controller.CreateVM(logger, vmDefinition)
			})

			Context("when creating the vm with Free succeeds", func() {
				It("returns no error", func() {
					Expect(fakeVirtualGuestDB.InsertVirtualGuestToPoolCallCount()).To(Equal(1))
					_, actualVmDefinition := fakeVirtualGuestDB.InsertVirtualGuestToPoolArgsForCall(0)
					Expect(actualVmDefinition).To(Equal(vmDefinition))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when creating the vm fails", func() {
				BeforeEach(func() {
					fakeVirtualGuestDB.InsertVirtualGuestToPoolReturns(errors.New("kaboom"))
				})

				It("responds with an error", func() {
					Expect(err).To(MatchError("kaboom"))
				})
			})
		})
	})

	Describe("UpdateVM", func() {
		Context("when the updateVM request is normal", func() {
			var (
				vmDefinition *models.VM
				err   error
			)

			BeforeEach(func() {
				vmDefinition = &models.VM{
					Cid: 1234567,
					CPU: 4,
					MemoryMb: 1024,
					PublicVlan: 1234567,
					PrivateVlan: 12345678,
					State: models.StateFree,
				}
			})

			JustBeforeEach(func() {
				err = controller.UpdateVM(logger, vmDefinition)
			})

			Context("when updating the vm with Free succeeds", func() {
				It("returns no error", func() {
					Expect(fakeVirtualGuestDB.UpdateVirtualGuestInPoolCallCount()).To(Equal(1))
					_, actualVmDefinition := fakeVirtualGuestDB.UpdateVirtualGuestInPoolArgsForCall(0)
					Expect(actualVmDefinition).To(Equal(vmDefinition))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when updating the vm fails", func() {
				BeforeEach(func() {
					fakeVirtualGuestDB.UpdateVirtualGuestInPoolReturns(errors.New("kaboom"))
				})

				It("responds with an error", func() {
					Expect(err).To(MatchError("kaboom"))
				})
			})
		})
	})

	Describe("UpdateVMWithState", func() {
		Context("when the updateVMWithState request is normal", func() {
			var (
				vmState models.VMState
				err      error
				cid  int32
			)

			BeforeEach(func() {
				cid = 1234567
			})

			Context("when updating the vm with Free succeeds", func() {
				It("returns no error", func() {
					vmState = models.VMState{
						State:  models.StateFree,
					}
					err = controller.UpdateVMWithState(logger, cid, &vmState.State)
					Expect(fakeVirtualGuestDB.ChangeVirtualGuestToFreeCallCount()).To(Equal(1))
					_, actualCid := fakeVirtualGuestDB.ChangeVirtualGuestToFreeArgsForCall(0)
					Expect(actualCid).To(Equal(cid))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when updating the vm with Provisioing succeeds", func() {

				It("returns no error", func() {
					vmState = models.VMState{
						State:  models.StateProvisioning,
					}
					err = controller.UpdateVMWithState(logger, cid, &vmState.State)
					Expect(fakeVirtualGuestDB.ChangeVirtualGuestToProvisionCallCount()).To(Equal(1))
					_, actualCid := fakeVirtualGuestDB.ChangeVirtualGuestToProvisionArgsForCall(0)
					Expect(actualCid).To(Equal(cid))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when updating the vm with Using succeeds", func() {
				It("returns no error", func() {
					vmState = models.VMState{
						State: models.StateUsing,
					}
					err = controller.UpdateVMWithState(logger, cid, &vmState.State)
					Expect(fakeVirtualGuestDB.ChangeVirtualGuestToUseCallCount()).To(Equal(1))
					_, actualCid := fakeVirtualGuestDB.ChangeVirtualGuestToUseArgsForCall(0)
					Expect(actualCid).To(Equal(cid))
					Expect(err).NotTo(HaveOccurred())
				})
			})


			Context("when updating the vm with Using fails", func() {
				It("responds with an error", func() {
					vmState = models.VMState{
						State: models.StateUsing,
					}
					fakeVirtualGuestDB.ChangeVirtualGuestToUseReturns(errors.New("kaboom"))
					err = controller.UpdateVMWithState(logger, cid, &vmState.State)
					Expect(err).To(MatchError("kaboom"))
				})
			})
		})
	})
})