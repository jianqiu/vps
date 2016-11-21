package sqldb

import (
	"database/sql"
	"strings"
	"github.com/go-openapi/strfmt"

	"github.com/jianqiu/vps/models"
	"code.cloudfoundry.org/lager"
)

func (db *SQLDB) VirtualGuests(logger lager.Logger, filter models.VMFilter) ([]*models.VM, error) {
	logger = logger.Session("vms-by-filter", lager.Data{"filter": filter})
	logger.Debug("starting")
	defer logger.Debug("complete")

	wheres := []string{}
	values := []interface{}{}

	if filter.CPU > 0 {
		wheres = append(wheres, "cpu = ?")
		values = append(values, filter.CPU)
	}

	if filter.MemoryMb > 0 {
		wheres = append(wheres, "memory_mb = ?")
		values = append(values, filter.MemoryMb)
	}

	if filter.PrivateVlan >0 {
		wheres = append(wheres, "private_vlan = ?")
		values = append(values, filter.PrivateVlan)
	}

	if filter.PublicVlan > 0 {
		wheres = append(wheres, "public_vlan = ?")
		values = append(values, filter.PublicVlan)
	}

	if filter.DeploymentName != "" {
		wheres = append(wheres, "deployment_name = ?")
		values = append(values, filter.DeploymentName)
	}

	switch filter.State {
	case models.StateUsing:
		wheres = append(wheres, "state = ?")
		values = append(values, "using")
	case models.StateProvisioning:
		wheres = append(wheres, "state = ?")
		values = append(values, "provisioning")
	case models.StateFree:
		wheres = append(wheres, "state = ?")
		values = append(values, "free")
	default:
	}

	rows, err := db.all(logger, db.db, virtualGuests,
		virtualGuestColumns, NoLockRow,
		strings.Join(wheres, " AND "), values...,
	)
	if err != nil {
		logger.Error("failed-query", err)
		return nil, db.convertSQLError(err)
	}
	defer rows.Close()

	results := []*models.VM{}
	for rows.Next() {
		vm, err := db.fetchVirtualGuest(logger, rows, db.db)
		if err != nil {
			logger.Error("failed-fetch", err)
			return nil, err
		}
		results = append(results, vm)
	}

	if rows.Err() != nil {
		logger.Error("failed-getting-next-row", rows.Err())
		return nil, db.convertSQLError(rows.Err())
	}

	return results, nil
}

func (db *SQLDB) OrderVirtualGuestToProvision(logger lager.Logger, filter models.VMFilter) (*models.VM, error) {
	logger = logger.Session("order-free-vm", lager.Data{"filter": filter})
	logger.Debug("starting")
	defer logger.Debug("complete")

	var vm *models.VM
	var err error

	err = db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		vm, err = db.fetchOneVMWithFilter(logger, filter, tx)
		if err != nil {
			logger.Error("failed-locking-vm", err)
			return err
		}

		if err = vm.ValidateTransitionTo(models.StateProvisioning); err != nil {
			logger.Error("failed-to-transition-vm-to-provisioning", err)
			return err
		}

		logger.Info("starting")
		defer logger.Info("complete")
		now := db.clock.Now().UnixNano()
		_, err = db.update(logger, tx, virtualGuests,
			SQLAttributes{
				"state":      "provisioning",
				"updated_at": now,
			},
			"cid = ?", vm.Cid,
		)
		if err != nil {
			return db.convertSQLError(err)
		}

		return nil
	})

	return vm, err
}

func (db *SQLDB) VirtualGuestByCID(logger lager.Logger, cid int32) (*models.VM, error) {
	logger = logger.Session("vm-by-cid", lager.Data{"cid": cid})
	logger.Debug("starting")
	defer logger.Debug("complete")

	row, _ := db.one(logger, db.db, virtualGuests,
		virtualGuestColumns, NoLockRow,
		"cid = ?", cid,
	)
	return db.fetchVirtualGuest(logger, row, db.db)
}

