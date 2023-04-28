package controllers

import (
	"fmt"
	"net/http"

	"github.com/Hustle299/Project-0/views"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// Create a signup page (New User)
func (u *Users) New(w http.ResponseWriter, R *http.Request) {
	err := u.NewView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

// function to submit the signup form - meaning create a new user
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Email is", form.Email)
	fmt.Fprintln(w, "Password is", form.Password)
}
