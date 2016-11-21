package sqldb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jianqiu/vps/models"
	"code.cloudfoundry.org/lager"
)

const (
	MySQL    = "mysql"
	Postgres = "postgres"
)

type RowLock bool

const (
	LockRow   RowLock = true
	NoLockRow RowLock = false
)

type SQLAttributes map[string]interface{}

type ColumnList []string

const (
	virtualGuests = "virtual_guests"
)

var (
	virtualGuestColumns = ColumnList{
		virtualGuests + ".cid",
		virtualGuests + ".hostname",
		virtualGuests + ".ip",
		virtualGuests + ".cpu",
		virtualGuests + ".memory_mb",
		virtualGuests + ".private_vlan",
		virtualGuests + ".public_vlan",
		virtualGuests + ".deployment_name",
		virtualGuests + ".state",
	}
)

func (db *SQLDB) CreateConfigurationsTable(logger lager.Logger) error {
	_, err := db.db.Exec(`
		CREATE TABLE IF NOT EXISTS configurations(
			id VARCHAR(255) PRIMARY KEY,
			value VARCHAR(255)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// Takes in a query that uses question marks to represent unbound SQL parameters
// and converts those to '$1, $2', etc. if the DB flavor is postgres.
// Takes in a query that uses MEDIUMTEXT to create table columns and converts
// those to TEXT if the DB flavor is postgres
// e.g., `SELECT * FROM table_name WHERE col = ? AND col2 = ?` becomes
//       `SELECT * FROM table_name WHERE col = $1 AND col2 = $2`
// e.g., `CREATE TABLE desired_lrps(
//	 annotation MEDIUMTEXT
// )` becomes
// `CREATE TABLE desired_lrps(
//	 annotation TEXT
// )`
func RebindForFlavor(query, flavor string) string {
	if flavor == MySQL {
		return query
	}
	if flavor != Postgres {
		panic(fmt.Sprintf("Unrecognized DB flavor '%s'", flavor))
	}

	strParts := strings.Split(query, "?")
	for i := 1; i < len(strParts); i++ {
		strParts[i-1] = fmt.Sprintf("%s$%d", strParts[i-1], i)
	}
	return strings.Replace(strings.Join(strParts, ""), "MEDIUMTEXT", "TEXT", -1)
}

func (db *SQLDB) countVirtualGuestsByState(logger lager.Logger, q Queryable) (deletedCount, usingCount int) {
	var query string
	switch db.flavor {
	case Postgres:
		query = `
			SELECT
				COUNT(*) FILTER (WHERE state = $1) AS deleted_virtualGuests,
				COUNT(*) FILTER (WHERE state = $2) AS using_virtualGuests,
			FROM virtual_guests
		`
	case MySQL:
		query = `
			SELECT
				COUNT(IF(state = ?, 1, NULL)) AS deleted_virtualGuests,
				COUNT(IF(state = ?, 1, NULL)) AS using_virtualGuests,
			FROM virtual_guests
		`
	default:
		// totally shouldn't happen
		panic("database flavor not implemented: " + db.flavor)
	}

	row := db.db.QueryRow(query, models.StateFree, models.StateUsing)
	err := row.Scan(&deletedCount, &usingCount)
	if err != nil {
		logger.Error("failed-counting-virtuals", err)
	}
	return
}

// SELECT <columns> FROM <table> WHERE ... LIMIT 1 [FOR UPDATE]
func (db *SQLDB) one(logger lager.Logger, q Queryable, table string,
columns ColumnList, lockRow RowLock,
wheres string, whereBindings ...interface{},
) *sql.Row {
	query := fmt.Sprintf("SELECT %s FROM %s\n", strings.Join(columns, ", "), table)

	if len(wheres) > 0 {
		query += "WHERE " + wheres
	}

	query += "\nLIMIT 1"

	if lockRow {
		query += "\nFOR UPDATE"
	}

	logger.Info("one", lager.Data{"query": db.rebind(query)})
	logger.Info("one", lager.Data{"bindings": whereBindings})

	return q.QueryRow(db.rebind(query), whereBindings...)
}

// SELECT <columns> FROM <table> WHERE ... [FOR UPDATE]
func (db *SQLDB) all(logger lager.Logger, q Queryable, table string,
columns ColumnList, lockRow RowLock,
wheres string, whereBindings ...interface{},
) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT %s FROM %s\n", strings.Join(columns, ", "), table)

	if len(wheres) > 0 {
		query += "WHERE " + wheres
	}

	if lockRow {
		query += "\nFOR UPDATE"
	}

	return q.Query(db.rebind(query), whereBindings...)
}

func (db *SQLDB) upsert(logger lager.Logger, q Queryable, table string, keyAttributes, updateAttributes SQLAttributes) (sql.Result, error) {
	columns := make([]string, 0, len(keyAttributes)+len(updateAttributes))
	keyNames := make([]string, 0, len(keyAttributes))
	updateBindings := make([]string, 0, len(updateAttributes))
	bindingValues := make([]interface{}, 0, len(keyAttributes)+2*len(updateAttributes))

	keyBindingValues := make([]interface{}, 0, len(keyAttributes))
	nonKeyBindingValues := make([]interface{}, 0, len(updateAttributes))

	for column, value := range keyAttributes {
		columns = append(columns, column)
		keyNames = append(keyNames, column)
		keyBindingValues = append(keyBindingValues, value)
	}

	for column, value := range updateAttributes {
		columns = append(columns, column)
		updateBindings = append(updateBindings, fmt.Sprintf("%s = ?", column))
		nonKeyBindingValues = append(nonKeyBindingValues, value)
	}

	insertBindings := questionMarks(len(keyAttributes) + len(updateAttributes))

	var query string
	switch db.flavor {
	case Postgres:
		bindingValues = append(bindingValues, nonKeyBindingValues...)
		bindingValues = append(bindingValues, keyBindingValues...)
		bindingValues = append(bindingValues, keyBindingValues...)
		bindingValues = append(bindingValues, nonKeyBindingValues...)

		insert := fmt.Sprintf(`
				INSERT INTO %s
					(%s)
				SELECT %s`,
			table,
			strings.Join(columns, ", "),
			insertBindings)

		// TODO: Add where clause with key values.
		// Alternatively upgrade to postgres 9.5 :D
		whereClause := []string{}
		for _, key := range keyNames {
			whereClause = append(whereClause, fmt.Sprintf("%s = ?", key))
		}

		upsert := fmt.Sprintf(`
				UPDATE %s SET
					%s
				WHERE %s
				`,
			table,
			strings.Join(updateBindings, ", "),
			strings.Join(whereClause, " AND "),
		)

		query = fmt.Sprintf(`
				WITH upsert AS (%s RETURNING *)
				%s WHERE NOT EXISTS
				(SELECT * FROM upsert)
				`,
			upsert,
			insert)

		result, err := q.Exec(fmt.Sprintf("LOCK TABLE %s IN SHARE ROW EXCLUSIVE MODE", table))
		if err != nil {
			return result, err
		}

	case MySQL:
		bindingValues = append(bindingValues, keyBindingValues...)
		bindingValues = append(bindingValues, nonKeyBindingValues...)
		bindingValues = append(bindingValues, nonKeyBindingValues...)

		query = fmt.Sprintf(`
				INSERT INTO %s
					(%s)
				VALUES (%s)
				ON DUPLICATE KEY UPDATE
					%s
			`,
			table,
			strings.Join(columns, ", "),
			insertBindings,
			strings.Join(updateBindings, ", "),
		)
	default:
		// totally shouldn't happen
		panic("database flavor not implemented: " + db.flavor)
	}
	return q.Exec(db.rebind(query), bindingValues...)
}

// INSERT INTO <table> (...) VALUES ...
func (db *SQLDB) insert(logger lager.Logger, q Queryable, table string, attributes SQLAttributes) (sql.Result, error) {
	attributeCount := len(attributes)
	if attributeCount == 0 {
		return nil, nil
	}

	query := fmt.Sprintf("INSERT INTO %s\n", table)
	attributeNames := make([]string, 0, attributeCount)
	attributeBindings := make([]string, 0, attributeCount)
	bindings := make([]interface{}, 0, attributeCount)

	for column, value := range attributes {
		attributeNames = append(attributeNames, column)
		attributeBindings = append(attributeBindings, "?")
		bindings = append(bindings, value)
	}
	query += fmt.Sprintf("(%s)", strings.Join(attributeNames, ", "))
	query += fmt.Sprintf("VALUES (%s)", strings.Join(attributeBindings, ", "))

	return q.Exec(db.rebind(query), bindings...)
}

// UPDATE <table> SET ... WHERE ...
func (db *SQLDB) update(logger lager.Logger, q Queryable, table string, updates SQLAttributes, wheres string, whereBindings ...interface{}) (sql.Result, error) {
	updateCount := len(updates)
	if updateCount == 0 {
		return nil, nil
	}

	query := fmt.Sprintf("UPDATE %s SET\n", table)
	updateQueries := make([]string, 0, updateCount)
	bindings := make([]interface{}, 0, updateCount+len(whereBindings))

	for column, value := range updates {
		updateQueries = append(updateQueries, fmt.Sprintf("%s = ?", column))
		bindings = append(bindings, value)
	}
	query += strings.Join(updateQueries, ", ") + "\n"
	if len(wheres) > 0 {
		query += "WHERE " + wheres
		bindings = append(bindings, whereBindings...)
	}

	return q.Exec(db.rebind(query), bindings...)
}

// DELETE FROM <table> WHERE ...
func (db *SQLDB) delete(logger lager.Logger, q Queryable, table string, wheres string, whereBindings ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s\n", table)

	if len(wheres) > 0 {
		query += "WHERE " + wheres
	}

	return q.Exec(db.rebind(query), whereBindings...)
}

func (db *SQLDB) rebind(query string) string {
	return RebindForFlavor(query, db.flavor)
}

func questionMarks(count int) string {
	if count == 0 {
		return ""
	}
	return strings.Repeat("?, ", count-1) + "?"
}
