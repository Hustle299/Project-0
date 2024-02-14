package main

import (
	"fmt"
	"net/http"

	"github.com/Hustle299/Project-0/controllers"
	"github.com/Hustle299/Project-0/middleware"
	"github.com/Hustle299/Project-0/models"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "duytung299"
	dbname   = "photogarage_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	services.AutoMigrate()

	r := mux.NewRouter()

	// Tao cac controller
	staticsC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	//Tao 1 middleware
	userMw := middleware.User{
		UserServices: services.User,
	}
	requireUserMw := middleware.RequireUser{}

	// galleriesC.New la 1 handler nen dung apply
	newGallery := requireUserMw.Apply(galleriesC.New)
	// galleriecsC.Create la 1 http.HandlerFunc dung applyfn
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)

	//cac duong dan cho trang tinh
	r.Handle("/", staticsC.Home).Methods("GET")
	r.Handle("/contact", staticsC.Contact).Methods("GET")

	//cac duong dan lien quan user
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	//cac duong dan lien quan den galleries
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries/new", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET").Name(controllers.IndexGalleries)
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")

	//duong dan lien quan image
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	http.ListenAndServe(":8080", userMw.Apply(r))
}
