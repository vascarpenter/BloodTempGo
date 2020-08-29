package main

import (
	"BloodTemp/m/routes"
	"html/template"
	"io"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	_ "github.com/mattn/go-oci8"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Template はHTMLテンプレートを利用するためのRenderer Interfaceです。
type Template struct {
	templates *template.Template
}

// Render はHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込みます。
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Echo instance
	e := echo.New()

	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
	}
	t := &Template{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} uri=${uri} path=${path} status=${status}\n",
	}))
	e.Use(middleware.Recover())

	var keystore string
	if keystore = os.Getenv("COOKIE_SEED"); keystore == "" {
		keystore = "secret key"
	}
	var store = sessions.NewCookieStore([]byte(keystore))
	e.Use(session.Middleware(store))

	// Routes
	e.Static("/css", "./static/css")
	e.Static("/img", "./static/img")
	e.Static("/javascript", "./static/javascript")
	e.Static("/icon", "./static/icon")
	e.GET("/", routes.IndexRouter)
	e.POST("/", routes.IndexRouterPost)

	// Start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "3001"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
