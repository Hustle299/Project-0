package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

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
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}
