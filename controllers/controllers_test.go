package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cesarbmathec/medical-exams-backend/config"
	"github.com/cesarbmathec/medical-exams-backend/dtos"
	"github.com/cesarbmathec/medical-exams-backend/middleware"
	"github.com/cesarbmathec/medical-exams-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	file := fmt.Sprintf("file:memdb_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Patient{},
		&models.ExamCategory{},
		&models.SampleType{},
		&models.ExamType{},
		&models.ExamParameter{},
		&models.Order{},
		&models.OrderExam{},
		&models.ExamResult{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	config.DB = db
	return db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.RecoveryMiddleware())

	api := r.Group("/api/v1")
	api.POST("/login", Login)
	api.POST("/register", Register)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/patients", CreatePatient)
	protected.GET("/patients", GetPatients)
	protected.POST("/orders", CreateOrder)
	protected.GET("/orders", GetOrders)
	protected.GET("/lab/exams/catalog", GetExamCatalog)

	return r
}

func seedAuthData(t *testing.T, db *gorm.DB) models.User {
	role := models.Role{Name: "admin", Description: "admin"}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}

	user := models.User{
		Username: "admin",
		Email:    "admin@test.com",
		Password: "Admin123!",
		FullName: "Admin",
		RoleID:   role.ID,
		IsActive: true,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func getToken(t *testing.T, r *gin.Engine, username, password string) string {
	body := dtos.LoginRequest{Username: username, Password: password}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("login failed: %d", resp.Code)
	}

	var parsed struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("unmarshal login: %v", err)
	}
	if parsed.Data.Token == "" {
		t.Fatal("empty token")
	}
	return parsed.Data.Token
}

func TestAuthLoginAndRegister(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	db := setupTestDB(t)
	seedAuthData(t, db)
	r := setupRouter()

	register := dtos.RegisterRequest{
		Username: "user1",
		Email:    "user1@test.com",
		Password: "Admin123!",
		FullName: "User One",
		RoleID:   1,
	}
	payload, _ := json.Marshal(register)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("register failed: %d", resp.Code)
	}

	token := getToken(t, r, "admin", "Admin123!")
	if token == "" {
		t.Fatal("expected token")
	}
}

func TestPatientsEndpoints(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	db := setupTestDB(t)
	seedAuthData(t, db)
	r := setupRouter()

	token := getToken(t, r, "admin", "Admin123!")

	patient := dtos.CreatePatientRequest{
		DocumentType:   "cedula",
		DocumentNumber: "V12345678",
		FirstName:      "Maria",
		LastName:       "Delgado",
		DateOfBirth:    time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		Gender:         "F",
	}
	payload, _ := json.Marshal(patient)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/patients", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create patient failed: %d", resp.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	r.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("list patients failed: %d", listResp.Code)
	}
}

func TestOrdersAndLabCatalog(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	db := setupTestDB(t)
	seedAuthData(t, db)
	r := setupRouter()

	category := models.ExamCategory{Name: "Hematologia", Code: "HEM"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	sample := models.SampleType{Name: "Sangre"}
	if err := db.Create(&sample).Error; err != nil {
		t.Fatalf("create sample: %v", err)
	}
	examType := models.ExamType{
		Code:         "HB",
		Name:         "Hemoglobina",
		CategoryID:   category.ID,
		SampleTypeID: sample.ID,
		BasePrice:    10,
	}
	if err := db.Create(&examType).Error; err != nil {
		t.Fatalf("create exam type: %v", err)
	}
	patient := models.Patient{
		DocumentType:   "cedula",
		DocumentNumber: "V98765432",
		FirstName:      "Luis",
		LastName:       "Perez",
		DateOfBirth:    time.Date(1992, 7, 10, 0, 0, 0, 0, time.UTC),
		Gender:         "M",
		CreatedBy:      1,
	}
	if err := db.Create(&patient).Error; err != nil {
		t.Fatalf("create patient: %v", err)
	}

	token := getToken(t, r, "admin", "Admin123!")

	order := dtos.CreateOrderRequest{
		PatientID: patient.ID,
		Priority:  "normal",
		Exams: []dtos.OrderExamRequest{
			{ExamTypeID: examType.ID, Price: 10},
		},
	}
	payload, _ := json.Marshal(order)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create order failed: %d", resp.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	r.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("list orders failed: %d", listResp.Code)
	}

	catReq := httptest.NewRequest(http.MethodGet, "/api/v1/lab/exams/catalog", nil)
	catReq.Header.Set("Authorization", "Bearer "+token)
	catResp := httptest.NewRecorder()
	r.ServeHTTP(catResp, catReq)
	if catResp.Code != http.StatusOK {
		t.Fatalf("lab catalog failed: %d", catResp.Code)
	}
}
