package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Hustle299/Project-0/context"
	"github.com/Hustle299/Project-0/models"
	"github.com/Hustle299/Project-0/views"
	"github.com/gorilla/mux"
)

const (
	ShowGallery     = "show_gallery"
	IndexGalleries  = "index_galleries"
	EditGallery     = "edit_gallery"
	maxMultipartMem = 1 << 20
)

type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryServices
	r         *mux.Router
	is        models.ImageService
}

// tao mot gallery instance trong controller
func NewGalleries(gs models.GalleryServices, is models.ImageService, r *mux.Router) *Galleries {
	return &Galleries{
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		is:        is,
		r:         r,
	}
}

// struct de nhan cac du lieu nhap vao tu gallery form
type GalleryForm struct {
	Title string `schema:"title"`
}

// function submitform tao 1 galleries moi (POST /galleries).
// su dung parseForm function truyen data tu web request vao GalleryForm tao o tren
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}

	user := context.User(r.Context())

	//tao gallery moi voi ten va user ID
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	//goi GalleryService’s Create method truyen du lieu tu instance gallery dang co vao DB
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}

	//redirect sau khi tao xong
	url, err := g.r.Get(EditGallery).URL("id", strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// GET /galleries/:id (xem 1 gallery)
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {

	//we reuse the function "galleryByID"
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// The galleryByID method will already render the error
		// for us, so we just need to return here.
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, r, vd)
}

// tim gallery theo ID
func (g *Galleries) galleryByID(w http.ResponseWriter,
	r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	images, _ := g.is.ByGalleryID(gallery.ID)
	gallery.Images = images
	return gallery, nil
}

// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		// If there is an error we are going to render the
		// EditView again with an alert message.
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	// If there is an error our alert will be an error. Otherwise
	// we will still render an alert, but instead it will be
	// a success message.
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Gallery successfully updated!",
		}
	}
	// Error or not, we are going to render the EditView with
	// our updated information.
	g.EditView.Render(w, r, vd)
}

func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	// Lookup the gallery using the galleryByID we wrote earlier
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// If there is an error the galleryByID will have rendered
		// it for us already.
		return
	}
	// We also need to retrieve the user and verify they have
	// permission to delete this gallery. This means we will
	// need to use the RequireUser middleware on any routes
	// mapped to this method.
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		// If there is an error we want to set an alert and
		// render the edit page with the error. We also need
		// to set the Yield to gallery so that the EditView
		// is rendered correctly.
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
		return
	}
	// redirect to the index page
	url, err := g.r.Get(IndexGalleries).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}

// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	//send data galleries to yield template
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

// POST galleries/:id/images (upload image)
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		// If we can't parse the form just render an error alert on the// edit gallery page.
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	//iterates over uploaded file
	files := r.MultipartForm.File["images"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()

		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Images successfully uploaded!",
	}
	g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/images/:filename/delete
func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery or image", http.StatusForbidden)
		return
	}
	// Get the filename from the path
	filename := mux.Vars(r)["filename"]
	// Build the Image model
	i := models.Image{Filename: filename,
		GalleryID: gallery.ID,
	}
	// Try to delete the image.
	err = g.is.Delete(&i)
	if err != nil {
		// Render the edit page with any errors.
		var vd views.Data
		vd.Yield = gallery
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	// If all goes well, redirect to the edit gallery page.
	url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}
