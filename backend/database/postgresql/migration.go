package postgresql

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"strings"
)

func Run() error {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", DbUser, DbPassword, DbHost, DbPort, DbName, DbSslMode)
	m, err := migrate.New("file://backend/migrations", dbUrl)
	if err != nil {
		return err
	}
	defer m.Close()

	return applyMigrations(m)
	//return applyDownMigrations(m)
}

func applyMigrations(m *migrate.Migrate) error {
	err := m.Up()
	if err == nil || err == migrate.ErrNoChange {
		return nil
	}

	currVersion, _, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return nil
		}
		return err
	}

	// migration error can be handled by setting previous version and returning error
	prev := int(currVersion) - 1
	if err := m.Force(prev); err != nil {
		return err
	}

	return errors.New("database not up to date; migration(s) still pending")
}

func applyDownMigrations(m *migrate.Migrate) error {

	err := m.Steps(-1)
	if err == nil || err == migrate.ErrNoChange || strings.EqualFold(err.Error(), "file does not exist") {
		return nil
	}

	currVersion, _, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return nil
		}
		return err
	}

	// migration error can be handled by setting next version and returning error
	if err := m.Force(int(currVersion) + 1); err != nil {
		return err
	}

	return errors.New("forced version, Dirty state reset to false, migration not applied")
}

func fixDirtyMigration(m *migrate.Migrate) error {
	currVersion, _, err := m.Version()
	if err != nil || err == migrate.ErrNoChange {
		return err
	}

	err = m.Force(int(currVersion))
	if err == nil || err == migrate.ErrNoChange {
		return nil
	}
	return err
}
