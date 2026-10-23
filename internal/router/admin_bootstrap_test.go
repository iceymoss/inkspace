package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func bootstrapHTTPTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("INKSPACE_BOOTSTRAP_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set INKSPACE_BOOTSTRAP_TEST_MYSQL_DSN for isolated bootstrap HTTP tests")
	}
	parsed, err := mysqlDriver.ParseDSN(dsn)
	if err != nil || parsed.DBName != "inkspace_bootstrap_test" {
		t.Fatalf("bootstrap HTTP tests require inkspace_bootstrap_test: %v", err)
	}
	if os.Getenv("INKSPACE_TEST_ALLOW_DROP") != "1" {
		t.Skip("set INKSPACE_TEST_ALLOW_DROP=1 to reset isolated bootstrap fixtures")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.SystemInitialization{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&models.User{}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Where("1 = 1").Delete(&models.SystemInitialization{}).Error; err != nil {
		t.Fatal(err)
	}
	oldDB, oldConfig := database.DB, config.AppConfig
	database.DB = db
	config.AppConfig = &config.Config{}
	config.AppConfig.JWT.Secret = "bootstrap-test-user-token"
	config.AppConfig.JWT.AdminSecret = "bootstrap-test-admin-token"
	config.AppConfig.JWT.ExpireHours = 1
	t.Cleanup(func() {
		if err := db.Unscoped().Where("1 = 1").Delete(&models.User{}).Error; err != nil {
			t.Errorf("clean bootstrap test users: %v", err)
		}
		if err := db.Where("1 = 1").Delete(&models.SystemInitialization{}).Error; err != nil {
			t.Errorf("clean bootstrap test state: %v", err)
		}
		database.DB = oldDB
		config.AppConfig = oldConfig
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				t.Errorf("close bootstrap test DB: %v", err)
			}
		}
	})
	return db
}

type bootstrapHTTPResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		State            string `json:"state"`
		SecretConfigured bool   `json:"secret_configured"`
		Token            string `json:"token"`
	} `json:"data"`
}

func bootstrapHTTPCall(t *testing.T, engine *gin.Engine, method, path, body string) bootstrapHTTPResponse {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	var response bootstrapHTTPResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode HTTP bootstrap response (%d): %v", recorder.Code, err)
	}
	return response
}

func TestBootstrapHTTPSubmissionAndReplay(t *testing.T) {
	db := bootstrapHTTPTestDB(t)
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "isolated-bootstrap-secret")
	if err := db.Create(&models.SystemInitialization{ID: 1, State: models.InitializationPending}).Error; err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	engine := SetupAdminRouter()
	status := bootstrapHTTPCall(t, engine, http.MethodGet, "/api/admin/auth/bootstrap/status", "")
	if status.Code != 0 || status.Data.State != "pending" || !status.Data.SecretConfigured || status.Data.Token != "" {
		t.Fatalf("unexpected pending status: %+v", status)
	}
	body := `{"secret":"isolated-bootstrap-secret","username":"initial-admin","email":"initial@example.test","password":"a-strong-password","confirm_password":"a-strong-password"}`
	invalidRequests := map[string]string{
		"wrong secret":          `{"secret":"wrong","username":"initial-admin","email":"initial@example.test","password":"a-strong-password","confirm_password":"a-strong-password"}`,
		"password confirmation": `{"secret":"isolated-bootstrap-secret","username":"initial-admin","email":"initial@example.test","password":"a-strong-password","confirm_password":"different"}`,
	}
	for name, invalid := range invalidRequests {
		t.Run(name, func(t *testing.T) {
			if got := bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/bootstrap", invalid); got.Code == 0 {
				t.Fatal("accepted invalid initialization")
			}
		})
	}
	if got := bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/bootstrap", body); got.Code != 0 {
		t.Fatalf("initialization failed: %+v", got)
	}
	login := bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/login", `{"username":"initial-admin","password":"a-strong-password"}`)
	if login.Code != 0 || login.Data.Token == "" {
		t.Fatalf("first administrator cannot log in: %+v", login)
	}
	if got := bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/bootstrap", body); got.Code == 0 {
		t.Fatal("bootstrap replay succeeded")
	}
	status = bootstrapHTTPCall(t, engine, http.MethodGet, "/api/admin/auth/bootstrap/status", "")
	if status.Code != 0 || status.Data.State != "initialized" || status.Data.SecretConfigured {
		t.Fatalf("bootstrap reopened after completion: %+v", status)
	}
}

