package main


	_, err = psql.Insert("build_resource_config_version_outputs").
		Columns("resource_id", "build_id", "version_md5", "name").
		Values(resource.ID(), strconv.Itoa(b.id), sq.Expr(fmt.Sprintf("md5('%s')", versionJSON)), sq.Expr(fmt.Sprintf("'%s'", outputName))).
		Suffix("ON CONFLICT DO NOTHING").
		RunWith(tx).
		Exec()

	_, err = psql.Insert("build_resource_config_version_inputs").
		Columns("build_id", "resource_id", "version_md5", "name").
		Values(buildID, input.ResourceID, sq.Expr(fmt.Sprintf("md5('%s')", versionJSON)), input.Name).
		Suffix("ON CONFLICT DO NOTHING").
		RunWith(tx).
		Exec()
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"code.cloudfoundry.org/lager"
		Where(sq.Eq{
			"v.resource_config_scope_id": r.id,
		}).
		Where(sq.Expr(fmt.Sprintf("v.version_md5 = md5('%s')", versionByte))).
		RunWith(r.conn).
		QueryRow()

