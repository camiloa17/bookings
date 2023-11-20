package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/camiloa17/bookings/internal/config"
	"github.com/camiloa17/bookings/internal/models"
	"github.com/justinas/nosurf"
)

//var templateCache = make(map[string]*template.Template)

var app *config.AppConfig
var pathToTemplates = "./templates"
var functions = template.FuncMap{}

// NewRenderer sets the config for the template package
func NewRenderer(appConfig *config.AppConfig) {
	app = appConfig
}

// AddDefaultData returns default data we want on every page
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func Template(w http.ResponseWriter, r *http.Request, templateName string, templateData *models.TemplateData) error {
	var templateCache map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		tc, err := CreateTemplateCache()

		if err != nil {
			templateCache = app.TemplateCache
			log.Println(err)
			return errors.New("can't get templates from cache")
		} else {
			templateCache = tc
		}

	}

	template, ok := templateCache[templateName]

	if !ok {
		log.Println("no template available in cache")
		return errors.New("can't get templates from cache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, r)

	err := template.Execute(buf, templateData)

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all files names *.page.gohtml

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	// range through all page files
	for _, page := range pages {
		// removes the path and leave the last part of the url
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		templatesRoute := fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates)
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
