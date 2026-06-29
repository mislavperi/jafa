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
	"github.com/mislavperi/jafa/server/internal/domain/apperr"
	"github.com/mislavperi/jafa/server/internal/domain/dto"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	"github.com/mislavperi/jafa/server/internal/server/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// stubExpenseService implements the ExpenseService interface the controller
// depends on. Only the methods a test exercises need a fn; the rest return zero
// values. Controller tests assert HTTP behaviour (status, body, param parsing,
// error mapping) — service logic is covered by the service's own tests.
var _ ExpenseService = (*stubExpenseService)(nil)

type stubExpenseService struct {
	getAllExpensesFn func(ctx context.Context, userID int64) ([]models.Expense, error)
	deleteExpenseFn  func(ctx context.Context, userID, id int64) error
	updateExpenseFn  func(ctx context.Context, input dto.UpdateExpenseInput) (models.Expense, error)
	createExpenseFn  func(ctx context.Context, input dto.CreateExpenseInput) (models.Expense, error)
	getDailySpendFn  func(ctx context.Context, userID int64, months int32) ([]models.DailySpend, error)
	getByIdFn        func(ctx context.Context, userID, id int64) (models.Expense, error)
}

func (s *stubExpenseService) CreateExpense(ctx context.Context, input dto.CreateExpenseInput) (models.Expense, error) {
	return s.createExpenseFn(ctx, input)
}
func (s *stubExpenseService) BulkCreateExpenses(_ context.Context, _ int64, _ []dto.BulkExpenseItem) ([]models.Expense, error) {
	return nil, nil
}
func (s *stubExpenseService) GetAllExpenses(ctx context.Context, userID int64) ([]models.Expense, error) {
	return s.getAllExpensesFn(ctx, userID)
}
func (s *stubExpenseService) GetAllEntries(_ context.Context, _ int64) ([]models.Expense, error) {
	return nil, nil
}
func (s *stubExpenseService) GetTotalSpendThisMonth(_ context.Context, _ int64) (models.MonthlyTotal, error) {
	return models.MonthlyTotal{}, nil
}
func (s *stubExpenseService) GetTotalIncomeThisMonth(_ context.Context, _ int64) (models.MonthlyTotal, error) {
	return models.MonthlyTotal{}, nil
}
func (s *stubExpenseService) GetDailySpend(ctx context.Context, userID int64, months int32) ([]models.DailySpend, error) {
	return s.getDailySpendFn(ctx, userID, months)
}
func (s *stubExpenseService) GetFirstExpenseDate(_ context.Context, _ int64) (string, error) {
	return "", nil
}
func (s *stubExpenseService) GetDailySpendForMonth(_ context.Context, _ int64, _, _ int32) ([]models.DailySpend, error) {
	return nil, nil
}
func (s *stubExpenseService) GetExpensesByMonth(_ context.Context, _ int64, _, _ int32) ([]models.Expense, error) {
	return nil, nil
}
func (s *stubExpenseService) UpdateExpense(ctx context.Context, input dto.UpdateExpenseInput) (models.Expense, error) {
	return s.updateExpenseFn(ctx, input)
}
func (s *stubExpenseService) DeleteExpense(ctx context.Context, userID, id int64) error {
	return s.deleteExpenseFn(ctx, userID, id)
}
func (s *stubExpenseService) GetById(ctx context.Context, userID, id int64) (models.Expense, error) {
	return s.getByIdFn(ctx, userID, id)
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
	stub := &stubExpenseService{
		getAllExpensesFn: func(_ context.Context, userID int64) ([]models.Expense, error) {
			if userID != 42 {
				t.Errorf("userID = %d, want 42", userID)
			}
			return []models.Expense{{Id: 1, Name: "Coffee"}}, nil
		},
	}
	ec := NewExpenseController(stub)
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
	stub := &stubExpenseService{
		getAllExpensesFn: func(_ context.Context, _ int64) ([]models.Expense, error) {
			return nil, errors.New("db down")
		},
	}
	ec := NewExpenseController(stub)
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
	stub := &stubExpenseService{
		deleteExpenseFn: func(_ context.Context, userID, id int64) error {
			if userID != 7 || id != 5 {
				t.Errorf("Delete(user=%d,id=%d), want (7,5)", userID, id)
			}
			return nil
		},
	}
	ec := NewExpenseController(stub)
	r := newRouter(7, http.MethodDelete, "/expense/:id", ec.DeleteExpense())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/expense/5", nil))

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", w.Code)
	}
}

