package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/camiloa17/bookings/pkg/config"
	"github.com/camiloa17/bookings/pkg/models"
)

//var templateCache = make(map[string]*template.Template)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(appConfig *config.AppConfig) {
	app = appConfig
}

// AddDefaultData returns default data we want on every page
func AddDefaultData(templateDate *models.TemplateData) *models.TemplateData {
	return templateDate
}

func RenderTemplate(responseWriter http.ResponseWriter, templateName string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		tc, err := CreateTemplateCache()

		if err != nil {
			templateCache = app.TemplateCache
			log.Println(err)
		} else {
			templateCache = tc
		}

	}

	template, ok := templateCache[templateName]

	if !ok {
		log.Fatal("no template available in cache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	err := template.Execute(buf, templateData)

	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(responseWriter)

	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all files names *.page.gohtml

	pages, err := filepath.Glob("./templates/*.page.gohtml")

	if err != nil {
		return myCache, err
	}

	// range through all page files
	for _, page := range pages {
		// removes the path and leave the last part of the url
		name := filepath.Base(page)

		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		templatesRoute := "./templates/*.layout.gohtml"
		matches, err := filepath.Glob(templatesRoute)
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(templatesRoute)
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet

	}

	return myCache, nil
}
