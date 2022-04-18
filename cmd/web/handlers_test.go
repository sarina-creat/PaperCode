package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	// Create a new instance of our application struct. For now, this just
	// contains a couple of mock loggers (which discard anything written to
	// them).
	app := newTestApplication(t)

	// We then use the httptest.NewTLSServer() function to create a new test
	// server, passing in the value returned by our app.routes() method as the
	// handler for the server. This starts up a HTTPS server which listens on a
	// randomly-chosen port of your local machine for the duration of the test.
	// Notice that we defer a call to ts.Close() to shutdown the server when
	// the test finishes.

	//ts := httptest.NewTLSServer(app.routes())
	//defer ts.Close()
	//// The network address that the test server is listening on is contained
	//// in the ts.URL field. We can use this along with the ts.Client().Get()
	//// method to make a GET /ping request against the test server. This
	//// returns a http.Response struct containing the response.
	//rs, err := ts.Client().Get(ts.URL + "/ping")
	//if err != nil {
	//	t.Fatal(err)
	//}
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	statusCode, _, body := ts.get(t, "/ping")
	// Check
	if statusCode != 200 {
		t.Errorf("want %d, got %d", 200, statusCode)
	}

	if string(body) != "OK" {
		t.Errorf("want OK, got %s", string(body))
	}
}

func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct{
		name string
		urlPath string
		wantCode int
		wantBody []byte
	} {
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"NegativeId", "/snippet/-1", http.StatusNotFound, nil},
		{"DecimalID", "/snippet/1.23", http.StatusNotFound, nil},
		{"StringID", "/snippet/foo", http.StatusNotFound, nil},
		{"EmptyID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to got %d", tt.wantBody)
			}
		})
	}

}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name string
		userName string
		userEmail string
		userPassword string
		csrfToken string
		wantCode int
		wantBody []byte
	} {
		{"Valid submission", "Bob", "bob@example.com", "vaildPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field can not be blank")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field can not be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field can not be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Addresses is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			if code != tt.wantCode {
				t.Log(code)
				t.Errorf("want %d, got %d.", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Log(body)
				t.Errorf("want body %s to contain %q.", body, tt.wantBody)
			}
		})
	}
	//t.Log(csrfToken)
}

func TestCreatSnippetForm(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/snippet/creat")
		if code != http.StatusSeeOther {
			t.Errorf("want %d; got %d", http.StatusSeeOther, code)
		}

		if headers.Get("Location") != "/user/login" {
			t.Errorf("want %s; got %s", "/user/login", headers.Get("Location"))
		}
	})

	t.Run("Authenticated", func(t *testing.T) {
		//Authenticate the user...
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "")
		form.Add("csrf_token", csrfToken)
		ts.postForm(t, "/user/login", form)

		//被重定向了，测试无法通过。。。
		code, _, body := ts.get(t, "/snippet/creat")

		if code != http.StatusOK {
			t.Errorf("want %d, got %d", http.StatusOK, code)
		}

		formTag := "<form action='/snippet/creat' method='POST'"
		if !bytes.Contains(body, []byte(formTag)) {
			t.Errorf("want body %s to contain %q", body, formTag)
		}
	})

}