func TestDeleteExpense_NotFound(t *testing.T) {
	stub := &stubExpenseService{
		deleteExpenseFn: func(_ context.Context, _, _ int64) error { return apperr.ErrExpenseNotFound },
	}
	ec := NewExpenseController(stub)
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
	stub := &stubExpenseService{
		updateExpenseFn: func(_ context.Context, _ dto.UpdateExpenseInput) (models.Expense, error) {
			return models.Expense{}, apperr.ErrExpenseNotFound
		},
	}
	ec := NewExpenseController(stub)
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
	stub := &stubExpenseService{
		updateExpenseFn: func(_ context.Context, input dto.UpdateExpenseInput) (models.Expense, error) {
			if input.ID != 10 || input.UserID != 3 {
				t.Errorf("wrong input: %+v", input)
			}
			return models.Expense{Id: 10, Name: "Dinner"}, nil
		},
	}
	ec := NewExpenseController(stub)
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
	stub := &stubExpenseService{
		createExpenseFn: func(_ context.Context, input dto.CreateExpenseInput) (models.Expense, error) {
			if input.Name != "Groceries" {
				t.Errorf("Name = %q, want Groceries", input.Name)
			}
			return models.Expense{Id: 99, Name: "Groceries"}, nil
		},
	}
	ec := NewExpenseController(stub)
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
	stub := &stubExpenseService{
		createExpenseFn: func(_ context.Context, input dto.CreateExpenseInput) (models.Expense, error) {
			if input.InstallmentCount == nil || *input.InstallmentCount != 4 {
				t.Errorf("InstallmentCount = %v, want 4", input.InstallmentCount)
			}
			return models.Expense{Id: 1, Name: "Phone", InstallmentPlan: &models.InstallmentPlan{Count: 4, PaymentAmount: 50}}, nil
		},
	}
	ec := NewExpenseController(stub)
	r := newRouter(1, http.MethodPost, "/expense/", ec.CreateExpense())

	body := jsonBody(t, map[string]any{"name": "Phone", "amount": 1.0, "cost": 200.0, "installmentCount": 4})
	req := httptest.NewRequest(http.MethodPost, "/expense/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", w.Code)
	}
	var got models.Expense
	json.Unmarshal(w.Body.Bytes(), &got)
	if got.InstallmentPlan == nil || got.InstallmentPlan.Count != 4 || got.InstallmentPlan.PaymentAmount != 50 {
		t.Errorf("InstallmentPlan = %+v, want {4 50}", got.InstallmentPlan)
	}
}

// installmentCount<2 is rejected by request binding before the service runs.
func TestCreateExpense_InstallmentCountTooLow(t *testing.T) {
	stub := &stubExpenseService{
		createExpenseFn: func(_ context.Context, _ dto.CreateExpenseInput) (models.Expense, error) {
			t.Fatal("service should not be reached; binding rejects count<2")
			return models.Expense{}, nil
		},
	}
	ec := NewExpenseController(stub)
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

// The controller maps a service-side validation error (ErrInvalidStartDate) to 400.
func TestCreateExpense_InvalidStartDate(t *testing.T) {
	stub := &stubExpenseService{
		createExpenseFn: func(_ context.Context, _ dto.CreateExpenseInput) (models.Expense, error) {
			return models.Expense{}, apperr.ErrInvalidStartDate
		},
	}
	ec := NewExpenseController(stub)
	r := newRouter(1, http.MethodPost, "/expense/", ec.CreateExpense())

	body := jsonBody(t, map[string]any{
		"name": "Rent", "amount": 1.0, "cost": 100.0,
		"recurringSchedule": map[string]any{
			"interval": "monthly", "dayOfMonth": 1, "startDate": "not-a-date",
		},
	})
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
	stub := &stubExpenseService{
		getDailySpendFn: func(_ context.Context, _ int64, months int32) ([]models.DailySpend, error) {
			if months != 3 {
				t.Errorf("months = %d, want 3", months)
			}
			return []models.DailySpend{{Day: "2024-01-01", Total: 30}}, nil
		},
	}
	ec := NewExpenseController(stub)
	r := newRouter(1, http.MethodGet, "/expense-stats/daily-spend", ec.GetDailySpend())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense-stats/daily-spend?months=3", nil))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

// ---- GetExpenseById ----

func TestGetExpenseById_NotFound(t *testing.T) {
	stub := &stubExpenseService{
		getByIdFn: func(_ context.Context, _, _ int64) (models.Expense, error) {
			return models.Expense{}, apperr.ErrExpenseNotFound
		},
	}
	ec := NewExpenseController(stub)
	r := newRouter(1, http.MethodGet, "/expense/:id", ec.GetExpenseById())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense/99", nil))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestGetExpenseById_OK(t *testing.T) {
	stub := &stubExpenseService{
		getByIdFn: func(_ context.Context, _, id int64) (models.Expense, error) {
			return models.Expense{Id: id, Name: "Lunch"}, nil
		},
	}
	ec := NewExpenseController(stub)
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
