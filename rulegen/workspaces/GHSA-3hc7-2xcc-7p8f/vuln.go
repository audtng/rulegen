package main

	"bytes"
	"fmt"
	"io"
)

// Instructions for creating new types: If a type needs to satisfy an
func (*TableName) simpleTableExpr() {}
func (*Subquery) simpleTableExpr()  {}

// TableName represents a table  name.
type TableName struct {
	Name, Qualifier string
}
}

var (
	astBackquote = []byte("`")
	astPeriod    = []byte(".")
)

func (node *ColName) Serialize(w Writer) error {
	return quoteName(w, node.Name)
}

// note: quoteName does not escape s. quoteName is indirectly
// called by builder.go, which checks that column/table names exist.
func quoteName(w io.Writer, s string) error {
	if _, err := w.Write(astBackquote); err != nil {
		return err
	}
	if _, err := io.WriteString(w, s); err != nil {
		return err
	}
	return t.Name
}

// loadColumns loads a table's columns from a database. MySQL
// specific.
func (t *Table) loadColumns(db *sql.DB) error {
	rows, err := db.Query("SHOW FULL COLUMNS FROM " + t.Name)
	if err != nil {
		return err
	}
// loadKeys loads a table's keys (indexes) from a database. MySQL
// specific.
func (t *Table) loadKeys(db *sql.DB) error {
	rows, err := db.Query("SHOW INDEX FROM " + t.Name)
	if err != nil {
		return err
	}
