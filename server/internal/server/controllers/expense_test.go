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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mislavperi/jafa/server/internal/domain/mappers"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	"github.com/mislavperi/jafa/server/internal/domain/services"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	"github.com/mislavperi/jafa/server/internal/server/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// Ensure stubExpenseQuerier satisfies the interface at compile time.
var _ services.ExpenseQuerier = (*stubExpenseQuerier)(nil)

// stubExpenseQuerier implements services.ExpenseQuerier so handlers run against a
// real services.ExpenseService. Only the methods a test exercises need a stub fn;
// the rest return zero values.
type stubExpenseQuerier struct {
	getAllExpensesFn func(ctx context.Context, userID int64) ([]psql.Expense, error)
	softDeleteFn     func(ctx context.Context, arg psql.SoftDeleteExpenseParams) (int64, error)
	updateExpenseFn  func(ctx context.Context, arg psql.UpdateExpenseParams) (psql.Expense, error)
	createExpenseFn  func(ctx context.Context, arg psql.CreateExpenseParams) (psql.Expense, error)
	getDailySpendFn  func(ctx context.Context, arg psql.GetDailySpendParams) ([]psql.GetDailySpendRow, error)
	getExpenseByIdFn func(ctx context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error)
}

func (s *stubExpenseQuerier) GetAllExpenses(ctx context.Context, userID int64) ([]psql.Expense, error) {
	return s.getAllExpensesFn(ctx, userID)
}
func (s *stubExpenseQuerier) GetAllEntries(_ context.Context, _ int64) ([]psql.Expense, error) {
	return nil, nil
}
func (s *stubExpenseQuerier) GetExpenseById(ctx context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error) {
	return s.getExpenseByIdFn(ctx, arg)
}
func (s *stubExpenseQuerier) CreateExpense(ctx context.Context, arg psql.CreateExpenseParams) (psql.Expense, error) {
	return s.createExpenseFn(ctx, arg)
}
func (s *stubExpenseQuerier) UpdateExpense(ctx context.Context, arg psql.UpdateExpenseParams) (psql.Expense, error) {
	return s.updateExpenseFn(ctx, arg)
}
func (s *stubExpenseQuerier) SoftDeleteExpense(ctx context.Context, arg psql.SoftDeleteExpenseParams) (int64, error) {
	return s.softDeleteFn(ctx, arg)
}
func (s *stubExpenseQuerier) GetTotalSpendThisMonth(_ context.Context, _ int64) (pgtype.Numeric, error) {
	return pgtype.Numeric{}, nil
}
func (s *stubExpenseQuerier) GetTotalIncomeThisMonth(_ context.Context, _ int64) (pgtype.Numeric, error) {
	return pgtype.Numeric{}, nil
}
func (s *stubExpenseQuerier) GetDailySpend(ctx context.Context, arg psql.GetDailySpendParams) ([]psql.GetDailySpendRow, error) {
	return s.getDailySpendFn(ctx, arg)
}
func (s *stubExpenseQuerier) GetExpensesByMonth(_ context.Context, _ psql.GetExpensesByMonthParams) ([]psql.Expense, error) {
	return nil, nil
}
func (s *stubExpenseQuerier) GetFirstExpenseDate(_ context.Context, _ int64) (interface{}, error) {
	return nil, nil
}
func (s *stubExpenseQuerier) GetDailySpendForMonth(_ context.Context, _ psql.GetDailySpendForMonthParams) ([]psql.GetDailySpendForMonthRow, error) {
	return nil, nil
}
func (s *stubExpenseQuerier) UpsertTag(_ context.Context, _ psql.UpsertTagParams) (psql.Tag, error) {
	return psql.Tag{}, nil
}
func (s *stubExpenseQuerier) AddTagToExpense(_ context.Context, _ psql.AddTagToExpenseParams) error {
	return nil
}
func (s *stubExpenseQuerier) WithTx(_ pgx.Tx) services.ExpenseQuerier { return s }

