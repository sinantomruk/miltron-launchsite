package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sinantomruk/miltron-launchsite/models"
)

func parseResponse(res *http.Response, v interface{}) error {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func doReq(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 205 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}
	return res, nil
}

func (a *App) doReq(method, url string, v interface{}) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", a.apiKey)
	for i := 0; i < 3; i++ {
		res, err := doReq(req)
		if err != nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if err = parseResponse(res, v); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("connection error, %v", req)
}

func (a *App) Rockets() ([]*models.RocketInfo, error) {
	url := fmt.Sprintf("%s/rockets", a.baseUrl)
	rockets := []*models.RocketInfo{}
	if err := a.doReq(http.MethodGet, url, &rockets); err != nil {
		return nil, err
	}
	return rockets, nil
}

func (a *App) LaunchRocket(rocketID string) (*models.RocketInfo, error) {
	url := fmt.Sprintf("%s/rocket/%s/status/launched", a.baseUrl, rocketID)
	rocket := &models.RocketInfo{}
	if err := a.doReq(http.MethodPut, url, rocket); err != nil {
		return nil, err
	}
	return rocket, nil
}

func (a *App) DeployRocket(rocketID string) (*models.RocketInfo, error) {
	url := fmt.Sprintf("%s/rocket/%s/status/deployed", a.baseUrl, rocketID)
	rocket := &models.RocketInfo{}
	if err := a.doReq(http.MethodPut, url, rocket); err != nil {
		return nil, err
	}
	return rocket, nil
}

func (a *App) CancelLaunch(rocketID string) (*models.RocketInfo, error) {
	url := fmt.Sprintf("%s/rocket/%s/status/launched", a.baseUrl, rocketID)
	rocket := &models.RocketInfo{}
	if err := a.doReq(http.MethodDelete, url, rocket); err != nil {
		return nil, err
	}
	return rocket, nil
}

func (a *App) Weather() (*models.Weather, error) {
	url := fmt.Sprintf("%s/weather", a.baseUrl)
	weather := &models.Weather{}
	if err := a.doReq(http.MethodGet, url, weather); err != nil {
		return nil, err
	}
	return weather, nil
}
