package models

import "time"

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

func (c Category) BgClass() string {
	return c.Color + "-100"
}

func (c Category) TextClass() string {
	return c.Color + "-700"
}

func (c Category) DotClass() string {
	return c.Color + "-500"
}
