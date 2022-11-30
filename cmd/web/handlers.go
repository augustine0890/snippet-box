package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/pkg/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Initialise a slice containing the paths to 2 files. The ORDER MATTERS
	// The base template must be first
	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }
	// // Use the template.ParseFiles function to read the template file into a
	// // template set.
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// 	return
	// }
	//
	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// we then execute the methong on the template.
	// the last param to Execute() is the dynamic data
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// }

	// Call the newTemplateDate() helper to get a templateDate struct
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the new render helper
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

// snippetView handler
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	// Extract the value of the id
	// convert it to an interger using `strconv.Atoi`, if it cant convert it
	// or its less that 1 return 404
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		// Current best practice for checking for specific error types
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
