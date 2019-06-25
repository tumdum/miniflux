// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package database // import "miniflux.app/database"

import (
	"database/sql"
	"strconv"

	"miniflux.app/logger"
)

const schemaVersion = 23

// Migrate executes database migrations.
func Migrate(db *sql.DB) error {
	var currentVersion int
	db.QueryRow(`select version from schema_version`).Scan(&currentVersion)

	logger.Debug("Current schema version:", currentVersion)
	logger.Debug("Latest schema version:", schemaVersion)

	for version := currentVersion + 1; version <= schemaVersion; version++ {
		logger.Debug("Migrating to version: %v", version)

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		rawSQL := SqlMap["schema_version_"+strconv.Itoa(version)]
		// fmt.Println(rawSQL)
		_, err = tx.Exec(rawSQL)
		if err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec(`delete from schema_version`); err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec(`insert into schema_version (version) values($1)`, version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
