package ui

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sinantomruk/miltron-launchsite/models"
)

type WeatherContainer struct {
	dataFetcher func() (*models.Weather, error)
	Container   *fyne.Container
	// Weather Info
	temperature *widget.Label
	humidity    *widget.Label
	pressure    *widget.Label
	time        *widget.Label
	// Precipitation
	probability *widget.Label
	rain        *widget.Label
	snow        *widget.Label
	sleet       *widget.Label
	hail        *widget.Label
	// Wind
	direction *widget.Label
	angle     *widget.Label
	speed     *widget.Label
}

func NewWeatherContainer(ctx context.Context, dataFetcher func() (*models.Weather, error)) (*WeatherContainer, error) {
	wc := &WeatherContainer{
		dataFetcher: dataFetcher,
		temperature: widget.NewLabel(""),
		humidity:    widget.NewLabel(""),
		pressure:    widget.NewLabel(""),
		time:        widget.NewLabel(""),
		probability: widget.NewLabel(""),
		rain:        widget.NewLabel(""),
		snow:        widget.NewLabel(""),
		sleet:       widget.NewLabel(""),
		hail:        widget.NewLabel(""),
		direction:   widget.NewLabel(""),
		angle:       widget.NewLabel(""),
		speed:       widget.NewLabel(""),
	}

	precipitation := widget.NewFormItem("Precipitation:", container.NewVBox(
		wc.probability,
		wc.rain,
		wc.snow,
		wc.sleet,
		wc.hail,
	))

	wind := widget.NewFormItem("Wind:", container.NewVBox(
		wc.direction,
		wc.angle,
		wc.speed,
	))

	f := widget.NewForm(precipitation, wind)

	title := canvas.NewText("Weather", color.White)
	title.TextSize = 36
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	wc.Container = container.NewVBox(
		title,
		wc.temperature,
		wc.humidity,
		wc.pressure,
		wc.time,
		f,
	)
	return wc, nil
}

func (wc *WeatherContainer) UpdateData(ctx context.Context) {
	for {
		wi, err := wc.dataFetcher()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Millisecond * 100)
			continue
		}

		parsedTime, err := time.Parse(timeLayout, wi.Time)
		if err != nil {
			fmt.Println(err)
		}
		parsedTime = parsedTime.Local()

		wc.temperature.SetText(fmt.Sprintf("Temperature: %.2f°C", wi.Temperature))
		wc.humidity.SetText(fmt.Sprintf("Humidity: %.2f%%", wi.Humidity*100))
		wc.pressure.SetText(fmt.Sprintf("Pressure: %.2f", wi.Pressure))
		wc.time.SetText(fmt.Sprintf("Time: %s", parsedTime.Format(time.RFC1123)))
		wc.probability.SetText(fmt.Sprintf("Probability: %.2f%%", wi.Precipitation.Probability*100))
		wc.rain.SetText(fmt.Sprintf("Rain: %v", wi.Precipitation.Rain))
		wc.snow.SetText(fmt.Sprintf("Snow: %v", wi.Precipitation.Snow))
		wc.sleet.SetText(fmt.Sprintf("Sleet: %v", wi.Precipitation.Sleet))
		wc.hail.SetText(fmt.Sprintf("Hail: %v", wi.Precipitation.Hail))
		wc.direction.SetText(fmt.Sprintf("Direction: %s", wi.Wind.Direction))
		wc.angle.SetText(fmt.Sprintf("Angle: %.2f°", wi.Wind.Angle))
		wc.speed.SetText(fmt.Sprintf("Speed: %.2f km/h", wi.Wind.Speed))

		time.Sleep(time.Second)
	}
}
