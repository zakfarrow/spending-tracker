package models

import "time"

type Period struct {
	Month int
	Year  int
}

func (p Period) MonthName() string {
	return time.Month(p.Month).String()
}

func (p Period) Prev() Period {
	if p.Month == 1 {
		return Period{Month: 12, Year: p.Year - 1}
	}
	return Period{Month: p.Month - 1, Year: p.Year}
}

func (p Period) Next() Period {
	if p.Month == 12 {
		return Period{Month: 1, Year: p.Year + 1}
	}
	return Period{Month: p.Month + 1, Year: p.Year}
}

func (p Period) DaysInMonth() int {
	return time.Date(p.Year, time.Month(p.Month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func CurrentPeriod() Period {
	now := time.Now()
	return Period{Month: int(now.Month()), Year: now.Year()}
}
