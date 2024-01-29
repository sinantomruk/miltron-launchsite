package ui

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sinantomruk/miltron-launchsite/app"
	"github.com/sinantomruk/miltron-launchsite/models"
)

const timeLayout = "2006-01-02T15:04:05.999999"

type RocketController interface {
	Rockets() ([]*models.RocketInfo, error)
	LaunchRocket(rocketID string) (*models.RocketInfo, error)
	DeployRocket(rocketID string) (*models.RocketInfo, error)
	CancelLaunch(rocketID string) (*models.RocketInfo, error)
}

type RocketsContainer struct {
	rocketController RocketController
	Container        *widget.Accordion
	window           fyne.Window
}

func updateDetails(ctx context.Context, telemetryChan <-chan models.TelemetryData, labels map[string]*widget.Label) {
	for {
		data := <-telemetryChan
		labels["altitude"].SetText(fmt.Sprintf("Altitude: %.2f", data.Altitude))
		labels["speed"].SetText(fmt.Sprintf("Speed: %.2f", data.Speed))
		labels["acceleration"].SetText(fmt.Sprintf("Acceleration: %.2f", data.Acceleration))
		labels["thrust"].SetText(fmt.Sprintf("Thrust: %.2f", data.Thrust))
		labels["temperature"].SetText(fmt.Sprintf("Temperature: %.2f", data.Temperature))
	}
}

func (rc *RocketsContainer) onButtonClick(buttons map[string]*widget.Button, button string, statusLabel *widget.Label, historyDetails *widget.Form, rocketId string) func() {
	return func() {
		var ri *models.RocketInfo
		var err error
		if button == "deploy" {
			ri, err = rc.rocketController.DeployRocket(rocketId)
			if err != nil {
				return
			}
			buttons["launch"].Show()
			buttons["cancel"].Hide()
			buttons["deploy"].Hide()
		} else if button == "launch" {
			ri, err = rc.rocketController.LaunchRocket(rocketId)
			if err != nil {
				return
			}
			buttons["launch"].Hide()
			buttons["cancel"].Show()
			buttons["deploy"].Hide()
		} else if button == "cancel" {
			ri, err = rc.rocketController.CancelLaunch(rocketId)
			if err != nil {
				return
			}
			buttons["launch"].Hide()
			buttons["cancel"].Hide()
			buttons["deploy"].Hide()
		}
		if ri != nil {
			historyDetails.Items = []*widget.FormItem{}
			for k, v := range ri.Timestamps {
				if v != "" {
					parsedTime, err := time.Parse(timeLayout, v)
					if err != nil {
						fmt.Println(err)
					}
					parsedTime = parsedTime.Local()
					historyDetails.Append(k, widget.NewLabel(parsedTime.Format(time.RFC1123)))
				}
			}
			statusLabel.SetText(ri.Status)
		}
	}
}

func (rc *RocketsContainer) showDetails(ri *models.RocketInfo) func() {
	return func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		labels := map[string]*widget.Label{
			"description":  widget.NewLabel(fmt.Sprintf("Description: %s", ri.Payload.Description)),
			"mass":         widget.NewLabel(fmt.Sprintf("Mass: %.2f", ri.Mass)),
			"weight":       widget.NewLabel(fmt.Sprintf("Weight: %d", ri.Payload.Weight)),
			"altitude":     widget.NewLabel(fmt.Sprintf("Altitude: %.2f", ri.Altitude)),
			"speed":        widget.NewLabel(fmt.Sprintf("Speed: %.2f", ri.Speed)),
			"acceleration": widget.NewLabel(fmt.Sprintf("Acceleration: %.2f", ri.Acceleration)),
			"thrust":       widget.NewLabel(fmt.Sprintf("Thrust: %.2f", ri.Thrust)),
			"temperature":  widget.NewLabel(fmt.Sprintf("Temperature: %.2f", ri.Temperature)),
		}

		telemetryChan := make(chan models.TelemetryData)
		go app.ReadTelemetry(ctx, ri.Telemetry.Host, ri.Telemetry.Port, telemetryChan)
		go updateDetails(ctx, telemetryChan, labels)

		statusLabel := widget.NewLabel(ri.Status)

		historyDetails := widget.NewForm()
		for k, v := range ri.Timestamps {
			if v != "" {
				parsedTime, err := time.Parse(timeLayout, v)
				if err != nil {
					fmt.Println(err)
				}
				parsedTime = parsedTime.Local()
				historyDetails.Append(k, widget.NewLabel(parsedTime.Format(time.RFC1123)))
			}
		}
		history := widget.NewAccordion(widget.NewAccordionItem("History", historyDetails))

		rows := container.NewVBox(
			labels["description"],
			labels["mass"],
			labels["weight"],
			labels["altitude"],
			labels["speed"],
			labels["acceleration"],
			labels["thrust"],
			labels["temperature"],
			container.NewGridWithColumns(2, statusLabel, history),
		)

		details := dialog.NewCustomWithoutButtons(ri.Model, rows, rc.window)
		details.Resize(fyne.NewSize(400, 400))

		closeFunc := func() {
			details.Hide()
			rc.refreshContainer()
		}

		buttons := map[string]*widget.Button{
			"launch": widget.NewButton("Launch", details.Hide),
			"deploy": widget.NewButton("Deploy", details.Hide),
			"cancel": widget.NewButton("Cancel", details.Hide),
			"close":  widget.NewButton("Close", closeFunc),
		}

		if ri.Status == models.RocketStatusWaiting {
			buttons["launch"].Hide()
			buttons["cancel"].Hide()
			buttons["deploy"].Show()
		} else if ri.Status == models.RocketStatusDeployed {
			buttons["launch"].Show()
			buttons["cancel"].Hide()
			buttons["deploy"].Hide()
		} else if ri.Status == models.RocketStatusLaunched {
			buttons["launch"].Hide()
			buttons["cancel"].Show()
			buttons["deploy"].Hide()
		} else if ri.Status == models.RocketStatusCancelled {
			buttons["launch"].Hide()
			buttons["cancel"].Hide()
			buttons["deploy"].Hide()
		}

		buttons["launch"].OnTapped = rc.onButtonClick(buttons, "launch", statusLabel, historyDetails, ri.ID)
		buttons["cancel"].OnTapped = rc.onButtonClick(buttons, "cancel", statusLabel, historyDetails, ri.ID)
		buttons["deploy"].OnTapped = rc.onButtonClick(buttons, "deploy", statusLabel, historyDetails, ri.ID)
		details.SetButtons([]fyne.CanvasObject{buttons["cancel"], buttons["launch"], buttons["deploy"], buttons["close"]})

		details.Show()
	}
}

