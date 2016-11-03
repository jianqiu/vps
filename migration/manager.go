package migration

import (
	"database/sql"
	"os"
	"time"

	"github.com/jianqiu/vps/db"
	"github.com/jianqiu/vps/db/sqldb"
	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/runtimeschema/metric"
)

const (
	migrationDuration = metric.Duration("MigrationDuration")
)

type Manager struct {
	logger         lager.Logger
	sqlDB          db.DB
	rawSQLDB       *sql.DB
	migrationsDone chan<- struct{}
	clock          clock.Clock
	databaseDriver string
}

func NewManager(
logger lager.Logger,
sqlDB db.DB,
rawSQLDB *sql.DB,
migrationsDone chan<- struct{},
clock clock.Clock,
databaseDriver string,
) Manager {
	return Manager{
		logger:         logger,
		sqlDB:          sqlDB,
		rawSQLDB:       rawSQLDB,
		migrationsDone: migrationsDone,
		clock:          clock,
		databaseDriver: databaseDriver,
	}
}

func (m Manager) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	logger := m.logger.Session("migration-manager")
	logger.Info("starting")

	errorChan := make(chan error)
	go m.performMigration(logger, errorChan, ready)
	defer logger.Info("exited")

	select {
	case err := <-errorChan:
		logger.Error("migration-failed", err)
		return err
	case <-signals:
		logger.Info("migration-interrupt")
		return nil
	}
}

func (m *Manager) performMigration(
logger lager.Logger,
errorChan chan error,
readyChan chan<- struct{},
) {
	migrateStart := m.clock.Now()

	logger.Info("running-migration", lager.Data{
		"init":   "virtual_guest_db",
	})

        if checkTables(logger, m.rawSQLDB) {
	       err := createTables(logger, m.rawSQLDB, m.databaseDriver)
	       if err != nil {
		       errorChan <- err
		       return
	       }

	       err = createIndices(logger, m.rawSQLDB)
	       if err != nil {
		       errorChan <- err
		       return
	       }
        }

	logger.Debug("migrations-finished")

	err := migrationDuration.Send(time.Since(migrateStart))
	if err != nil {
		logger.Error("failed-to-send-migration-duration-metric", err)
	}

	m.finish(logger, readyChan)
}

func (m *Manager) finish(logger lager.Logger, ready chan<- struct{}) {
	close(ready)
	close(m.migrationsDone)
	logger.Info("finished-migrations")
}


func checkTables(logger lager.Logger, db *sql.DB) bool {
	var value int
	db.QueryRow("SELECT 1 FROM information_schema.tables WHERE table_name = 'virtual_guests' LIMIT 1").Scan(&value)
	if value == 0 {
		// check whether the table exists before truncating
		return true
	}

	return false
}

func createTables(logger lager.Logger, db *sql.DB, flavor string) error {
	var createTablesSQL = []string{
		sqldb.RebindForFlavor(createVirtualGuestsSQL, flavor),
	}

	logger.Info("creating-tables")
	for _, query := range createTablesSQL {
		logger.Info("creating the table", lager.Data{"query": query})
		_, err := db.Exec(query)
		if err != nil {
			logger.Error("failed-creating-tables", err)
			return err
		}
		logger.Info("created the table", lager.Data{"query": query})
	}

	return nil
}

func createIndices(logger lager.Logger, db *sql.DB) error {
	logger.Info("creating-indices")
	createIndicesSQL := []string{}
	createIndicesSQL = append(createIndicesSQL, createVirtualGuestsIndices...)

	for _, query := range createIndicesSQL {
		logger.Info("creating the index", lager.Data{"query": query})
		_, err := db.Exec(query)
		if err != nil {
			logger.Error("failed-creating-index", err)
			return err
		}
		logger.Info("created the index", lager.Data{"query": query})
	}

	return nil
}

const createVirtualGuestsSQL = `CREATE TABLE virtual_guests(
	cid INT PRIMARY KEY,
	hostname VARCHAR(255) NOT NULL,
	ip VARCHAR(255) NOT NULL,
	cpu INT NOT NULL,
	memory_mb INT NOT NULL,
	private_vlan INT NOT NULL,
	public_vlan INT NOT NULL,
	deployment_name VARCHAR(255) NOT NULL DEFAULT '',
	state VARCHAR(255) NOT NULL,
	updated_at BIGINT DEFAULT 0,
	created_at BIGINT DEFAULT 0
);`

var createVirtualGuestsIndices = []string{
	`CREATE INDEX virtual_guests_cid_idx ON virtual_guests (cid)`,
	`CREATE INDEX virtual_guests_state_idx ON virtual_guests (state)`,
}



