package app

import (
	"github.com/wailsapp/wails"
)

func basic() string {
	return "Hello World!"
}

func Run() error {

	//js := mewn.String("./frontend/dist/app.js")
	//css := mewn.String("./frontend/dist/app.css")
	wails.BuildMode = "bridge"
	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    768,
		Title:     "onecup",
		JS:        "./app/frontend/dist/app.js",
		CSS:       "./app/frontend/dist/app.css",
		Colour:    "#131313",
		Resizable: true,
	})

	app.Bind(basic)
	err := app.Run()
	if err != nil {
		return err
	}

	return nil
}
