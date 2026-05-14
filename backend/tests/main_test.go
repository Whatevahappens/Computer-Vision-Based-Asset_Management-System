package tests

import (
	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/model"
	"backend/internal/router"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testRouter      *gin.Engine
	testDB          *gorm.DB
	custodianToken  string
	accountantToken string
	employeeToken   string
	adminToken      string
	testLocationID  string
	testAssetID     string
	testSessionID   string
	testDeptID      string
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// SQLite in-memory
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect sqlite: %v", err)
	}

	// Migrate all models
	testDB.AutoMigrate(
		&model.User{},
		&model.Department{},
		&model.Location{},
		&model.AssetModel{},
		&model.Asset{},
		&model.AssetHistory{},
		&model.Document{},
		&model.AuditSession{},
		&model.AuditFinding{},
		&model.AuditEvidence{},
		&model.AuditSummary{},
		&model.Notification{},
	)

	// Override the global DB used by repositories
	// You need to make database.DB assignable or add a setter
	// For now, assuming database.DB is a package-level var:
	database.DB = testDB

	// Set JWT secret
	middleware.SetJWTSecret("test-secret-key")

	// Seed test data
	seedTestData()

	// Setup router
	testRouter = gin.New()
	router.SetupRoutes(testRouter)

	os.Exit(m.Run())
}

func seedTestData() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 12)

	// Department
	testDeptID = uuid.New().String()
	testDB.Create(&model.Department{
		ID:   testDeptID,
		Name: "IT Department",
	})

	// Location
	testLocationID = uuid.New().String()
	testDB.Create(&model.Location{
		ID:       testLocationID,
		Name:     "Room A-301",
		Building: "A Building",
		Floor:    "3",
		Room:     "301",
	})

	// Admin
	adminID := uuid.New().String()
	testDB.Create(&model.User{
		ID: adminID, FirstName: "Test", LastName: "Admin",
		Email: "admin@test.com", Username: "admin",
		PasswordHash: string(hash), Status: model.UserActive,
		Role: model.Admin,
	})
	adminToken, _ = middleware.GenerateToken(adminID, "admin", string(model.Admin), 24)

	// Custodian
	custodianID := uuid.New().String()
	testDB.Create(&model.User{
		ID: custodianID, FirstName: "Test", LastName: "Custodian",
		Email: "custodian@test.com", Username: "custodian",
		PasswordHash: string(hash), Status: model.UserActive,
		Role: model.AssetCustodian,
	})
	custodianToken, _ = middleware.GenerateToken(custodianID, "custodian", string(model.AssetCustodian), 24)

	// Accountant
	accountantID := uuid.New().String()
	testDB.Create(&model.User{
		ID: accountantID, FirstName: "Test", LastName: "Accountant",
		Email: "accountant@test.com", Username: "accountant",
		PasswordHash: string(hash), Status: model.UserActive,
		Role: model.Accountant,
	})
	accountantToken, _ = middleware.GenerateToken(accountantID, "accountant", string(model.Accountant), 24)

	// Employee
	employeeID := uuid.New().String()
	testDB.Create(&model.User{
		ID: employeeID, FirstName: "Test", LastName: "Employee",
		Email: "employee@test.com", Username: "employee",
		PasswordHash: string(hash), Status: model.UserActive,
		Role: model.Employee,
	})
	employeeToken, _ = middleware.GenerateToken(employeeID, "employee", string(model.Employee), 24)

	fmt.Println("Test data seeded")
}
