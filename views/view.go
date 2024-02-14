package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/Hustle299/Project-0/context"
)

// This file is to render out a new view combining the main page and the layouts

// A struct to hold a View consist of a template and layouts which later on will
// be passed variables in NewView function below to create a page
type View struct {
	Template *template.Template
	Layout   string
}

// Variables used in layoutFiles function
var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

// This function is used to collect all layout gohtml files in layout folder
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// Add the prefix "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// Add the suffix ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

// Create a new View storing the parsed template ready to be template-execute
func NewView(layout string, files ...string) *View {
	// Add the prefix and the suffix
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
			"persists, please email support@lenslocked.com",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}