func TestBootstrapConcurrentSubmissionsCreateOnlyOneAdmin(t *testing.T) {
	db := bootstrapHTTPTestDB(t)
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "isolated-bootstrap-secret")
	if err := db.Create(&models.SystemInitialization{ID: 1, State: models.InitializationPending}).Error; err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	engine := SetupAdminRouter()
	var wg sync.WaitGroup
	start := make(chan struct{})
	outcomes := make(chan int, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			body := fmt.Sprintf(`{"secret":"isolated-bootstrap-secret","username":"admin-%d","email":"admin-%d@example.test","password":"a-strong-password","confirm_password":"a-strong-password"}`, i, i)
			outcomes <- bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/bootstrap", body).Code
		}(i)
	}
	close(start)
	wg.Wait()
	close(outcomes)
	var success int
	for code := range outcomes {
		if code == 0 {
			success++
		}
	}
	var admins int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&admins).Error; err != nil {
		t.Fatal(err)
	}
	if success != 1 || admins != 1 {
		t.Fatalf("successful requests=%d persisted admins=%d, want 1 each", success, admins)
	}
}

func TestBootstrapHTTPRejectsOccupiedUsernamesAndEmails(t *testing.T) {
	db := bootstrapHTTPTestDB(t)
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "isolated-bootstrap-secret")
	if err := db.Create(&models.SystemInitialization{ID: 1, State: models.InitializationPending}).Error; err != nil {
		t.Fatal(err)
	}
	password, err := utils.HashPassword("ordinary-password")
	if err != nil {
		t.Fatal(err)
	}
	ordinary := models.User{Username: "occupied-name", Email: "occupied@example.test", Password: password, Role: "user", Status: 1}
	if err := db.Create(&ordinary).Error; err != nil {
		t.Fatal(err)
	}
	engine := SetupAdminRouter()
	occupiedRequests := map[string]string{
		"username": `{"secret":"isolated-bootstrap-secret","username":"occupied-name","email":"new@example.test","password":"a-strong-password","confirm_password":"a-strong-password"}`,
		"email":    `{"secret":"isolated-bootstrap-secret","username":"new-name","email":"occupied@example.test","password":"a-strong-password","confirm_password":"a-strong-password"}`,
	}
	for name, body := range occupiedRequests {
		t.Run(name, func(t *testing.T) {
			got := bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/bootstrap", body)
			if got.Code == 0 || got.Message == "" {
				t.Fatalf("accepted occupied account field: %+v", got)
			}
		})
	}
	var admins int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&admins).Error; err != nil {
		t.Fatal(err)
	}
	if admins != 0 {
		t.Fatalf("created %d admins despite collisions", admins)
	}
	var state models.SystemInitialization
	if err := db.First(&state, 1).Error; err != nil || state.State != models.InitializationPending {
		t.Fatalf("changed state after collisions: %+v, %v", state, err)
	}
}

func TestPendingBootstrapKeepsPublicBlogAndUserRegistrationAvailable(t *testing.T) {
	db := bootstrapHTTPTestDB(t)
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "isolated-bootstrap-secret")
	if err := db.Create(&models.SystemInitialization{ID: 1, State: models.InitializationPending}).Error; err != nil {
		t.Fatal(err)
	}

	blog := SetupUserRouter()
	health := httptest.NewRecorder()
	blog.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("public blog health returned HTTP %d", health.Code)
	}

	registration := bootstrapHTTPCall(t, blog, http.MethodPost, "/api/register",
		`{"username":"ordinary-user","email":"ordinary@example.test","password":"ordinary-password","nickname":"Ordinary"}`)
	if registration.Code != 0 {
		t.Fatalf("ordinary registration failed while bootstrap pending: %+v", registration)
	}
	var ordinary models.User
	if err := db.Where("username = ?", "ordinary-user").First(&ordinary).Error; err != nil {
		t.Fatal(err)
	}
	if ordinary.Role != "user" {
		t.Fatalf("registered role=%q, want user", ordinary.Role)
	}

	admin := SetupAdminRouter()
	body := `{"secret":"isolated-bootstrap-secret","username":"initial-admin","email":"initial@example.test","password":"a-strong-password","confirm_password":"a-strong-password"}`
	if got := bootstrapHTTPCall(t, admin, http.MethodPost, "/api/admin/auth/bootstrap", body); got.Code != 0 {
		t.Fatalf("ordinary registration blocked later bootstrap: %+v", got)
	}
}

