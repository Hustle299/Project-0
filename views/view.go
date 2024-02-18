package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/Hustle299/Project-0/context"
)

// Render ra 1 page bang file chinh va cac layout

// Struct giu 1 template va cac layout
// dung NewView function de tao 1 trang
type View struct {
	Template *template.Template
	Layout   string
}

// Variables su dung layoutFiles function
var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// gom toan bo file layout
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// them dau "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// them duoi".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

// Tao 1 View chua template de dung ham template excecute
func NewView(layout string, files ...string) *View {
	// Add dau duoi
	addTemplatePath(files)
	addTemplateExt(files)

	// The slice of string variables: "files" - is used to hold names of
	// the template we need to use in ParseFiles function below including layouts
	files = append(files, layoutFiles()...)

	// Parse all the files we collect and store it in variable "t"
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// Excecute the "View" we create in the NewView function
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		// We need to do this so we can access the data in a var
		// with the type Data.
		vd = d
	default:
		// If the data IS NOT of the type Data, we create one
		// and set the data to the Yield field like before.
		vd = Data{
			Yield: data,
		}
	}
	// Lookup and set the user to the User field
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(w, "Something went wrong. If the problem "+
			"persists, please email teo2992001@gmail.com",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}
