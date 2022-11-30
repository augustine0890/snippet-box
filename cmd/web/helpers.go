package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Create an newTemplateData() help, which returns a pointer to a templateData
// struct Initialised with the current year.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// The serverError helper writes an error message and stacktrace
// Then sends generic 500
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// This allows us to report the file name and line number of the file that
	// the error originated from (not this helper) "one step back"
	app.errorLog.Output(2, trace)
	app.errorLog.Panic(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError heler sends a specic status code and description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the correct template set from the cache based on the page
	// based on the name. If no entry point exists in the cache wih the name
	// provided, then create a new error and call the serverError
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// We are making the render a 2 step process
	// Step 1: make a trial render by writing the template to a buffer, if it fails respond with an error
	// Step 2: write the contents to our `http.ResponseWriter`

	// Initialise a new buffer
	buf := new(bytes.Buffer)

	// Write teh template to the buffer, instead of straight to the
	// http.ResponseWriter, if there is an serverError
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to write the HTTP status code
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note this is another
	// time where we pass our http.ResponseWriter to a function that takes an io.Writer
	buf.WriteTo(w)
}
