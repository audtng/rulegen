package main

	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/VividCortex/mysqlerr"
	"github.com/grafana/grafana/pkg/setting"

	"github.com/go-sql-driver/mysql"
	tsdb.RegisterTsdbQueryEndpoint("mysql", newMysqlQueryEndpoint)
}

func characterEscape(s string, escapeChar string) string {
	return strings.Replace(s, escapeChar, url.QueryEscape(escapeChar), -1)
}

func newMysqlQueryEndpoint(datasource *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	logger := log.New("tsdb.mysql")

	if strings.HasPrefix(datasource.Url, "/") {
		protocol = "unix"
	}

	cnnstr := fmt.Sprintf("%s:%s@%s(%s)/%s?collation=utf8mb4_unicode_ci&parseTime=true&loc=UTC&allowNativePasswords=true",
		characterEscape(datasource.User, ":"),
		characterEscape(datasource.DecryptedPassword(), "@"),
		protocol,
		characterEscape(datasource.Url, ")"),
		characterEscape(datasource.Database, "?"),
	)

	tlsConfig, err := datasource.GetTLSConfig()
