package main

import (
	"alexedwards.net/snippetbox/pkg/forms"
	"alexedwards.net/snippetbox/pkg/models"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	app.notFound(w) // Use the NotFound() helper
	//	//http.NotFound(w, r)
	//	return
	//}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Use the new render helper.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

	//for _, snippet := range s {
	//	fmt.Fprintf(w, "%v\n", snippet)
	//}

	//data := &templateData{Snippets: s}
	//
	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//	//app.errorLog.Println(err.Error())
	//	//http.Error(w, "Internal Server Error",500)
	//	return
	//}
	//
	//err = ts.Execute(w, data)
	//if err != nil {
	//	app.serverError(w, err)
	//	//app.errorLog.Println(err.Error())
	//	//http.Error(w, "Internal Server Error",500)
	//}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		//http.NotFound(w, r)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Use the PopString() method to retrieve the value for the "flash" key.
	// PopString() also deletes the key and value from the session data, so it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.

	//flash := app.session.PopString(r, "flash")

	//Use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		//Flash: flash,
		Snippet: s,
	})

	//data := &templateData{Snippet:s}
	//
	//file := []string{
	//	"./ui/html/show.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//// Parse the template files...
	//ts, err := template.ParseFiles(file...)
	//if err != nil {
	//	app.serverError(w, err)
	//	//app.errorLog.Println(err.Error())
	//	//http.Error(w, "Internal Server Error",500)
	//	return
	//}
	//
	//err = ts.Execute(w, data)
	//if err != nil {
	//	app.serverError(w, err)
	//}

	//fmt.Fprintf(w, "%v", s)
	//w.Write([]byte("Display a specific Snippet..."))
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	//w.WriteHeader(http.StatusMethodNotAllowed)
	//	//w.Write([]byte("Method Not Allowed"))
	//
	//	//Use the http.Error() function to send a 405 status code and "Method
	//	// Not Alloowed" string as the response body.
	//
	//	//Suppressing System-Generated Headers
	//	//w.Header()["Date"] = nil
	//
	//	app.clientError(w, http.StatusMethodNotAllowed)
	//	//http.Error(w, "Method Not Allowed", 405)
	//	return
	//}

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError helper to send
	// a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the relevant data fields
	// from the r.PostForm map.
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "creat.page.tmpl", &templateData{Form: form})
		return
	}
	//title := r.PostForm.Get("title")
	//content := r.PostForm.Get("content")
	//expires := r.PostForm.Get("expires")

	////Initialize a map to hold any validation errors.
	//errors := make(map[string]string)

	////Validate the data is correct or not.
	//// Check that the title field is not blank and is not more than 100 characters
	//if strings.TrimSpace(title) == "" {
	//	errors["title"] = "this field cannot be blank."
	//} else if utf8.RuneCountInString(title) > 100 {
	//	errors["title"] = "This field is too long (maximum is 100 characters."
	//}
	//
	////Check that the content field isn't blank.
	//if strings.TrimSpace(content) == "" {
	//	errors["content"] = "this field cannot be blank."
	//}
	//
	//// Check the expires field isn't blank and matches one of the permitted
	//// values ("1", "7" or "365").
	//if strings.TrimSpace(expires) == "" {
	//	errors["expires"] = "this field cannot be blank"
	//} else if expires != "365" && expires != "7" && expires != "1" {
	//	errors["expires"] = "this field is invalid"
	//}
	//
	//if len(errors) > 0 {
	//	app.render(w, r, "creat.page.tmpl", &templateData{
	//		FormErrors: errors,
	//		FormData: r.PostForm,
	//	})
	//	return
	//}

	////database
	//title := "0 snail"
	//content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	//expires := "7"

	id, err := app.snippets.Insert(form.Get("title"),form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Use the Put() method to add a string value ("Your snippet was saved
	// successfully!") and the corresponding key ("flash") to the session
	// data. Note that if there's no existing session for the current user
	// (or their session has expired) then a new, empty, session for them
	// will automatically be created by the session middleware.
	app.session.Put(r, "false", "Snippet successfully created!")
	// Change the redirect to use the new semantic URL style of /snippet/:id
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	//w.Write([]byte("Creat a new  snippet..."))
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	//Retrieve the appropriate template set from the cache based on the page name
	//(like 'home.page.tmpl'). If no entry exists in the cache with the
	//provided name, call the serverError helper method that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	//Initialize a new buffer
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	buf.WriteTo(w)
	//time.Sleep(time.Second * 5)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "creat.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	//Parse the form data and
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	//If there are any errors, redisplay teh signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// Try to create a new user record in the database. If the email already exists
	// add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Addresses is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked and asking them to log in.
	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrecct")

			app.render(w, r, "login.page.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Add the ID of the current user to the session, so that they are now 'logged
	// in'.
	app.session.Put(r, "authenticateUserID", id)

	if url := app.session.PopString(r, "requestURL"); url != "" {
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}
	// Redirect the user to the create snippet page
	http.Redirect(w, r, "/snippet/creat", http.StatusSeeOther)
	log.Println("Done...............")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}


func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "about.page.tmpl", nil)
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")

	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "profile.page.tmpl", &templateData{
		User: user,
	})
}

func (app *application) changePasswordForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "password.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	//Parse the form data and
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper we made earlier.
	form := forms.New(r.PostForm)
	form.Required("currentPassword", "newPassword", "newPasswordConfirmation")
	form.MinLength("newPassword", 10)
	if form.Get("newPassword") != form.Get("newPasswordConfirmation") {
		form.Errors.Add("newPasswordConfirmation", "Password do not match")
	}

	//If there are any errors, redisplay teh signup form.
	if !form.Valid() {
		app.render(w, r, "password.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	userID := app.session.GetInt(r, "authenticatedUserID")
	err = app.users.ChangePassword(userID, form.Get("currentPassword"), form.Get("newPassword"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("currentPassword", "entered the wrong password")
			app.session.Put(r, "flash", "Password has been successfully changed!")
			}
		} else {
			app.serverError(w, err)
			return
		}
	app.session.Put(r, "flash", "Password has been successfully changed!")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
