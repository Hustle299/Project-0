package controllers

import (
	"github.com/Hustle299/Project-0/views"
)

// Function to serve static resource like homepage and contact page

type Static struct {
	Home    *views.View
	Contact *views.View
}

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "statics/homepage.gohtml"),
		Contact: views.NewView("bootstrap", "statics/contact.gohtml"),
	}
}
