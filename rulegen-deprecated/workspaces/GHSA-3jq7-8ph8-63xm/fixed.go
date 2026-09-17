package main


	"github.com/go-sql-driver/mysql"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/infra/fs"
	"github.com/grafana/grafana/pkg/infra/localcache"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	}

	sqlog.Info("Connecting to DB", "dbtype", ss.dbCfg.Type)
	if ss.dbCfg.Type == migrator.SQLITE && strings.HasPrefix(connectionString, "file:") {
		exists, err := fs.Exists(ss.dbCfg.Path)
		if err != nil {
			return nil, errutil.Wrapf(err, "can't check for existence of %q", ss.dbCfg.Path)
		}

		const perms = 0640
		if !exists {
			ss.log.Info("Creating SQLite database file", "path", ss.dbCfg.Path)
			f, err := os.OpenFile(ss.dbCfg.Path, os.O_CREATE|os.O_RDWR, perms)
			if err != nil {
				return nil, errutil.Wrapf(err, "failed to create SQLite database file %q", ss.dbCfg.Path)
			}
			if err := f.Close(); err != nil {
				return nil, errutil.Wrapf(err, "failed to create SQLite database file %q", ss.dbCfg.Path)
			}
		} else {
			fi, err := os.Lstat(ss.dbCfg.Path)
			if err != nil {
				return nil, errutil.Wrapf(err, "failed to stat SQLite database file %q", ss.dbCfg.Path)
			}
			m := fi.Mode() & os.ModePerm
			if m|perms != perms {
				ss.log.Warn("SQLite database file has broader permissions than it should",
					"path", ss.dbCfg.Path, "mode", m, "expected", os.FileMode(perms))
			}
		}
	}
	engine, err := xorm.NewEngine(ss.dbCfg.Type, connectionString)
	if err != nil {
		return nil, err

	exists, err := mg.x.IsTableExist(new(MigrationLog))
	if err != nil {
		return nil, errutil.Wrap("failed to check table existence", err)
	}
	if !exists {
		return logMap, nil
	}
}

func (mg *Migrator) Start() error {
	mg.Logger.Info("Starting DB migrations")

	logMap, err := mg.GetMigrationLog()
	if err != nil {
			_, err = sess.Insert(&record)
			return err
		})
		if err != nil {
			return errutil.Wrap("migration failed", err)
		}
	}

