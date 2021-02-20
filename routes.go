package main

import (
	"net/http"
	//use gorillamux
)

//use http.Handler instead of *http.ServeMux

func (app *Application) routes() http.Handler {
	//create a middleware chain containing our 'standard' middleware
	//which will be used for every request our application receives
	//standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	/*	//using pat instead of gorillamux ro route handlers
		mux := pat.New()
		mux.Get("/", http.HandlerFunc(app.home))
		mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
		mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
		mux.Get("/snippet/:id", http.HandlerFunc(app.viewSnippet)) //moved down
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("post/create", app.createSnippetForm)

	mux.HandleFunc("/post/create", app.createSnippet)
	mux.HandleFunc("/post", app.viewSnippet)

	fileServer := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	/*fileServer := http.FileServer(http.Dir("ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))*/

	//return standardMiddleware.Then(mux)
	//wrap the existing chain wil the recoverPanic middleware
	//wrap the existing chain with the logRequest middleware
	//Pass the servemux as the 'next' parameter to the secureHeaders middleware
	//because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	return app.logRequest(secureHeaders(mux))
}
