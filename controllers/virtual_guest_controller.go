package controllers

import (
	"code.cloudfoundry.org/lager"
	"github.com/jianqiu/vps/db"
	"github.com/jianqiu/vps/models"
)

type VirtualGuestController struct {
	db            db.VirtualGuestDB
}

func NewVirtualGuestController(
	db db.VirtualGuestDB,
) *VirtualGuestController {
	return &VirtualGuestController{
		db:            db,
	}
}

func (h *VirtualGuestController) AllVirtualGuests(logger lager.Logger) ([]*models.VM, error) {
	logger = logger.Session("vms")

	filter := models.VMFilter{}

	return h.db.VirtualGuests(logger, filter)
}

func (h *VirtualGuestController) VirtualGuests(logger lager.Logger, publicVlan, privateVlan, cpu, memory_mb int32, state models.State) ([]*models.VM, error) {
	logger = logger.Session("vms")

	filter := models.VMFilter{
		CPU: cpu,
		MemoryMb: memory_mb,
		PublicVlan: publicVlan,
		PrivateVlan: privateVlan,
		State: state,
	}

	return h.db.VirtualGuests(logger, filter)
}

func (h *VirtualGuestController) CreateVM(logger lager.Logger, vmDefinition *models.VM) error {
	var err error
	logger = logger.Session("create-vm")

	err = h.db.InsertVirtualGuestToPool(logger, vmDefinition)
	if err != nil {
		return err
	}

	return nil
}

func (h *VirtualGuestController) DeleteVM(logger lager.Logger, cid int32) error {
	logger = logger.Session("delete-vm")

	return h.db.DeleteVirtualGuestFromPool(logger, cid)
}

func (h *VirtualGuestController) UpdateVM(logger lager.Logger, cid int32, updateData *models.State) error {
	var err error
	logger = logger.Session("update-vm")

	switch *updateData {
	case models.StateUsing:
		err = h.db.ChangeVirtualGuestToUse(logger, cid)
	case models.StateFree:
		err = h.db.ChangeVirtualGuestToFree(logger, cid)
	case models.StateProvisioning:
		err = h.db.ChangeVirtualGuestToProvision(logger, cid)
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *VirtualGuestController) VirtualGuestByCid(logger lager.Logger, cid int32) (*models.VM, error) {
	logger = logger.Session("vm-by-guid")
	return h.db.VirtualGuestByCID(logger, cid)
}
