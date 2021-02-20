package main

import (
	"html/template"
	"path/filepath"
	"time"

	"arnoldcodes.com/snippetbox/pkg/forms"
	"arnoldcodes.com/snippetbox/pkg/models"
)

//define a templateData type to act as the holding structure for any
//dynamic data that we want to pass to our HTML templates.

type templateData struct {
	CurrentYear int
	Form        *forms.Form
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

//create a humanDate function which return a nicely formatted string
//representation of a time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")

}

//Initialize a template.FuncMap object and store it in a global variable
//This is essentially a string-keyed map which acts as a lookup between the names of a
//custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	//Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	//Use the filepath.Glob function to get a slice of all filepath with
	//the extension '.page.tmpl'. This essentially gives us a slice of all the
	//'page' tesmplates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	//loop through the pages one-by-one
	for _, page := range pages {
		//Extract the file name (like 'home.page.tmpl') from the full file path
		//and assign it to the name variable
		name := filepath.Base(page)

		//The template.FuncMap must be registered with the template set before
		//call ParseFiles() method. This means we have to use template.New()
		//create an empty template set, use the Funcs() method to register the
		//template.FuncMap, and then parse the files as normal
		// Parse the page template file in to a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		//Use the ParseGlob method to add any 'layout' template to the
		//template set (in our case, it' just the 'base' layout at the moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		//Use the ParseGlob method to add any 'partial' templates to the
		//template set (in our case, it's just the 'footer' partial at the moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		//Add the template set to the cache, using the name of the page
		//(like 'home.page.tmpl) as the key
		cache[name] = ts

	}
	//return the map.
	return cache, nil
}
