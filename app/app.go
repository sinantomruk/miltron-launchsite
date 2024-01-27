package app

type App struct {
	baseUrl string
	apiKey  string
}

func NewApp(baseUrl, apiKey string) *App {
	return &App{
		baseUrl: baseUrl,
		apiKey:  apiKey,
	}
}
