package controllers

import (
	"github.com/Hustle299/Project-0/views"
)

type Static struct {
	Home    *views.View
	Contact *views.View
}

// Function load cac trang tinh, khong thay doi
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "statics/homepage"),
		Contact: views.NewView("bootstrap", "statics/contact"),
	}
}
