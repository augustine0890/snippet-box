package main

import (
	"html/template"
	"path/filepath"
	"snippetbox/pkg/models"
	"time"
)

type templateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialise a template.FuncMap object to store it in a global variable.
// This is essentially a string-keyed map which actsas a lookup between
// the names of our custom funcs and the funcs themselves
// NOTE: custom template funcs can acceptas many params as needed
// but must only return a single value except for an error
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	// Initialise a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths
	// that match the pattern "./ui/html/pages/*.tmpl". This will eventually
	// gives us a slice of all the filepaths for our application page templates
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name and assign it to the name variable
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before
		// you call the ParseFiles(). The means we have to use template.New() to
		// create an empty template set, use the Funcs() to register the template.FuncMap
		// then pass the file as normal
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		cache[name] = ts
	}

	return cache, nil
}
