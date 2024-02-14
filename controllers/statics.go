package controllers

import (
	"github.com/Hustle299/Project-0/views"
)

type Static struct {
	Home    *views.View
	Contact *views.View
}

// Function to serve static resource like homepage and contact page through NewView function in view folder
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "statics/homepage"),
		Contact: views.NewView("bootstrap", "statics/contact"),
	}
}