func (db *SQLDB) VirtualGuestByIP(logger lager.Logger, ip string) (*models.VM, error) {
	logger = logger.Session("vm-by-ip", lager.Data{"ip": ip})
	logger.Debug("starting")
	defer logger.Debug("complete")

	row, _:= db.one(logger, db.db, virtualGuests,
		virtualGuestColumns, NoLockRow,
		"ip = ?", ip,
	)
	return db.fetchVirtualGuest(logger, row, db.db)
}

func (db *SQLDB) VirtualGuestsByDeployments(logger lager.Logger, names []string) ([]*models.VM, error) {
	logger = logger.Session("vms-by-deployment", lager.Data{"filter": names})
	logger.Debug("starting")
	defer logger.Debug("complete")

	wheres := []string{}
	values := []interface{}{}

	for _, name := range names {
		wheres = append(wheres, "deployment_name = ?")
		values = append(values, name)
	}

	rows, err := db.all(logger, db.db, virtualGuests,
		virtualGuestColumns, LockRow,
		strings.Join(wheres, " OR "), values...,
	)
	if err != nil {
		logger.Error("failed-query", err)
		return nil, db.convertSQLError(err)
	}
	defer rows.Close()

	results := []*models.VM{}
	for rows.Next() {
		vm, err := db.fetchVirtualGuest(logger, rows, db.db)
		if err != nil {
			logger.Error("failed-fetch", err)
			return nil, err
		}
		results = append(results, vm)
	}

	if rows.Err() != nil {
		logger.Error("failed-getting-next-row", rows.Err())
		return nil, db.convertSQLError(rows.Err())
	}

	return results, nil
}

func (db *SQLDB) VirtualGuestsByStates(logger lager.Logger, states []string) ([]*models.VM, error) {
	logger = logger.Session("vms-by-state", lager.Data{"filter": states})
	logger.Debug("starting")
	defer logger.Debug("complete")

	wheres := []string{}
	values := []interface{}{}

	for _, state := range states {
		switch state {
		case "using":
			wheres = append(wheres, "state = ?")
			values = append(values, "using")
		case "provisioning":
			wheres = append(wheres, "state = ?")
			values = append(values, "provisioning")
		case "free":
			wheres = append(wheres, "state = ?")
			values = append(values, "free")
		default:
		}
	}

	rows, err := db.all(logger, db.db, virtualGuests,
		virtualGuestColumns, LockRow,
		strings.Join(wheres, " OR "), values...,
	)
	if err != nil {
		logger.Error("failed-query", err)
		return nil, db.convertSQLError(err)
	}
	defer rows.Close()

	results := []*models.VM{}
	for rows.Next() {
		vm, err := db.fetchVirtualGuest(logger, rows, db.db)
		if err != nil {
			logger.Error("failed-fetch", err)
			return nil, err
		}
		results = append(results, vm)
	}

	if rows.Err() != nil {
		logger.Error("failed-getting-next-row", rows.Err())
		return nil, db.convertSQLError(rows.Err())
	}

	return results, nil
}

func (db *SQLDB) InsertVirtualGuestToPool(logger lager.Logger, virtualGuest *models.VM) error {
	logger = logger.Session("insert-vm-to-pool", lager.Data{"cid": virtualGuest.Cid})
	logger.Info("starting")
	defer logger.Info("complete")

	now := db.clock.Now().UnixNano()

	var stateString string

	switch virtualGuest.State {
	case models.StateUsing:
		stateString = "using"
	case models.StateProvisioning:
		stateString = "provisioning"
	case models.StateFree:
		stateString = "free"
	default:
		stateString = "unknown"
	}

	_, err := db.insert(logger, db.db, virtualGuests,
		SQLAttributes{
			"cid":               virtualGuest.Cid,
			"hostname":          virtualGuest.Hostname,
			"ip":		     virtualGuest.IP,
			"cpu":		     virtualGuest.CPU,
			"memory_mb":         virtualGuest.MemoryMb,
			"public_vlan":	     virtualGuest.PublicVlan,
			"private_vlan":      virtualGuest.PrivateVlan,
			"created_at":         now,
			"updated_at":         now,
			"deployment_name":    virtualGuest.DeploymentName,
			"state":              stateString,
		},
	)
	if err != nil {
		logger.Error("failed-inserting-vm", err)
		return db.convertSQLError(err)
	}

	return nil
}

