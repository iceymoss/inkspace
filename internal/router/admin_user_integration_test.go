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
	"time"

	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/iceymoss/inkspace/internal/config"
	"github.com/iceymoss/inkspace/internal/database"
	"github.com/iceymoss/inkspace/internal/models"
	"github.com/iceymoss/inkspace/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func adminUserTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("INKSPACE_ADMIN_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set INKSPACE_ADMIN_TEST_MYSQL_DSN to run admin account integration tests")
	}
	parsed, err := mysqlDriver.ParseDSN(dsn)
	if err != nil || parsed.DBName != "inkspace_admin_guard_test" {
		t.Fatalf("admin integration tests require database inkspace_admin_guard_test: %v", err)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Article{}, &models.Work{}, &models.Comment{}, &models.UserFollow{}, &models.Favorite{}); err != nil {
		t.Fatalf("migrate admin test tables: %v", err)
	}
	previousDB, previousConfig := database.DB, config.AppConfig
	database.DB = db
	config.AppConfig = &config.Config{}
	config.AppConfig.JWT.Secret = "admin-user-test-secret"
	config.AppConfig.JWT.AdminSecret = "admin-user-test-admin-secret"
	config.AppConfig.JWT.ExpireHours = 1
	t.Cleanup(func() {
		database.DB = previousDB
		config.AppConfig = previousConfig
	})
	return db
}

type adminUserTestResponse struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func adminUserRequest(t *testing.T, engine *gin.Engine, method, path, token, body string) adminUserTestResponse {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	var response adminUserTestResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode %s %s (%d): %v", method, path, recorder.Code, err)
	}
	return response
}

func adminUserFixture(t *testing.T, db *gorm.DB, name, role string, status int) models.User {
	t.Helper()
	password, err := utils.HashPassword("test-password")
	if err != nil {
		t.Fatal(err)
	}
	username := fmt.Sprintf("%s-%d", name, time.Now().UnixNano())
	user := models.User{Username: username, Email: username + "@example.test", Password: password, Role: role, Status: status}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create %s: %v", name, err)
	}
	if status == 0 {
		if err := db.Model(&user).Update("status", 0).Error; err != nil {
			t.Fatalf("disable %s: %v", name, err)
		}
	}
	t.Cleanup(func() {
		if err := db.Unscoped().Delete(&models.User{}, user.ID).Error; err != nil {
			t.Errorf("clean up %s: %v", username, err)
		}
	})
	return user
}

func adminUserLogin(t *testing.T, engine *gin.Engine, username string) string {
	t.Helper()
	response := adminUserRequest(t, engine, http.MethodPost, "/api/admin/auth/login", "", fmt.Sprintf(`{"username":%q,"password":"test-password"}`, username))
	if response.Code != 0 || response.Data.Token == "" {
		t.Fatalf("admin login for %s returned code %d", username, response.Code)
	}
	return response.Data.Token
}

func TestAdminUserLastLoginCapableAdmin(t *testing.T) {
	db := adminUserTestDB(t)
	gin.SetMode(gin.TestMode)
	engine := SetupAdminRouter()
	admin := adminUserFixture(t, db, "sole-admin", "admin", 1)
	adminUserFixture(t, db, "inactive-admin", "admin", 0)
	deleted := adminUserFixture(t, db, "deleted-admin", "admin", 1)
	if err := db.Delete(&deleted).Error; err != nil {
		t.Fatal(err)
	}
	token := adminUserLogin(t, engine, admin.Username)
	for _, action := range []struct {
		name, method, path, body string
	}{
		{"demote", http.MethodPut, fmt.Sprintf("/api/admin/users/%d/role", admin.ID), `{"role":"user"}`},
		{"disable", http.MethodPut, fmt.Sprintf("/api/admin/users/%d/status", admin.ID), `{"status":0}`},
		{"delete", http.MethodDelete, fmt.Sprintf("/api/admin/users/%d", admin.ID), ""},
	} {
		t.Run(action.name, func(t *testing.T) {
			response := adminUserRequest(t, engine, action.method, action.path, token, action.body)
			if response.Code == 0 {
				t.Fatal("removed the only login-capable admin")
			}
			if got := adminUserRequest(t, engine, http.MethodGet, "/api/admin/auth/profile", token, ""); got.Code != 0 {
				t.Fatalf("admin lost access after rejected operation: %d", got.Code)
			}
			adminUserLogin(t, engine, admin.Username)
		})
	}
}

func TestAdminUserChangesRevokeOldTokens(t *testing.T) {
	db := adminUserTestDB(t)
	gin.SetMode(gin.TestMode)
	engine := SetupAdminRouter()
	keeper := adminUserFixture(t, db, "keeper", "admin", 1)
	keeperToken := adminUserLogin(t, engine, keeper.Username)
	for _, action := range []struct {
		name, method, path, body string
	}{
		{"demote", http.MethodPut, "/role", `{"role":"user"}`},
		{"disable", http.MethodPut, "/status", `{"status":0}`},
		{"delete", http.MethodDelete, "", ""},
		{"keep-active", http.MethodPut, "/status", `{"status":1}`},
		{"keep-role", http.MethodPut, "/role", `{"role":"admin"}`},
	} {
		t.Run(action.name, func(t *testing.T) {
			user := adminUserFixture(t, db, "target-"+action.name, "admin", 1)
			oldToken := adminUserLogin(t, engine, user.Username)
			path := fmt.Sprintf("/api/admin/users/%d%s", user.ID, action.path)
			if got := adminUserRequest(t, engine, action.method, path, keeperToken, action.body); got.Code != 0 {
				t.Fatalf("change rejected with another active admin: %d", got.Code)
			}
			if got := adminUserRequest(t, engine, http.MethodGet, "/api/admin/auth/profile", oldToken, ""); (got.Code == 0) != (action.name == "keep-active" || action.name == "keep-role") {
				t.Fatalf("unexpected old token access after %s: code %d", action.name, got.Code)
			}
			if got := adminUserRequest(t, engine, http.MethodGet, "/api/admin/auth/profile", keeperToken, ""); got.Code != 0 {
				t.Fatalf("unchanged admin lost access: %d", got.Code)
			}
		})
	}
}

