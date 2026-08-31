package main


	exists, err := mg.x.IsTableExist(new(MigrationLog))
	if err != nil {
		return nil, err
	}

	if !exists {
		return logMap, nil
	}
}

func (mg *Migrator) Start() error {
	mg.Logger.Info("Starting DB migration")

	logMap, err := mg.GetMigrationLog()
	if err != nil {
			_, err = sess.Insert(&record)
			return err
		})

		if err != nil {
			return err
		}
	}


	"github.com/go-sql-driver/mysql"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/infra/localcache"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	}

	sqlog.Info("Connecting to DB", "dbtype", ss.dbCfg.Type)
	engine, err := xorm.NewEngine(ss.dbCfg.Type, connectionString)
	if err != nil {
		return nil, err