func (db *SQLDB) UpdateVirtualGuestInPool(logger lager.Logger, virtualGuest *models.VM) error {
	logger = logger.Session("update-vm-in-pool", lager.Data{"cid":virtualGuest.Cid})

	err := db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		_, err := db.fetchVMForUpdate(logger, virtualGuest.Cid, tx)
		if err != nil {
			logger.Error("failed-locking-vm", err)
			return err
		}

		logger.Info("starting")
		defer logger.Info("complete")
		now := db.clock.Now().UnixNano()

		_, err = db.update(logger, tx, virtualGuests,
			SQLAttributes{
				"hostname": virtualGuest.Hostname,
				"ip":  virtualGuest.IP,
				"cpu":  virtualGuest.CPU,
				"memory_mb":  virtualGuest.MemoryMb,
				"deployment_name":  virtualGuest.DeploymentName,
				"public_vlan":  virtualGuest.PublicVlan,
				"private_vlan":  virtualGuest.PrivateVlan,
				"updated_at": now,
			},
			"cid = ?", virtualGuest.Cid,
		)
		if err != nil {
			return db.convertSQLError(err)
		}

		return nil
	})

	return err
}

func (db *SQLDB) ChangeVirtualGuestToProvision(logger lager.Logger, cid int32) error {
	logger = logger.Session("update-vm-to-provisioning", lager.Data{"cid": cid})

	err := db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		vm, err := db.fetchVMForUpdate(logger, cid, tx)
		if err != nil {
			logger.Error("failed-locking-vm", err)
			return err
		}

		if err = vm.ValidateTransitionTo(models.StateProvisioning); err != nil {
			logger.Error("failed-to-transition-vm-to-provisioning", err)
			return err
		}

		logger.Info("starting")
		defer logger.Info("complete")
		now := db.clock.Now().UnixNano()
		_, err = db.update(logger, tx, virtualGuests,
			SQLAttributes{
				"state":      "provisioning",
				"updated_at": now,
			},
			"cid = ?", cid,
		)
		if err != nil {
			return db.convertSQLError(err)
		}

		return nil
	})

	return err
}

func (db *SQLDB) ChangeVirtualGuestToUse(logger lager.Logger, cid int32) error {
	logger = logger.Session("update-vm-to-use", lager.Data{"cid": cid})

	err := db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		vm, err := db.fetchVMForUpdate(logger, cid, tx)
		if err != nil {
			logger.Error("failed-locking-vm", err)
			return err
		}

		if err = vm.ValidateTransitionTo(models.StateUsing); err != nil {
			logger.Error("failed-to-transition-vm-to-running", err)
			return err
		}

		logger.Info("starting")
		defer logger.Info("complete")
		now := db.clock.Now().UnixNano()
		_, err = db.update(logger, tx, virtualGuests,
			SQLAttributes{
				"state":      "using",
				"updated_at": now,
			},
			"cid = ?", cid,
		)
		if err != nil {
			return db.convertSQLError(err)
		}

		return nil
	})

	return err
}

func (db *SQLDB) ChangeVirtualGuestToFree(logger lager.Logger, cid int32) error {
	logger = logger.Session("update-vm-to-free", lager.Data{"cid": cid})

	err := db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		vm, err := db.fetchVMForUpdate(logger, cid, tx)
		if err != nil {
			logger.Error("failed-locking-vm", err)
			return err
		}

		if err = vm.ValidateTransitionTo(models.StateFree); err != nil {
			logger.Error("failed-to-transition-vm-to-free", err)
			return err
		}

		logger.Info("starting")
		defer logger.Info("complete")
		now := db.clock.Now().UnixNano()
		_, err = db.update(logger, tx, virtualGuests,
			SQLAttributes{
				"state":      "free",
				"updated_at": now,
			},
			"cid = ?", cid,
		)
		if err != nil {
			return db.convertSQLError(err)
		}

		return nil
	})

	return err
}

