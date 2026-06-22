package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/server/internal/server/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// stubExpenseService implements services.ExpenseServicer for tests.
type stubExpenseService struct {
	getAllExpensesFn     func(userID int64) ([]models.Expense, error)
	getAllEntriesFn      func(userID int64) ([]models.Expense, error)
	getByIdFn            func(userID, id int64) (models.Expense, error)
	getTotalFn           func(userID int64) (models.MonthlyTotal, error)
	getIncomeFn          func(userID int64) (models.MonthlyTotal, error)
	getDailySpendFn      func(userID int64, months int32) ([]models.DailySpend, error)
	getFirstDateFn       func(userID int64) (string, error)
	getDailyForMonthFn   func(userID int64, year, month int32) ([]models.DailySpend, error)
	getExpensesByMonthFn func(userID int64, year, month int32) ([]models.Expense, error)
	createExpenseFn      func(input services.CreateExpenseInput) (models.Expense, error)
	bulkCreateExpensesFn func(userID int64, items []services.BulkExpenseItem) ([]models.Expense, error)
	updateExpenseFn      func(input services.UpdateExpenseInput) (models.Expense, error)
	deleteExpenseFn      func(userID, id int64) error
}

func (s *stubExpenseService) GetAllExpenses(_ context.Context, userID int64) ([]models.Expense, error) {
	return s.getAllExpensesFn(userID)
}
func (s *stubExpenseService) GetAllEntries(_ context.Context, userID int64) ([]models.Expense, error) {
	if s.getAllEntriesFn != nil {
		return s.getAllEntriesFn(userID)
	}
	return nil, nil
}
func (s *stubExpenseService) GetById(_ context.Context, userID, id int64) (models.Expense, error) {
	return s.getByIdFn(userID, id)
}
func (s *stubExpenseService) GetTotalSpendThisMonth(_ context.Context, userID int64) (models.MonthlyTotal, error) {
	return s.getTotalFn(userID)
}
func (s *stubExpenseService) GetTotalIncomeThisMonth(_ context.Context, userID int64) (models.MonthlyTotal, error) {
	if s.getIncomeFn != nil {
		return s.getIncomeFn(userID)
	}
	return models.MonthlyTotal{}, nil
}
func (s *stubExpenseService) GetDailySpend(_ context.Context, userID int64, months int32) ([]models.DailySpend, error) {
	return s.getDailySpendFn(userID, months)
}
func (s *stubExpenseService) GetFirstExpenseDate(_ context.Context, userID int64) (string, error) {
	return s.getFirstDateFn(userID)
}
func (s *stubExpenseService) GetDailySpendForMonth(_ context.Context, userID int64, year, month int32) ([]models.DailySpend, error) {
	return s.getDailyForMonthFn(userID, year, month)
}
func (s *stubExpenseService) GetExpensesByMonth(_ context.Context, userID int64, year, month int32) ([]models.Expense, error) {
	return s.getExpensesByMonthFn(userID, year, month)
}
func (s *stubExpenseService) CreateExpense(_ context.Context, input services.CreateExpenseInput) (models.Expense, error) {
	return s.createExpenseFn(input)
}
func (s *stubExpenseService) BulkCreateExpenses(_ context.Context, userID int64, items []services.BulkExpenseItem) ([]models.Expense, error) {
	return s.bulkCreateExpensesFn(userID, items)
}
func (s *stubExpenseService) UpdateExpense(_ context.Context, input services.UpdateExpenseInput) (models.Expense, error) {
	return s.updateExpenseFn(input)
}
func (s *stubExpenseService) DeleteExpense(_ context.Context, userID, id int64) error {
	return s.deleteExpenseFn(userID, id)
}

// newRouter builds a test router that injects a fixed user ID into the context.
func newRouter(userID int64, method, path string, handler gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(func(ctx *gin.Context) {
		ctx.Set(middleware.ContextUserIDKey, userID)
		ctx.Next()
	})
	r.Handle(method, path, handler)
	return r
}

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

// ---- GetAllExpenses ----

