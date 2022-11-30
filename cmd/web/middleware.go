package main

import (
	"fmt"
	"net/http"
)

// Becaise we want this to run on ALL requests we must execute it
// BEFORE the servemuc
// secureHeaders -> servemux -> application handler
// func secureHeaders(next http.Handler) http.Handler {
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		// w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		// w.Header().Set("X-Content-Type-Options", "nosniff")
		// w.Header().Set("X-Frame-Options", "deny")
		// w.Header().Set("X-XXS-Protection", "0")

		// next.ServeHTTP(w, r)
	// })
// }

// Flow
// logRequest <-> secureHeaders <-> servemux <-> application handler
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic as
		// Go unwinds the stack)
		defer func() {
			// use builtin recover function to check if there has been a panic
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}