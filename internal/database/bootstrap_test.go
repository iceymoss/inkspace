package database

import (
	"fmt"
	"os"
	"strings"
	"testing"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/iceymoss/inkspace/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Bootstrap startup tests use only a dedicated, explicitly named disposable schema.
func bootstrapTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("INKSPACE_BOOTSTRAP_STARTUP_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set INKSPACE_BOOTSTRAP_STARTUP_TEST_MYSQL_DSN to run bootstrap MySQL tests")
	}
	parsed, err := mysqlDriver.ParseDSN(dsn)
	if err != nil || parsed.DBName != "inkspace_bootstrap_startup_test" {
		t.Fatalf("bootstrap tests require inkspace_bootstrap_startup_test: %v", err)
	}
	if os.Getenv("INKSPACE_TEST_ALLOW_DROP") != "1" {
		t.Skip("set INKSPACE_TEST_ALLOW_DROP=1 to reset the dedicated bootstrap test schema")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("close bootstrap startup test DB: %v", err)
		}
	})
	return db
}

func resetBootstrapSchema(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Connection(func(conn *gorm.DB) error {
		if err := conn.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
			return err
		}
		var tables []string
		if err := conn.Raw(`SELECT table_name FROM information_schema.tables
			WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE'`).Scan(&tables).Error; err != nil {
			return err
		}
		for _, table := range tables {
			quoted := strings.ReplaceAll(table, "`", "``")
			if err := conn.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", quoted)).Error; err != nil {
				return err
			}
		}
		return conn.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	}); err != nil {
		t.Fatal(err)
	}
}

func TestBootstrapStartupRecognizesFreshAndLegacySchemas(t *testing.T) {
	db := bootstrapTestDB(t)
	old := DB
	DB = db
	t.Cleanup(func() { DB = old })
	tests := []struct {
		name   string
		legacy bool
	}{
		{name: "fresh schema"},
		{name: "legacy schema", legacy: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetBootstrapSchema(t, db)
			if tt.legacy {
				if err := db.Exec("CREATE TABLE users (id BIGINT UNSIGNED PRIMARY KEY)").Error; err != nil {
					t.Fatal(err)
				}
			}
			if err := PrepareInitialization(); err != nil {
				t.Fatal(err)
			}
			var state string
			if err := db.Table("system_initializations").Select("state").Where("id = 1").Scan(&state).Error; err != nil {
				t.Fatal(err)
			}
			want := "pending"
			if tt.legacy {
				want = "legacy_closed"
			}
			if state != want {
				t.Fatalf("state=%q, want %q", state, want)
			}
			if err := PrepareInitialization(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBootstrapStartupOrderUsesOnePersistentState(t *testing.T) {
	db := bootstrapTestDB(t)
	old := DB
	DB = db
	t.Cleanup(func() { DB = old })

	orders := [][]string{
		{"blog", "admin", "scheduler"},
		{"admin", "scheduler", "blog"},
		{"scheduler", "blog", "admin"},
	}
	for _, order := range orders {
		t.Run(strings.Join(order, "-"), func(t *testing.T) {
			resetBootstrapSchema(t, db)
			for _, program := range order {
				if err := PrepareInitialization(); err != nil {
					t.Fatalf("%s startup: %v", program, err)
				}
			}
			var states int64
			if err := db.Model(&models.SystemInitialization{}).Where("id = ? AND state = ?", 1, models.InitializationPending).Count(&states).Error; err != nil {
				t.Fatal(err)
			}
			if states != 1 {
				t.Fatalf("persistent pending states=%d, want 1", states)
			}
		})
	}
}

func TestBootstrapStartupKeepsHistoricalAdministratorsClosed(t *testing.T) {
	db := bootstrapTestDB(t)
	resetBootstrapSchema(t, db)
	if err := db.AutoMigrate(&models.SystemInitialization{}, &models.User{}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		state     string
		hasAdmin  bool
		status    int
		deleted   bool
		wantError bool
	}{
		{name: "legacy active admin", state: models.InitializationLegacyClosed, hasAdmin: true, status: 1},
		{name: "legacy disabled admin", state: models.InitializationLegacyClosed, hasAdmin: true},
		{name: "legacy deleted admin", state: models.InitializationLegacyClosed, hasAdmin: true, status: 1, deleted: true},
		{name: "completed without admin", state: models.InitializationCompleted},
		{name: "pending historical admin", state: models.InitializationPending, hasAdmin: true, status: 1, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Unscoped().Where("1 = 1").Delete(&models.User{}).Error; err != nil {
				t.Fatal(err)
			}
			if err := db.Where("1 = 1").Delete(&models.SystemInitialization{}).Error; err != nil {
				t.Fatal(err)
			}
			if err := db.Create(&models.SystemInitialization{ID: 1, State: tt.state}).Error; err != nil {
				t.Fatal(err)
			}
			if tt.hasAdmin {
				admin := models.User{Username: "historical-admin", Email: "historical@example.test", Password: "unused", Role: "admin", Status: 1}
				if err := db.Create(&admin).Error; err != nil {
					t.Fatal(err)
				}
				if tt.status == 0 {
					if err := db.Model(&admin).Update("status", 0).Error; err != nil {
						t.Fatal(err)
					}
				}
				if tt.deleted {
					if err := db.Delete(&admin).Error; err != nil {
						t.Fatal(err)
					}
				}
			}

			err := checkInitializationConsistency(db)
			if (err != nil) != tt.wantError {
				t.Fatalf("consistency error=%v, wantError=%v", err, tt.wantError)
			}
			var state models.SystemInitialization
			if err := db.First(&state, 1).Error; err != nil {
				t.Fatal(err)
			}
			if state.State != tt.state {
				t.Fatalf("state changed to %q, want %q", state.State, tt.state)
			}
		})
	}
}

func TestBootstrapStartupNeverRepairsMissingOrUnknownState(t *testing.T) {
	db := bootstrapTestDB(t)
	old := DB
	DB = db
	t.Cleanup(func() { DB = old })
	if err := db.AutoMigrate(&models.SystemInitialization{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Where("1 = 1").Delete(&models.SystemInitialization{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := PrepareInitialization(); err == nil {
		t.Fatal("reopened initialization with an orphaned state table")
	}
	if err := db.Create(&models.SystemInitialization{ID: 1, State: "unknown"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := PrepareInitialization(); err == nil {
		t.Fatal("accepted an unknown initialization state")
	}
}

func TestBootstrapMigrationUsesLockedConnectionWithSingleConnection(t *testing.T) {
	db := bootstrapTestDB(t)
	old := DB
	DB = db
	t.Cleanup(func() { DB = old })
	resetBootstrapSchema(t, db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := WithInitializationLock(func(conn *gorm.DB) error {
		if err := prepareInitialization(conn); err != nil {
			return err
		}
		if err := migrateDB(conn); err != nil {
			return err
		}
		return checkInitializationConsistency(conn)
	}); err != nil {
		t.Fatalf("single-connection startup migration failed: %v", err)
	}

	var state models.SystemInitialization
	if err := db.First(&state, 1).Error; err != nil {
		t.Fatal(err)
	}
	if state.State != models.InitializationPending {
		t.Fatalf("state=%q, want %q", state.State, models.InitializationPending)
	}
}