// newExpenseController builds a controller backed by a real ExpenseService over
// the given stub querier.
func newExpenseController(q services.ExpenseQuerier) *ExpenseController {
	return NewExpenseController(&services.ExpenseService{
		Queries: q,
		Mapper:  mappers.NewExpenseMapper(),
	})
}

// numeric builds a pgtype.Numeric from a decimal string, like the DB returns.
func numeric(t *testing.T, s string) pgtype.Numeric {
	t.Helper()
	var n pgtype.Numeric
	if err := n.Scan(s); err != nil {
		t.Fatalf("scan numeric %q: %v", s, err)
	}
	return n
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
	stub := &stubExpenseQuerier{
		getAllExpensesFn: func(_ context.Context, userID int64) ([]psql.Expense, error) {
			if userID != 42 {
				t.Errorf("userID = %d, want 42", userID)
			}
			return []psql.Expense{{ID: 1, Name: "Coffee", Amount: numeric(t, "3.50"), Cost: numeric(t, "3.50")}}, nil
		},
	}
	ec := newExpenseController(stub)
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
	stub := &stubExpenseQuerier{
		getAllExpensesFn: func(_ context.Context, _ int64) ([]psql.Expense, error) {
			return nil, errors.New("db down")
		},
	}
	ec := newExpenseController(stub)
	r := newRouter(1, http.MethodGet, "/expense/", ec.GetAllExpenses())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense/", nil))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

func TestGetAllExpenses_Unauthenticated(t *testing.T) {
	ec := newExpenseController(&stubExpenseQuerier{})
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
	stub := &stubExpenseQuerier{
		softDeleteFn: func(_ context.Context, arg psql.SoftDeleteExpenseParams) (int64, error) {
			if arg.UserID != 7 || arg.ID != 5 {
				t.Errorf("SoftDelete(user=%d,id=%d), want (7,5)", arg.UserID, arg.ID)
			}
			return 1, nil
		},
	}
	ec := newExpenseController(stub)
	r := newRouter(7, http.MethodDelete, "/expense/:id", ec.DeleteExpense())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/expense/5", nil))

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", w.Code)
	}
}

func TestDeleteExpense_NotFound(t *testing.T) {
	stub := &stubExpenseQuerier{
		softDeleteFn: func(_ context.Context, _ psql.SoftDeleteExpenseParams) (int64, error) { return 0, nil },
	}
	ec := newExpenseController(stub)
	r := newRouter(1, http.MethodDelete, "/expense/:id", ec.DeleteExpense())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/expense/99", nil))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestDeleteExpense_InvalidID(t *testing.T) {
	ec := newExpenseController(&stubExpenseQuerier{})
	r := newRouter(1, http.MethodDelete, "/expense/:id", ec.DeleteExpense())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/expense/abc", nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

// ---- UpdateExpense ----

func TestUpdateExpense_NotFound(t *testing.T) {
	stub := &stubExpenseQuerier{
		updateExpenseFn: func(_ context.Context, _ psql.UpdateExpenseParams) (psql.Expense, error) {
			return psql.Expense{}, pgx.ErrNoRows
		},
	}
	ec := newExpenseController(stub)
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
	ec := newExpenseController(&stubExpenseQuerier{})
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
	stub := &stubExpenseQuerier{
		updateExpenseFn: func(_ context.Context, arg psql.UpdateExpenseParams) (psql.Expense, error) {
			if arg.ID != 10 || arg.UserID != 3 {
				t.Errorf("wrong params: %+v", arg)
			}
			return psql.Expense{ID: 10, Name: "Dinner", Amount: numeric(t, "25.00"), Cost: numeric(t, "25.00")}, nil
		},
	}
	ec := newExpenseController(stub)
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
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, arg psql.CreateExpenseParams) (psql.Expense, error) {
			if arg.Name != "Groceries" {
				t.Errorf("Name = %q, want Groceries", arg.Name)
			}
			return psql.Expense{ID: 99, Name: "Groceries", Amount: numeric(t, "50.00"), Cost: numeric(t, "50.00")}, nil
		},
	}
	ec := newExpenseController(stub)
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
	var capturedCount pgtype.Int4
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, arg psql.CreateExpenseParams) (psql.Expense, error) {
			capturedCount = arg.InstallmentCount
			return psql.Expense{
				ID:               1,
				Name:             "Phone",
				Amount:           numeric(t, "1.00"),
				Cost:             numeric(t, "200.00"),
				InstallmentCount: arg.InstallmentCount,
			}, nil
		},
	}
	ec := newExpenseController(stub)
	r := newRouter(1, http.MethodPost, "/expense/", ec.CreateExpense())

	body := jsonBody(t, map[string]any{"name": "Phone", "amount": 1.0, "cost": 200.0, "installmentCount": 4})
	req := httptest.NewRequest(http.MethodPost, "/expense/", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", w.Code)
	}
	if !capturedCount.Valid || capturedCount.Int32 != 4 {
		t.Errorf("InstallmentCount param = %+v, want {4 true}", capturedCount)
	}
	var got models.Expense
	json.Unmarshal(w.Body.Bytes(), &got)
	if got.InstallmentPlan == nil || got.InstallmentPlan.Count != 4 || got.InstallmentPlan.PaymentAmount != 50 {
		t.Errorf("InstallmentPlan = %+v, want {4 50}", got.InstallmentPlan)
	}
}

