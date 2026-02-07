package models

import "sort"

type Summary struct {
	Income            float64
	TotalExpenses     float64
	Remaining         float64
	SavingsRate       float64
	DailyAllowance    float64
	CategoryBreakdown []CategoryTotal
}

type CategoryTotal struct {
	Category Category
	Total    float64
}

func CalculateSummary(income float64, expenses []Expense, daysInMonth int) Summary {
	var totalExpenses float64
	categoryTotals := make(map[int64]float64)
	categoryMap := make(map[int64]Category)

	for _, e := range expenses {
		totalExpenses += e.Amount
		if e.CategoryID != nil {
			categoryTotals[*e.CategoryID] += e.Amount
			if e.Category != nil {
				categoryMap[*e.CategoryID] = *e.Category
			}
		}
	}

	remaining := income - totalExpenses
	savingsRate := 0.0
	if income > 0 {
		savingsRate = (remaining / income) * 100
	}
	dailyAllowance := 0.0
	if daysInMonth > 0 {
		dailyAllowance = remaining / float64(daysInMonth)
	}
	if dailyAllowance < 0 {
		dailyAllowance = 0
	}

	breakdown := make([]CategoryTotal, 0, len(categoryTotals))
	for catID, total := range categoryTotals {
		breakdown = append(breakdown, CategoryTotal{
			Category: categoryMap[catID],
			Total:    total,
		})
	}
	sort.Slice(breakdown, func(i, j int) bool {
		return breakdown[i].Total > breakdown[j].Total
	})

	return Summary{
		Income:            income,
		TotalExpenses:     totalExpenses,
		Remaining:         remaining,
		SavingsRate:       savingsRate,
		DailyAllowance:    dailyAllowance,
		CategoryBreakdown: breakdown,
	}
}
