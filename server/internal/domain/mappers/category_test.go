package mappers

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	psql "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
)

func TestCategoryMapper_MapToDomain_Basic(t *testing.T) {
	mapper := NewCategoryMapper()
	var budget pgtype.Numeric
	if err := budget.Scan("500.000"); err != nil {
		t.Fatalf("Scan: %v", err)
	}

	cat := psql.Category{
		ID:       1,
		Name:     "Groceries",
		Icon:     "🛒",
		Color:    "#00ff00",
		Budget:   budget,
		Keywords: []string{"market", "food"},
	}

	got, err := mapper.MapToDomain(cat)
	if err != nil {
		t.Fatalf("MapToDomain: %v", err)
	}
	if got.ID != 1 {
		t.Errorf("ID = %d, want 1", got.ID)
	}
	if got.Name != "Groceries" {
		t.Errorf("Name = %q, want Groceries", got.Name)
	}
	if got.Budget < 499.9 || got.Budget > 500.1 {
		t.Errorf("Budget = %v, want ~500", got.Budget)
	}
	if len(got.Keywords) != 2 {
		t.Errorf("Keywords len = %d, want 2", len(got.Keywords))
	}
}

func TestCategoryMapper_MapToDomain_NilKeywords(t *testing.T) {
	mapper := NewCategoryMapper()
	var budget pgtype.Numeric
	budget.Scan("0")

	cat := psql.Category{
		ID:       2,
		Name:     "Other",
		Keywords: nil,
		Budget:   budget,
	}

	got, err := mapper.MapToDomain(cat)
	if err != nil {
		t.Fatalf("MapToDomain: %v", err)
	}
	if got.Keywords == nil {
		t.Error("Keywords is nil, want empty slice")
	}
	if len(got.Keywords) != 0 {
		t.Errorf("Keywords len = %d, want 0", len(got.Keywords))
	}
}

func TestCategoryMapper_MapManyToDomain(t *testing.T) {
	mapper := NewCategoryMapper()
	var budget pgtype.Numeric
	budget.Scan("100.000")

	cats := []psql.Category{
		{ID: 1, Name: "Dining", Budget: budget, Keywords: []string{"restaurant"}},
		{ID: 2, Name: "Transport", Budget: budget, Keywords: []string{"uber", "taxi"}},
	}

	got, err := mapper.MapManyToDomain(cats)
	if err != nil {
		t.Fatalf("MapManyToDomain: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Name != "Dining" || got[1].Name != "Transport" {
		t.Errorf("names = [%q, %q]", got[0].Name, got[1].Name)
	}
}

func TestCategoryMapper_MapManyToDomain_Empty(t *testing.T) {
	mapper := NewCategoryMapper()
	got, err := mapper.MapManyToDomain(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}
