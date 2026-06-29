package controllers

import (
	"context"

	"github.com/mislavperi/jafa/server/internal/domain/dto"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	requestmodels "github.com/mislavperi/jafa/server/internal/domain/models/request"
	"github.com/mislavperi/jafa/server/internal/domain/services"
)

// This file declares the service interfaces the controllers depend on. They are
// defined in the consumer (controller) package so each controller depends on an
// abstraction rather than a concrete *services.XService. The concrete services
// satisfy these, so bootstrap can wire them in unchanged.

// The concrete services satisfy the interfaces the controllers depend on.
var (
	_ ExpenseService     = (*services.ExpenseService)(nil)
	_ TagService         = (*services.TagService)(nil)
	_ AuthService        = (*services.AuthService)(nil)
	_ CategoryService    = (*services.CategoryService)(nil)
	_ PreferencesService = (*services.PreferencesService)(nil)
	_ ReportService      = (*services.ReportService)(nil)
)

// ExpenseService is the behaviour ExpenseController needs.
type ExpenseService interface {
	CreateExpense(ctx context.Context, input dto.CreateExpenseInput) (models.Expense, error)
	BulkCreateExpenses(ctx context.Context, userID int64, items []dto.BulkExpenseItem) ([]models.Expense, error)
	GetAllExpenses(ctx context.Context, userID int64) ([]models.Expense, error)
	GetAllEntries(ctx context.Context, userID int64) ([]models.Expense, error)
	GetTotalSpendThisMonth(ctx context.Context, userID int64) (models.MonthlyTotal, error)
	GetTotalIncomeThisMonth(ctx context.Context, userID int64) (models.MonthlyTotal, error)
	GetDailySpend(ctx context.Context, userID int64, months int32) ([]models.DailySpend, error)
	GetFirstExpenseDate(ctx context.Context, userID int64) (string, error)
	GetDailySpendForMonth(ctx context.Context, userID int64, year, month int32) ([]models.DailySpend, error)
	GetExpensesByMonth(ctx context.Context, userID int64, year, month int32) ([]models.Expense, error)
	UpdateExpense(ctx context.Context, input dto.UpdateExpenseInput) (models.Expense, error)
	DeleteExpense(ctx context.Context, userID, id int64) error
	GetById(ctx context.Context, userID, id int64) (models.Expense, error)
}

// TagService is the behaviour TagController needs.
type TagService interface {
	GetAllTags(ctx context.Context, userID int64) ([]models.Tag, error)
	CreateTag(ctx context.Context, userID int64, name string) (models.Tag, error)
	GetTagsForExpense(ctx context.Context, userID, expenseID int64) ([]models.Tag, error)
	AddTagToExpense(ctx context.Context, userID, expenseID, tagID int64) error
	RemoveTagFromExpense(ctx context.Context, userID, expenseID, tagID int64) error
}

// AuthService is the behaviour AuthController needs.
type AuthService interface {
	Login(ctx context.Context, username, password string) (models.User, error)
	Register(ctx context.Context, params requestmodels.RegisterRequest) (models.User, error)
	DeleteAccount(ctx context.Context, userID int64) error
}

// CategoryService is the behaviour CategoryController needs.
type CategoryService interface {
	List(ctx context.Context) ([]models.Category, error)
}

// PreferencesService is the behaviour PreferencesController needs.
type PreferencesService interface {
	Get(ctx context.Context, userID int64) (models.UserPreferences, error)
	Upsert(ctx context.Context, input dto.UpsertPreferencesInput) (models.UserPreferences, error)
}

// ReportService is the behaviour ReportController needs.
type ReportService interface {
	CategoryBreakdown(ctx context.Context, userID int64) ([]models.CategoryBreakdown, error)
	MonthlySpend(ctx context.Context, userID int64) ([]models.MonthlySpend, error)
}