func TestGetAllExpenses_OK(t *testing.T) {
	expenses := []models.Expense{{Id: 1, Name: "Coffee", Amount: 3.5}}
	svc := &stubExpenseService{
		getAllExpensesFn: func(userID int64) ([]models.Expense, error) {
			if userID != 42 {
				t.Errorf("userID = %d, want 42", userID)
			}
			return expenses, nil
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(42, http.MethodGet, "/expense/", ec.GetAllExpenses())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense/", nil))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	var got []models.Expense
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Coffee" {
		t.Errorf("body = %+v", got)
	}
}

func TestGetAllExpenses_ServiceError(t *testing.T) {
	svc := &stubExpenseService{
		getAllExpensesFn: func(_ int64) ([]models.Expense, error) {
			return nil, errors.New("db down")
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodGet, "/expense/", ec.GetAllExpenses())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense/", nil))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

func TestGetAllExpenses_Unauthenticated(t *testing.T) {
	ec := NewExpenseController(&stubExpenseService{})
	r := gin.New() // no user injected
	r.GET("/expense/", ec.GetAllExpenses())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense/", nil))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

// ---- DeleteExpense ----

func TestDeleteExpense_OK(t *testing.T) {
	svc := &stubExpenseService{
		deleteExpenseFn: func(userID, id int64) error {
			if userID != 7 || id != 5 {
				t.Errorf("DeleteExpense(%d,%d), want (7,5)", userID, id)
			}
			return nil
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(7, http.MethodDelete, "/expense/:id", ec.DeleteExpense())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/expense/5", nil))

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", w.Code)
	}
}

func TestDeleteExpense_NotFound(t *testing.T) {
	svc := &stubExpenseService{
		deleteExpenseFn: func(_, _ int64) error { return services.ErrExpenseNotFound },
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodDelete, "/expense/:id", ec.DeleteExpense())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/expense/99", nil))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestDeleteExpense_InvalidID(t *testing.T) {
	ec := NewExpenseController(&stubExpenseService{})
	r := newRouter(1, http.MethodDelete, "/expense/:id", ec.DeleteExpense())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/expense/abc", nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

// ---- UpdateExpense ----

func TestUpdateExpense_NotFound(t *testing.T) {
	svc := &stubExpenseService{
		updateExpenseFn: func(_ services.UpdateExpenseInput) (models.Expense, error) {
			return models.Expense{}, services.ErrExpenseNotFound
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodPatch, "/expense/:id", ec.UpdateExpense())

	body := jsonBody(t, map[string]any{"name": "X", "amount": 1.0, "cost": 1.0})
	req := httptest.NewRequest(http.MethodPatch, "/expense/10", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestUpdateExpense_BadBody(t *testing.T) {
	ec := NewExpenseController(&stubExpenseService{})
	r := newRouter(1, http.MethodPatch, "/expense/:id", ec.UpdateExpense())

	req := httptest.NewRequest(http.MethodPatch, "/expense/10", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestUpdateExpense_OK(t *testing.T) {
	updated := models.Expense{Id: 10, Name: "Dinner", Amount: 25}
	svc := &stubExpenseService{
		updateExpenseFn: func(inp services.UpdateExpenseInput) (models.Expense, error) {
			if inp.ID != 10 || inp.UserID != 3 {
				t.Errorf("wrong input: %+v", inp)
			}
			return updated, nil
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(3, http.MethodPatch, "/expense/:id", ec.UpdateExpense())

	body := jsonBody(t, map[string]any{"name": "Dinner", "amount": 25.0, "cost": 25.0})
	req := httptest.NewRequest(http.MethodPatch, "/expense/10", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	var got models.Expense
	json.Unmarshal(w.Body.Bytes(), &got)
	if got.Name != "Dinner" {
		t.Errorf("Name = %q, want Dinner", got.Name)
	}
}

// ---- CreateExpense ----

func TestCreateExpense_OK(t *testing.T) {
	created := models.Expense{Id: 99, Name: "Groceries", Amount: 50}
	svc := &stubExpenseService{
		createExpenseFn: func(inp services.CreateExpenseInput) (models.Expense, error) {
			if inp.Name != "Groceries" {
				t.Errorf("Name = %q, want Groceries", inp.Name)
			}
			return created, nil
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodPost, "/expense/", ec.CreateExpense())

	body := jsonBody(t, map[string]any{"name": "Groceries", "amount": 50.0, "cost": 50.0})
	req := httptest.NewRequest(http.MethodPost, "/expense/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", w.Code)
	}
}

func TestCreateExpense_WithInstallments(t *testing.T) {
	var capturedCount *int
	created := models.Expense{Id: 1, Name: "Phone", Cost: 200, InstallmentPlan: &models.InstallmentPlan{Count: 4, PaymentAmount: 50}}
	svc := &stubExpenseService{
		createExpenseFn: func(inp services.CreateExpenseInput) (models.Expense, error) {
			capturedCount = inp.InstallmentCount
			return created, nil
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodPost, "/expense/", ec.CreateExpense())

	body := jsonBody(t, map[string]any{"name": "Phone", "amount": 1.0, "cost": 200.0, "installmentCount": 4})
	req := httptest.NewRequest(http.MethodPost, "/expense/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", w.Code)
	}
	if capturedCount == nil || *capturedCount != 4 {
		t.Errorf("InstallmentCount = %v, want 4", capturedCount)
	}
	var got models.Expense
	json.Unmarshal(w.Body.Bytes(), &got)
	if got.InstallmentPlan == nil || got.InstallmentPlan.Count != 4 || got.InstallmentPlan.PaymentAmount != 50 {
		t.Errorf("InstallmentPlan = %+v, want {4 50}", got.InstallmentPlan)
	}
}

func TestCreateExpense_InstallmentCountTooLow(t *testing.T) {
	svc := &stubExpenseService{
		createExpenseFn: func(_ services.CreateExpenseInput) (models.Expense, error) {
			t.Fatal("service should not be reached; binding rejects count<2")
			return models.Expense{}, nil
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodPost, "/expense/", ec.CreateExpense())

	body := jsonBody(t, map[string]any{"name": "Phone", "amount": 1.0, "cost": 200.0, "installmentCount": 1})
	req := httptest.NewRequest(http.MethodPost, "/expense/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestCreateExpense_InvalidInstallmentFromService(t *testing.T) {
	svc := &stubExpenseService{
		createExpenseFn: func(_ services.CreateExpenseInput) (models.Expense, error) {
			return models.Expense{}, services.ErrInvalidInstallmentCount
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodPost, "/expense/", ec.CreateExpense())

	body := jsonBody(t, map[string]any{"name": "Phone", "amount": 1.0, "cost": 200.0, "installmentCount": 4})
	req := httptest.NewRequest(http.MethodPost, "/expense/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestCreateExpense_MissingFields(t *testing.T) {
	ec := NewExpenseController(&stubExpenseService{})
	r := newRouter(1, http.MethodPost, "/expense/", ec.CreateExpense())

	body := jsonBody(t, map[string]any{"name": "X"}) // missing amount and cost
	req := httptest.NewRequest(http.MethodPost, "/expense/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

// ---- GetDailySpend ----

func TestGetDailySpend_InvalidMonths(t *testing.T) {
	ec := NewExpenseController(&stubExpenseService{})
	r := newRouter(1, http.MethodGet, "/expense-stats/daily-spend", ec.GetDailySpend())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense-stats/daily-spend?months=abc", nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestGetDailySpend_OK(t *testing.T) {
	rows := []models.DailySpend{{Day: "2024-06-01", Total: 30}}
	svc := &stubExpenseService{
		getDailySpendFn: func(_ int64, months int32) ([]models.DailySpend, error) {
			if months != 3 {
				t.Errorf("months = %d, want 3", months)
			}
			return rows, nil
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodGet, "/expense-stats/daily-spend", ec.GetDailySpend())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense-stats/daily-spend?months=3", nil))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

// ---- GetExpenseById ----

func TestGetExpenseById_NotFound(t *testing.T) {
	svc := &stubExpenseService{
		getByIdFn: func(_, _ int64) (models.Expense, error) {
			return models.Expense{}, services.ErrExpenseNotFound
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodGet, "/expense/:id", ec.GetExpenseById())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense/99", nil))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestGetExpenseById_OK(t *testing.T) {
	svc := &stubExpenseService{
		getByIdFn: func(_, id int64) (models.Expense, error) {
			return models.Expense{Id: id, Name: "Lunch"}, nil
		},
	}
	ec := NewExpenseController(svc)
	r := newRouter(1, http.MethodGet, "/expense/:id", ec.GetExpenseById())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense/7", nil))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	var got models.Expense
	json.Unmarshal(w.Body.Bytes(), &got)
	if got.Name != "Lunch" {
		t.Errorf("Name = %q, want Lunch", got.Name)
	}
}