func (db *SQLDB) DeleteVirtualGuestFromPool(logger lager.Logger, cid int32) error {
	logger = logger.Session("delete-vm-from-pool", lager.Data{"cid": cid})
	logger.Info("starting")
	defer logger.Info("complete")

	return db.transact(logger, func(logger lager.Logger, tx *sql.Tx) error {
		vm, err := db.fetchVMForUpdate(logger, cid, tx)
		if err != nil {
			logger.Error("failed-locking-vm", err)
			return err
		}

		if vm.State != models.StateFree {
			err = models.NewTaskTransitionError(vm.State, models.StateFree)
			logger.Error("invalid-state-transition", err)
			return err
		}

		_, err = db.delete(logger, tx, virtualGuests, "cid = ?", cid)
		if err != nil {
			logger.Error("failed-removing-vm", err)
			return db.convertSQLError(err)
		}

		return nil
	})
}

func (db *SQLDB) fetchVMForUpdate(logger lager.Logger, cid int32, tx *sql.Tx) (*models.VM, error) {
	row, _:= db.one(logger, tx, virtualGuests,
		virtualGuestColumns, LockRow,
		"cid = ?", cid,
	)
	return db.fetchVirtualGuest(logger, row, tx)
}

func (db *SQLDB) fetchOneVMWithFilter(logger lager.Logger, filter models.VMFilter, tx *sql.Tx) (*models.VM, error) {
	wheres := []string{}
	values := []interface{}{}

	if filter.CPU > 0 {
		wheres = append(wheres, "cpu = ?")
		values = append(values, filter.CPU)
	}

	if filter.MemoryMb > 0 {
		wheres = append(wheres, "memory_mb = ?")
		values = append(values, filter.MemoryMb)
	}

	if filter.PrivateVlan >0 {
		wheres = append(wheres, "private_vlan = ?")
		values = append(values, filter.PrivateVlan)
	}

	if filter.PublicVlan > 0 {
		wheres = append(wheres, "public_vlan = ?")
		values = append(values, filter.PublicVlan)
	}

	if filter.DeploymentName != "" {
		wheres = append(wheres, "deployment_name = ?")
		values = append(values, filter.DeploymentName)
	}

	switch filter.State {
	case models.StateUsing:
		wheres = append(wheres, "state = ?")
		values = append(values, "using")
	case models.StateProvisioning:
		wheres = append(wheres, "state = ?")
		values = append(values, "provisioning")
	case models.StateFree:
		wheres = append(wheres, "state = ?")
		values = append(values, "free")
	default:
	}

	row, _ := db.one(logger, db.db, virtualGuests,
		virtualGuestColumns, LockRow,
		strings.Join(wheres, " AND "), values...,
	)
	return db.fetchVirtualGuest(logger, row, tx)
}

func (db *SQLDB) fetchVirtualGuest(logger lager.Logger, scanner RowScanner, tx Queryable) (*models.VM, error) {
	var hostname, deployment_name, state string
	var cpu, memory_mb, cid, public_vlan, private_vlan int32
	var ip strfmt.IPv4
	err := scanner.Scan(
		&cid,
		&hostname,
		&ip,
		&cpu,
		&memory_mb,
		&private_vlan,
		&public_vlan,
		&deployment_name,
		&state,
	)
	if err != nil {
		logger.Error("failed-scanning-row", err)
		return nil, models.ErrResourceNotFound
	}

	virtualGuest := &models.VM{
		Cid:              cid,
		Hostname:         hostname,
		IP:               ip,
		CPU:    	  cpu,
		MemoryMb: 	  memory_mb,
		PrivateVlan:      private_vlan,
		PublicVlan:       public_vlan,
		DeploymentName:   deployment_name,
	}
	switch state {
	case "free":
		virtualGuest.State = models.StateFree
	case "provisioning":
		virtualGuest.State = models.StateProvisioning
	case "using":
		virtualGuest.State = models.StateUsing
	default:
		virtualGuest.State = models.StateUnknown
	}

	return virtualGuest, nil
}
