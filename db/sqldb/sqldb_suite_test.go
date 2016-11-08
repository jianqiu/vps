package sqldb_test

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	thepackagedb "github.com/jianqiu/vps/db"
	"github.com/jianqiu/vps/db/sqldb"
	"github.com/jianqiu/vps/migration"
	"github.com/jianqiu/vps/test_helpers"
	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"

	_ "github.com/lib/pq"

	"testing"
)

var (
	db                                   *sql.DB
	sqlDB1                               *sqldb.SQLDB
	fakeClock                            *fakeclock.FakeClock
	logger                               *lagertest.TestLogger
	migrationProcess                     ifrit.Process
	dbDriverName, dbBaseConnectionString string
	dbFlavor                             string
)

func TestSql(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "SQL DB Suite")
}

var _ = BeforeSuite(func() {
	if !test_helpers.UseSQL() {
		return
	}

	var err error
	fakeClock = fakeclock.NewFakeClock(time.Now())
	logger = lagertest.NewTestLogger("sql-db")

	if test_helpers.UsePostgres() {
		dbDriverName = "postgres"
		dbBaseConnectionString = "postgres://diego:diego_pw@localhost/"
		dbFlavor = sqldb.Postgres
	} else if test_helpers.UseMySQL() {
		dbDriverName = "mysql"
		dbBaseConnectionString = "diego:diego_password@/"
		dbFlavor = sqldb.MySQL
	} else {
		panic("Unsupported driver")
	}

	// mysql must be set up on localhost as described in the CONTRIBUTING.md doc
	// in diego-release.
	db, err = sql.Open(dbDriverName, dbBaseConnectionString)
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Ping()).NotTo(HaveOccurred())

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE diego_%d", GinkgoParallelNode()))
	Expect(err).NotTo(HaveOccurred())

	db, err = sql.Open(dbDriverName, fmt.Sprintf("%sdiego_%d", dbBaseConnectionString, GinkgoParallelNode()))
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Ping()).NotTo(HaveOccurred())

	sqlDB1 = sqldb.NewSQLDB(db,fakeClock, dbFlavor)
	err = sqlDB1.CreateConfigurationsTable(logger)
	if err != nil {
		logger.Fatal("sql-failed-create-configurations-table", err)
	}

	// ensures sqlDB matches the db.DB interface
	var _ thepackagedb.DB = sqlDB1
})

var _ = BeforeEach(func() {
	if !test_helpers.UseSQL() {
		Skip("SQL Backend not available")
	}

	migrationsDone := make(chan struct{})

	migrationManager := migration.NewManager(logger,
		sqlDB1,
		db,
		migrationsDone,
		fakeClock,
		dbDriverName,
	)

	migrationProcess = ifrit.Invoke(migrationManager)

	Consistently(migrationProcess.Wait()).ShouldNot(Receive())
	Eventually(migrationsDone).Should(BeClosed())
})

var _ = AfterEach(func() {
	if test_helpers.UseSQL() {
		truncateTables(db)
	}
})

var _ = AfterSuite(func() {
	if test_helpers.UseSQL() {
		if migrationProcess != nil {
			migrationProcess.Signal(os.Kill)
		}

		Expect(db.Close()).NotTo(HaveOccurred())
		db, err := sql.Open(dbDriverName, dbBaseConnectionString)
		Expect(err).NotTo(HaveOccurred())
		Expect(db.Ping()).NotTo(HaveOccurred())
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE diego_%d", GinkgoParallelNode()))
		Expect(err).NotTo(HaveOccurred())
		Expect(db.Close()).NotTo(HaveOccurred())
	}
})

func truncateTables(db *sql.DB) {
	for _, query := range truncateTablesSQL {
		result, err := db.Exec(query)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RowsAffected()).To(BeEquivalentTo(0))
	}
}

var truncateTablesSQL = []string{
	"TRUNCATE TABLE virtual_guests",
}
