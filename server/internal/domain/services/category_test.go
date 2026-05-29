package services

import "testing"

import "github.com/mislavperi/jafa/server/internal/domain/models"

func TestCategorize(t *testing.T) {
	categories := []models.Category{
		{Name: "Groceries", Keywords: []string{"market", "food"}},
		{Name: "Dining", Keywords: []string{"restaurant", "coffee"}},
		{Name: "Other", Keywords: []string{}},
	}

	cases := map[string]string{
		"Whole Foods Market": "Groceries",
		"Corner Coffee Shop": "Dining",
		"Mystery Charge":     "Other",
		"FOOD truck":         "Groceries", // case-insensitive
	}

	for name, want := range cases {
		if got := Categorize(name, categories); got != want {
			t.Errorf("Categorize(%q) = %q, want %q", name, got, want)
		}
	}
}
