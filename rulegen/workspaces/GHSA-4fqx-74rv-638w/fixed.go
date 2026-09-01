package main


	_, err = psql.Insert("build_resource_config_version_outputs").
		Columns("resource_id", "build_id", "version_md5", "name").
		Values(resource.ID(), strconv.Itoa(b.id), sq.Expr("md5(?)", versionJSON), outputName).
		Suffix("ON CONFLICT DO NOTHING").
		RunWith(tx).
		Exec()

	_, err = psql.Insert("build_resource_config_version_inputs").
		Columns("build_id", "resource_id", "version_md5", "name").
		Values(buildID, input.ResourceID, sq.Expr("md5(?)", versionJSON), input.Name).
		Suffix("ON CONFLICT DO NOTHING").
		RunWith(tx).
		Exec()
import (
	"database/sql"
	"encoding/json"
	"time"

	"code.cloudfoundry.org/lager"
		Where(sq.Eq{
			"v.resource_config_scope_id": r.id,
		}).
		Where(sq.Expr("v.version_md5 = md5(?)", versionByte)).
		RunWith(r.conn).
		QueryRow()