func TestAdminUserStatusAndRolePromotions(t *testing.T) {
	db := adminUserTestDB(t)
	gin.SetMode(gin.TestMode)
	engine := SetupAdminRouter()
	keeper := adminUserFixture(t, db, "promotion-keeper", "admin", 1)
	actor := adminUserLogin(t, engine, keeper.Username)
	inactive := adminUserFixture(t, db, "promotion-inactive", "admin", 0)
	ordinary := adminUserFixture(t, db, "promotion-user", "user", 1)
	for _, action := range []struct {
		id         uint
		path, body string
	}{
		{inactive.ID, "/status", `{"status":1}`},
		{ordinary.ID, "/role", `{"role":"admin"}`},
	} {
		path := fmt.Sprintf("/api/admin/users/%d%s", action.id, action.path)
		if got := adminUserRequest(t, engine, http.MethodPut, path, actor, action.body); got.Code != 0 {
			t.Fatalf("promotion rejected: %d", got.Code)
		}
	}
	for _, name := range []string{inactive.Username, ordinary.Username} {
		token := adminUserLogin(t, engine, name)
		if got := adminUserRequest(t, engine, http.MethodGet, "/api/admin/auth/profile", token, ""); got.Code != 0 {
			t.Fatalf("promoted %s cannot access admin profile: %d", name, got.Code)
		}
	}
}

func TestAdminUserAuthDeniesRevokedAccountOnAnotherRoute(t *testing.T) {
	db := adminUserTestDB(t)
	gin.SetMode(gin.TestMode)
	engine := SetupAdminRouter()
	keeper := adminUserFixture(t, db, "alternate-keeper", "admin", 1)
	target := adminUserFixture(t, db, "alternate-target", "admin", 1)
	keeperToken := adminUserLogin(t, engine, keeper.Username)
	targetToken := adminUserLogin(t, engine, target.Username)
	path := fmt.Sprintf("/api/admin/users/%d/role", target.ID)
	if got := adminUserRequest(t, engine, http.MethodPut, path, keeperToken, `{"role":"user"}`); got.Code != 0 {
		t.Fatalf("demotion failed: %d", got.Code)
	}
	if got := adminUserRequest(t, engine, http.MethodGet, "/api/admin/knowledge/overview", targetToken, ""); got.Code != 403 {
		t.Fatalf("revoked admin accessed another protected route: %d", got.Code)
	}
}

func TestConcurrentAdminUserChangesPreserveAnAdmin(t *testing.T) {
	db := adminUserTestDB(t)
	gin.SetMode(gin.TestMode)
	engine := SetupAdminRouter()
	for _, action := range []struct {
		name, method, path, body string
	}{
		{"demote", http.MethodPut, "/role", `{"role":"user"}`},
		{"disable", http.MethodPut, "/status", `{"status":0}`},
		{"delete", http.MethodDelete, "", ""},
		{"demote-and-disable", http.MethodPut, "/role", `{"role":"user"}`},
	} {
		t.Run(action.name, func(t *testing.T) {
			first := adminUserFixture(t, db, "first-"+action.name, "admin", 1)
			second := adminUserFixture(t, db, "second-"+action.name, "admin", 1)
			tokens := []string{adminUserLogin(t, engine, first.Username), adminUserLogin(t, engine, second.Username)}
			ids := []uint{first.ID, second.ID}
			start := make(chan struct{})
			results := make(chan struct {
				index int
				code  int
			}, 2)
			var wg sync.WaitGroup
			for i := 0; i < 2; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					<-start
					path := fmt.Sprintf("/api/admin/users/%d%s", ids[i], action.path)
					method, body := action.method, action.body
					if action.name == "demote-and-disable" && i == 1 {
						path = fmt.Sprintf("/api/admin/users/%d/status", ids[i])
						method, body = http.MethodPut, `{"status":0}`
					}
					results <- struct {
						index int
						code  int
					}{i, adminUserRequest(t, engine, method, path, tokens[i], body).Code}
				}(i)
			}
			close(start)
			wg.Wait()
			close(results)
			var succeeded int
			outcomes := [2]int{}
			for result := range results {
				outcomes[result.index] = result.code
				if result.code == 0 {
					succeeded++
				}
			}
			if succeeded != 1 {
				t.Fatalf("successful changes = %d, want 1", succeeded)
			}
			for i, id := range ids {
				var user models.User
				if err := db.Unscoped().First(&user, id).Error; err != nil {
					t.Fatal(err)
				}
				stillAdmin := user.Role == "admin" && user.Status == 1 && !user.DeletedAt.Valid
				if (outcomes[i] == 0) == stillAdmin {
					t.Fatalf("account %d: response code %d disagrees with persisted role=%s status=%d deleted=%v", id, outcomes[i], user.Role, user.Status, user.DeletedAt.Valid)
				}
			}
			var active int64
			if err := db.Model(&models.User{}).Where("role = ? AND status = ?", "admin", 1).Count(&active).Error; err != nil {
				t.Fatal(err)
			}
			if active != 1 {
				t.Fatalf("login-capable administrators = %d, want 1", active)
			}
		})
	}
}
