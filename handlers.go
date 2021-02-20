package main

import (
	"fmt"
	"net/http"
	"strconv"

	"arnoldcodes.com/snippetbox/pkg/forms"
	"arnoldcodes.com/snippetbox/pkg/models"
)

//home handler
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) //use the notFound() helper
		return
	}
	//panic("oops! something went wrong")
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

	/*for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}
	//create an instance of a templateData struct holding the slice of snppets
	data := &templateData{Snippets: s}
	files := []string{
		"ui/html/home.tmpl",
		"ui/html/base.tmpl",
		"ui/html/footerPartial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) //use the serverError() helper
		return
	}
	ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err) //user the serverError() helper
	}*/
}

//View Snippet
func (app *Application) viewSnippet(w http.ResponseWriter, r *http.Request) {
	//Pat doesn't strip the colon(:) from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) //use the notFound() helper
		return
	}
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	//use the new render helper
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
	/*fmt.Fprint(w, "%v", s)*/
	/*
		//create an instance of templateData struct holding the snippet data
		data := &templateData{Snippet: s}
		//initialize a slice containing the paths to the pages to be parsed
		files := []string{
			"ui/html/show.tmpl",
			"ui/html/base.tmpl",
			"ui/html/footerPartial.tmpl",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, err) //use the serverError() helper
			return
		}
		ts.Execute(w, data)
		if err != nil {
			app.serverError(w, err) //user the serverError() helper
		}*/
}

//New createSnippetForm handler
func (app *Application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	//Pass a new empty form.Form object to the template
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
	/*err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	//Initializee a map to hold any validation errors.
	errors := make(map[string]string)

	//check that the title field is nt bank and is not more than 100 characters long
	// If it fails either of those checks, add a message to the errors
	//map using the field name as the key.
	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This is too long (maximum is 100 characters)"
	}
	// check that the content field isn't blank
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}
	//check the expires field isn't blank and matched one of the permitted
	//values ("1", "7" or "365").
	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	//If there are any errors, dump them in a plain text HTTP response and retrive
	//from the handler
	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})

		return
	}
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	//app.render(w, r, "create.page.tmpl", nil)*/
}

//Create a snippet
func (app *Application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//create a new forms.Form struct containing the POSTed data from the form, the use the validation
	//methods to check the content
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	// If the form isn't valid, redisplay the template passing in the form.Form
	//form.Form object as the data
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	//Because the form data (with type url.Values) has been anonymously embedded
	//in the form.Form struct, we can use the Get() method to retrieve
	//teh validated value for a particular form field.
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	//http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", id), http.StatusSeeOther)
	/*	//First we call r.ParseForm() which adds any data in POST request bodies
		//to the r.PostForm map. This also works in the same way for PUT nad PATCH request.
		//If there are any errors, we use or app.ClietError helper to a 400 Bad Request response to the user.
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		//Use the r.PostForm.Get() method to retrieve the relevant data fields
		// from the r.PostForm map.
		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")
		expires := r.PostForm.Get("expires")

		//validate form inputs
		errors := make(map[string]string)

		//validate title input
		if strings.TrimSpace(title) == "" {
			errors[title] = "This field cannot be empty"
		} else if utf8.RuneCountInString(title) > 100 {
			errors[title] = "Title too long (only a maximum of 100 characters is allowed)"
		}
		if strings.TrimSpace(content) == "" {
			errors[content] = "this field cannot be empty"
		}
		if strings.TrimSpace(expires) == "" {
			errors[expires] = "This field cannot be empty"
		} else if expires != "365" && expires != "7" && expires != "1" {
			errors[expires] = "This field is invalid"
		}

		//If there are any validation errors, re-display the create.page.tmpl
		//template passing in the validation errors and previously submitted
		//r.PostForm data.
		if len(errors) > 0 {
			app.render(w, r, "create.page.tmpl", &templateData{
				FormErrors: errors,
				FormData:   r.PostForm,
			})
			return
		}

		//Create a new snippet record in the database using the fom data
		id, err := app.snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	*/
	/*if r.Method != "POST" {
		w.Header().Set("Allow", "Post")
		app.clientError(w, http.StatusMethodNotAllowed) //use the clientError() helper
		return
	}
	//some variables holding dummy data
	title := "Hey Arnold"
	content := "Exe,\n my name is Arnold,\n and I am an aspiring Go and Java programmer!\n\n-arnoldcodes.com"
	expires := "7"

	//pass the data to the SnippetModel.Insert() method, receiving the
	//ID of the new record back.

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	*/
}