func TestBootstrapRollsBackAdministratorWhenCompletionFails(t *testing.T) {
	db := bootstrapHTTPTestDB(t)
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "isolated-bootstrap-secret")
	if err := db.Create(&models.SystemInitialization{ID: 1, State: models.InitializationPending}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TRIGGER reject_bootstrap_completion BEFORE UPDATE ON system_initializations
		FOR EACH ROW SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'forced completion failure'`).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Exec("DROP TRIGGER IF EXISTS reject_bootstrap_completion").Error; err != nil {
			t.Errorf("drop bootstrap failure trigger: %v", err)
		}
	})

	engine := SetupAdminRouter()
	body := `{"secret":"isolated-bootstrap-secret","username":"initial-admin","email":"initial@example.test","password":"a-strong-password","confirm_password":"a-strong-password"}`
	if got := bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/bootstrap", body); got.Code == 0 {
		t.Fatal("bootstrap succeeded despite forced completion failure")
	}

	var admins int64
	if err := db.Unscoped().Model(&models.User{}).Where("role = ?", "admin").Count(&admins).Error; err != nil {
		t.Fatal(err)
	}
	var state models.SystemInitialization
	if err := db.First(&state, 1).Error; err != nil {
		t.Fatal(err)
	}
	if admins != 0 || state.State != models.InitializationPending || state.CompletedAt != nil {
		t.Fatalf("partial bootstrap persisted: admins=%d state=%q completed_at=%v", admins, state.State, state.CompletedAt)
	}
}

func TestBootstrapHTTPRejectsMissingSecretAndClosedState(t *testing.T) {
	db := bootstrapHTTPTestDB(t)
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "")
	if err := db.Create(&models.SystemInitialization{ID: 1, State: models.InitializationPending}).Error; err != nil {
		t.Fatal(err)
	}
	engine := SetupAdminRouter()
	body := `{"secret":"any-secret","username":"first-admin","email":"first@example.test","password":"a-strong-password","confirm_password":"a-strong-password"}`
	if got := bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/bootstrap", body); got.Code == 0 {
		t.Fatal("created an administrator without a configured secret")
	}
	status := bootstrapHTTPCall(t, engine, http.MethodGet, "/api/admin/auth/bootstrap/status", "")
	if status.Code != 0 || status.Data.State != "pending" || status.Data.SecretConfigured {
		t.Fatalf("missing-secret status incorrect: %+v", status)
	}
	if err := db.Model(&models.SystemInitialization{}).Where("id = 1").Update("state", models.InitializationCompleted).Error; err != nil {
		t.Fatal(err)
	}
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "any-secret")
	status = bootstrapHTTPCall(t, engine, http.MethodGet, "/api/admin/auth/bootstrap/status", "")
	if status.Code != 0 || status.Data.State != "unavailable" {
		t.Fatalf("completed state without active admin reopened: %+v", status)
	}
	if got := bootstrapHTTPCall(t, engine, http.MethodPost, "/api/admin/auth/bootstrap", body); got.Code == 0 {
		t.Fatal("reopened completed initialization without an admin")
	}
}

func TestBootstrapStatusFailsClosedForHistoricalAndUnavailableAdministrators(t *testing.T) {
	db := bootstrapHTTPTestDB(t)
	t.Setenv("INKSPACE_BOOTSTRAP_SECRET", "isolated-bootstrap-secret")
	engine := SetupAdminRouter()

	tests := []struct {
		name      string
		state     string
		hasAdmin  bool
		status    int
		deleted   bool
		wantState string
	}{
		{name: "legacy active administrator", state: models.InitializationLegacyClosed, hasAdmin: true, status: 1, wantState: "initialized"},
		{name: "legacy disabled administrator", state: models.InitializationLegacyClosed, hasAdmin: true, status: 0, wantState: "unavailable"},
		{name: "legacy deleted administrator", state: models.InitializationLegacyClosed, hasAdmin: true, status: 1, deleted: true, wantState: "unavailable"},
		{name: "pending with deleted historical administrator", state: models.InitializationPending, hasAdmin: true, status: 1, deleted: true, wantState: "unavailable"},
		{name: "completed without administrator", state: models.InitializationCompleted, wantState: "unavailable"},
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
				user := models.User{Username: "historical-admin", Email: "historical@example.test", Password: "unused", Role: "admin", Status: 1}
				if err := db.Create(&user).Error; err != nil {
					t.Fatal(err)
				}
				if tt.status == 0 {
					if err := db.Model(&user).Update("status", 0).Error; err != nil {
						t.Fatal(err)
					}
				}
				if tt.deleted {
					if err := db.Delete(&user).Error; err != nil {
						t.Fatal(err)
					}
				}
			}

			got := bootstrapHTTPCall(t, engine, http.MethodGet, "/api/admin/auth/bootstrap/status", "")
			if got.Code != 0 || got.Data.State != tt.wantState {
				t.Fatalf("status=%+v, want state %q", got, tt.wantState)
			}
		})
	}
}

func TestBootstrapStatusFailsClosedWithoutDatabase(t *testing.T) {
	old := database.DB
	database.DB = nil
	t.Cleanup(func() { database.DB = old })

	engine := SetupAdminRouter()
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/admin/auth/bootstrap/status", nil))
	var body struct {
		Code int `json:"code"`
		Data struct {
			State string `json:"state"`
		} `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code == 0 || body.Data.State == "pending" {
		t.Fatalf("database failure opened bootstrap: %s", response.Body.String())
	}
}
