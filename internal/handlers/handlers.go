package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/camiloa17/bookings/internal/config"
	"github.com/camiloa17/bookings/internal/forms"
	"github.com/camiloa17/bookings/internal/models"
	"github.com/camiloa17/bookings/internal/render"
)

// Repository pattern
var Repo *Repository

// Repository is teh repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(rs *Repository) {
	Repo = rs
}

// Home is the home page handler
func (rs *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	rs.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (rs *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	// get the ip from the session and add it to the stringMap
	remoteIP := rs.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	// send the template to the front end
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{StringMap: stringMap})
}

// Reservation is the reservation page handler
func (rs *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["reservation"] = models.Reservation{}

	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (rs *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)
	requiredValues := []string{"first_name", "last_name", "email", "phone"}
	form.Required(requiredValues...)

	form.MinLength("first_name", 10)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

}

// Generals is the generals page handler renders the room page
func (rs *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.gohtml", &models.TemplateData{})
}

// Majors is the majors page handler which renders the room page
func (rs *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.gohtml", &models.TemplateData{})
}

// Availability is the availability page handler
func (rs *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability is the availability page handler
func (rs *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handler request for availability and send JSON
func (rs *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact is the contact page handler
func (rs *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}