func TestCreateExpense_InstallmentCountTooLow(t *testing.T) {
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, _ psql.CreateExpenseParams) (psql.Expense, error) {
			t.Fatal("query should not be reached; binding rejects count<2")
			return psql.Expense{}, nil
		},
	}
	ec := newExpenseController(stub)
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

// TestCreateExpense_InvalidStartDate exercises the controller's mapping of a
// service-side validation error (ErrInvalidStartDate) to 400. The real service
// rejects the unparseable start date before reaching the querier.
func TestCreateExpense_InvalidStartDate(t *testing.T) {
	stub := &stubExpenseQuerier{
		createExpenseFn: func(_ context.Context, _ psql.CreateExpenseParams) (psql.Expense, error) {
			t.Fatal("query should not be reached; service rejects bad start date")
			return psql.Expense{}, nil
		},
	}
	ec := newExpenseController(stub)
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
	ec := newExpenseController(&stubExpenseQuerier{})
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
	ec := newExpenseController(&stubExpenseQuerier{})
	r := newRouter(1, http.MethodGet, "/expense-stats/daily-spend", ec.GetDailySpend())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense-stats/daily-spend?months=abc", nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestGetDailySpend_OK(t *testing.T) {
	stub := &stubExpenseQuerier{
		getDailySpendFn: func(_ context.Context, arg psql.GetDailySpendParams) ([]psql.GetDailySpendRow, error) {
			if arg.Months != 3 {
				t.Errorf("months = %d, want 3", arg.Months)
			}
			return []psql.GetDailySpendRow{{Day: pgtype.Date{Valid: true}, Total: numeric(t, "30.00")}}, nil
		},
	}
	ec := newExpenseController(stub)
	r := newRouter(1, http.MethodGet, "/expense-stats/daily-spend", ec.GetDailySpend())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense-stats/daily-spend?months=3", nil))

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

// ---- GetExpenseById ----

func TestGetExpenseById_NotFound(t *testing.T) {
	stub := &stubExpenseQuerier{
		getExpenseByIdFn: func(_ context.Context, _ psql.GetExpenseByIdParams) (psql.Expense, error) {
			return psql.Expense{}, pgx.ErrNoRows
		},
	}
	ec := newExpenseController(stub)
	r := newRouter(1, http.MethodGet, "/expense/:id", ec.GetExpenseById())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/expense/99", nil))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestGetExpenseById_OK(t *testing.T) {
	stub := &stubExpenseQuerier{
		getExpenseByIdFn: func(_ context.Context, arg psql.GetExpenseByIdParams) (psql.Expense, error) {
			return psql.Expense{ID: arg.ID, Name: "Lunch", Amount: numeric(t, "0.00"), Cost: numeric(t, "0.00")}, nil
		},
	}
	ec := newExpenseController(stub)
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
