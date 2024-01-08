package database

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/HarveyJhuang1010/blockhw/internal/config"
)

const (
	checkTestDBStmt  = "SELECT * FROM pg_database WHERE datname = \"%s\";"
	createTestDBStmt = "CREATE DATABASE \"%s\";"
	dropDatabaseStmt = "DROP DATABASE IF EXISTS \"%s\";"
)

// InitTestPostgresSQL init test database for testing
func InitTestPostgresSQL(t *testing.T, cfg *config.DatabaseConfig, testDBName string) (*DB, error) {
	// Init Database
	Initialize(cfg)

	// Check test db exist
	isExist, err := checkDBExist(t, testDBName)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	if !isExist {
		// Create test db
		t.Logf(createTestDBStmt, testDBName)
		if err := createDB(t, cfg, testDBName); err != nil {
			t.Error(err)
			return nil, err
		}
	}

	// close local db connection
	Finalize()

	// conn to test db
	cfg.DBName = testDBName
	Initialize(cfg)

	db := GetDB()

	return db, nil
}

// checkDBExist check database exist
func checkDBExist(t *testing.T, testDBName string) (bool, error) {
	var (
		dbCli   = GetDB()
		isExist bool
	)

	t.Logf(checkTestDBStmt, testDBName)
	rs := dbCli.Raw(fmt.Sprintf(checkTestDBStmt, testDBName))
	if err := rs.Error; err != nil {
		t.Error(err)
		return isExist, err
	}

	// if not create test db
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		return isExist, nil
	}

	return true, nil
}

// createDB create database
func createDB(t *testing.T, cfg *config.DatabaseConfig, testDBName string) error {
	// #nosec
	cmd := exec.Command("psql",
		"-h", cfg.Host,
		"-p", fmt.Sprintf("%d", cfg.Port),
		"-d", cfg.DBName,
		"-U", cfg.User,
		"-c", fmt.Sprintf(createTestDBStmt, testDBName),
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("err=%v msg=%s", err, out)
		return err
	}

	return nil
}

// RemoveTestDB drop test database and disconnect
func RemoveTestDB(t *testing.T, cfg *config.DatabaseConfig, testDBName string) error {
	Finalize()

	// #nosec
	cmd := exec.Command("psql",
		"-h", cfg.Host,
		"-p", fmt.Sprintf("%d", cfg.Port),
		"-d", cfg.DBName,
		"-U", cfg.User,
		"-c", fmt.Sprintf(dropDatabaseStmt, testDBName),
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("err=%v msg=%s cmd=%s", err, out, cmd.String())
		return err
	}

	return nil
}