func (rc *RocketsContainer) refreshContainer() {
	rockets, err := rc.rocketController.Rockets()
	if err != nil {
		fmt.Println(err)
		return
	}
	items := []*widget.AccordionItem{}

	index := -1
	for i, item := range rc.Container.Items {
		if item.Open {
			index = i
		}
	}

	for i, rocket := range rockets {
		rocketStatusText := canvas.NewText(rocket.Status, color.White)
		if rocket.Status == models.RocketStatusWaiting {
			rocketStatusText.Color = color.RGBA{189, 116, 8, 255}
		} else if rocket.Status == models.RocketStatusLaunched || rocket.Status == models.RocketStatusDeployed {
			rocketStatusText.Color = color.RGBA{0, 140, 6, 255}
		} else if rocket.Status == models.RocketStatusFailed || rocket.Status == models.RocketStatusCancelled {
			rocketStatusText.Color = color.RGBA{239, 24, 24, 255}
		}

		desc := widget.NewLabel(rocket.Payload.Description)
		desc.Wrapping = fyne.TextWrapWord

		detail := container.NewVBox(
			desc,
			container.NewHBox(rocketStatusText, widget.NewButton("Show Details", rc.showDetails(rocket))),
		)
		item := widget.NewAccordionItem(rocket.Model, detail)
		if i == index {
			item.Open = true
		}
		items = append(items, item)
	}

	rc.Container.Items = items
	rc.Container.Refresh()
}

func NewRocketsContainer(ctx context.Context, rocketController RocketController, w fyne.Window) (*RocketsContainer, error) {
	rc := RocketsContainer{
		rocketController: rocketController,
		window:           w,
	}
	rockets, err := rc.rocketController.Rockets()
	if err != nil {
		return nil, err
	}
	items := []*widget.AccordionItem{}

	for _, rocket := range rockets {
		rocketStatusText := canvas.NewText(rocket.Status, color.White)
		if rocket.Status == models.RocketStatusWaiting {
			rocketStatusText.Color = color.RGBA{189, 116, 8, 255}
		} else if rocket.Status == models.RocketStatusLaunched || rocket.Status == models.RocketStatusDeployed {
			rocketStatusText.Color = color.RGBA{0, 140, 6, 255}
		} else if rocket.Status == models.RocketStatusFailed || rocket.Status == models.RocketStatusCancelled {
			rocketStatusText.Color = color.RGBA{239, 24, 24, 255}
		}

		desc := widget.NewLabel(rocket.Payload.Description)
		desc.Wrapping = fyne.TextWrapWord

		detail := container.NewVBox(
			desc,
			container.NewHBox(rocketStatusText, widget.NewButton("Show Details", rc.showDetails(rocket))),
		)
		items = append(items, widget.NewAccordionItem(rocket.Model, detail))
	}

	rc.Container = widget.NewAccordion(items...)

	return &rc, nil
}
