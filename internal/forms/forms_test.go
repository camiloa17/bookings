package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFrom_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when should have been valid.")
	}
}

func TestForm_Required(t *testing.T) {
	posteData := url.Values{}
	form := New(posteData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form show valid when required fields missing")
	}

	posteData = url.Values{}
	posteData.Add("a", "a")
	posteData.Add("b", "b")
	posteData.Add("c", "c")

	form = New(posteData)
	form.Required("a", "b", "c")

	isValid := form.Valid()

	if !isValid {
		t.Error("form show invalid when required fields are present")
	}
}

func TestForm_MinLength(t *testing.T) {
	posteData := url.Values{}
	posteData.Add("c", "ccc")

	form := New(posteData)

	form.MinLength("c", 4)
	if form.Valid() {
		t.Error("form show valid when minLength should be 4 and it was 3")
	}
	isError := form.Errors.Get("c")
	if isError == "" {
		t.Error("should have an error but did not get one.")
	}
	posteData = url.Values{}
	posteData.Add("a", "a")
	posteData.Add("b", "bb")

	form = New(posteData)

	form.MinLength("b", 2)
	isValid := form.Valid()

	isError = form.Errors.Get("b")
	if isError != "" {
		t.Error("should not have an error but did get one.")
	}
	if !isValid {
		t.Error("form show invalid when minLength should be ok")
	}
}

func TestForm_IsEmail(t *testing.T) {
	posteData := url.Values{}
	posteData.Add("c", "ccc")
	form := New(posteData)

	form.IsEmail("c")
	if form.Valid() {
		t.Error("form show valid when c should not be an email")
	}

	posteData = url.Values{}
	posteData.Add("b", "bs@gmail.com")

	form = New(posteData)

	form.IsEmail("b")
	isValid := form.Valid()

	if !isValid {
		t.Error("form show invalid when minLength should be ok")
	}
}

func TestForm_Has(t *testing.T) {
	posteData := url.Values{}
	form := New(posteData)

	isValuePresent := form.Has("a")
	if isValuePresent {
		t.Error("has show present when it should be missing")
	}

	posteData = url.Values{}
	posteData.Add("a", "a")

	form = New(posteData)
	isValuePresent = form.Has("a")

	if !isValuePresent {
		t.Error("has show present when it should be  present")
	}
}
