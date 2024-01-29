package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	rocketApp "github.com/sinantomruk/miltron-launchsite/app"
	"github.com/sinantomruk/miltron-launchsite/ui"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	myApp := rocketApp.NewApp("http://0.0.0.0:5000", "API_KEY_1")
	a := fyneApp.New()
	w := a.NewWindow("Rocket Launcher")
	w.Resize(fyne.NewSize(1024, 720))

	weatherContainer, err := ui.NewWeatherContainer(ctx, myApp.Weather)
	if err != nil {
		fmt.Println(err)
		return
	}
	go weatherContainer.UpdateData(ctx)

	rocketContainer, err := ui.NewRocketsContainer(ctx, myApp, w)
	if err != nil {
		fmt.Println(err)
		return
	}
	cont := container.NewHSplit(rocketContainer.Container, weatherContainer.Container)

	w.SetContent(cont)
	w.ShowAndRun()
}
