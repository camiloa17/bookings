package handlers

import (
	"net/http"

	"github.com/camiloa17/bookings/pkg/config"
	"github.com/camiloa17/bookings/pkg/models"
	"github.com/camiloa17/bookings/pkg/render"
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
	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{})
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
	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{StringMap: stringMap})
}

// Reservation is the reservation page handler
func (rs *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.gohtml", &models.TemplateData{})
}

// Generals is the generals page handler renders the room page
func (rs *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.gohtml", &models.TemplateData{})
}

// Majors is the majors page handler which renders the room page
func (rs *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.gohtml", &models.TemplateData{})
}

// Availability is the availability page handler
func (rs *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.gohtml", &models.TemplateData{})
}

// Contact is the contact page handler
func (rs *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.gohtml", &models.TemplateData{})
}
