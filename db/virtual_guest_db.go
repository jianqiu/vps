package db

import (
	"github.com/jianqiu/vps/models"
	"code.cloudfoundry.org/lager"
)

//go:generate counterfeiter . VirtualGuestDB
type VirtualGuestDB interface {
	VirtualGuests(logger lager.Logger, filter models.VMFilter) ([]*models.VM, error)
	VirtualGuestByCID(logger lager.Logger, cid int32) (*models.VM, error)
	VirtualGuestByIP(logger lager.Logger, ip string) (*models.VM, error)

	InsertVirtualGuestToPool(logger lager.Logger,virtualGuest *models.VM) error
	ChangeVirtualGuestToProvision(logger lager.Logger, cid int32) error
	ChangeVirtualGuestToUse(logger lager.Logger, cid int32) error
	ChangeVirtualGuestToFree(logger lager.Logger, cid int32) error
	DeleteVirtualGuestFromPool(logger lager.Logger, cid int32) error
}

