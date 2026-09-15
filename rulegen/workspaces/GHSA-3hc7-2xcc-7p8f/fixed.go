package main

	"bytes"
	"fmt"
	"io"
	"strings"
)

// Instructions for creating new types: If a type needs to satisfy an
func (*TableName) simpleTableExpr() {}
func (*Subquery) simpleTableExpr()  {}

// TableName represents a table name.
type TableName struct {
	Name, Qualifier string
}
}

var (
	astBackquoteStr       = "`"
	astDoubleBackquoteStr = "``"
	astBackquote          = []byte(astBackquoteStr)
	astPeriod             = []byte(".")
)

func (node *ColName) Serialize(w Writer) error {
	return quoteName(w, node.Name)
}

// note: quoteName escapes any backquote (`) characters in s. quoteName is indirectly
// called by builder.go, which checks that column/table names exist.
func quoteName(w io.Writer, s string) error {
	if _, err := w.Write(astBackquote); err != nil {
		return err
	}
	s = strings.ReplaceAll(s, astBackquoteStr, astDoubleBackquoteStr)
	if _, err := io.WriteString(w, s); err != nil {
		return err
	}
	return t.Name
}

// Enclose a name in backquotes and escape any internal backquotes.
func quoteNameStr(s string) string {
	return astBackquoteStr + strings.ReplaceAll(s, astBackquoteStr, astDoubleBackquoteStr) + astBackquoteStr
}

// loadColumns loads a table's columns from a database. MySQL
// specific.
func (t *Table) loadColumns(db *sql.DB) error {
	rows, err := db.Query("SHOW FULL COLUMNS FROM " + quoteNameStr(t.Name))
	if err != nil {
		return err
	}
// loadKeys loads a table's keys (indexes) from a database. MySQL
// specific.
func (t *Table) loadKeys(db *sql.DB) error {
	rows, err := db.Query("SHOW INDEX FROM " + quoteNameStr(t.Name))
	if err != nil {
		return err
	}
