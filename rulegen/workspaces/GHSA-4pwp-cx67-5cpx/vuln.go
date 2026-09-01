package main

	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/VividCortex/mysqlerr"

	"github.com/grafana/grafana/pkg/setting"

	"github.com/go-sql-driver/mysql"
	tsdb.RegisterTsdbQueryEndpoint("mysql", newMysqlQueryEndpoint)
}

func newMysqlQueryEndpoint(datasource *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
	logger := log.New("tsdb.mysql")

	if strings.HasPrefix(datasource.Url, "/") {
		protocol = "unix"
	}
	cnnstr := fmt.Sprintf("%s:%s@%s(%s)/%s?collation=utf8mb4_unicode_ci&parseTime=true&loc=UTC&allowNativePasswords=true",
		datasource.User,
		datasource.DecryptedPassword(),
		protocol,
		datasource.Url,
		datasource.Database,
	)

	tlsConfig, err := datasource.GetTLSConfig()
