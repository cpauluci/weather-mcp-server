package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var apiKey string

const (
	apiUrl                   = "https://api.hgbrasil.com/weather"
	conditionSlugPlaceholder = "{condition_slug}"
	moonPhasePlaceholder     = "{moon_phase}"
	assetsConditionsUrl      = "https://assets.hgbrasil.com/weather/icons/conditions/{condition_slug}.svg"
	assetsMoonUrl            = "https://assets.hgbrasil.com/weather/icons/moon/{moon_phase}.png"
)

type ForecastLatLonInput struct {
	Latitude  float64 `json:"latitude" jsonschema:"Latitude da localização, ex: -23.550196048274127"`
	Longitude float64 `json:"longitude" jsonschema:"Longitude da localização, ex: -46.63396252883539"`
}

type ForecastCityInput struct {
	City string `json:"city" jsonschema:"Nome da cidade, ex: Santo André, SP"`
}

type WeatherResponse struct {
	By            string         `json:"by"`
	ValidKey      bool           `json:"valid_key"`
	Results       WeatherResults `json:"results"`
	ExecutionTime float64        `json:"execution_time"`
	FromCache     bool           `json:"from_cache"`
}

type WeatherResults struct {
	Temp          int           `json:"temp"`
	Date          string        `json:"date"`
	Time          string        `json:"time"`
	ConditionCode string        `json:"condition_code"`
	Description   string        `json:"description"`
	Currently     string        `json:"currently"`
	Woeid         int           `json:"woeid"`
	City          string        `json:"city"`
	ImgID         string        `json:"img_id"`
	Humidity      int           `json:"humidity"`
	Cloudiness    float64       `json:"cloudiness"`
	Rain          float64       `json:"rain"`
	WindSpeedy    string        `json:"wind_speedy"`
	WindDirection int           `json:"wind_direction"`
	WindCardinal  string        `json:"wind_cardinal"`
	Sunrise       string        `json:"sunrise"`
	Sunset        string        `json:"sunset"`
	MoonPhase     string        `json:"moon_phase"`
	ConditionSlug string        `json:"condition_slug"`
	CityName      string        `json:"city_name"`
	Timezone      string        `json:"timezone"`
	Forecast      []ForecastDay `json:"forecast"`
	Cref          string        `json:"cref"`
}

type ForecastDay struct {
	Date            string  `json:"date"`
	FullDate        string  `json:"full_date"`
	Weekday         string  `json:"weekday"`
	Max             int     `json:"max"`
	Min             int     `json:"min"`
	Humidity        int     `json:"humidity"`
	Cloudiness      float64 `json:"cloudiness"`
	Rain            float64 `json:"rain"`
	RainProbability int     `json:"rain_probability"`
	WindSpeedy      string  `json:"wind_speedy"`
	Sunrise         string  `json:"sunrise"`
	Sunset          string  `json:"sunset"`
	MoonPhase       string  `json:"moon_phase"`
	Description     string  `json:"description"`
	Condition       string  `json:"condition"`
}

func init() {
	// precedence: flags > environment variables
	flag.StringVar(&apiKey, "API_KEY", "", "API key for HgBrasilWeather API")
	flag.Parse()

	if apiKey == "" {
		apiKey, ok := os.LookupEnv("API_KEY")
		if !ok {
			log.Fatalf("API_KEY not found in environment variables")
		}
		if apiKey == "" {
			log.Fatalf("API_KEY is empty")
		}
	}

}

func makeHgBrasilRequest[T any](ctx context.Context, url string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := http.DefaultClient
	log.Println("Making request to", url)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	// to debug remove the comment below
	// log.Println("Response: ", string(body))
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func formatForecast(result *WeatherResults) string {
	forecastStr := ""
	for _, forecast := range result.Forecast {
		logoCondicao := strings.Replace(assetsConditionsUrl, conditionSlugPlaceholder, forecast.Condition, 1)
		logoFaseLua := strings.Replace(assetsMoonUrl, moonPhasePlaceholder, forecast.MoonPhase, 1)
		forecastStr += fmt.Sprintf(`
	Previsão para os próximos dias:
	Data: %s
	Temperatura Max: %d°C | Min: %d°C
	Descrição: %s
	Logo da Condição: %s
	Humidade: %d%%
	Nebulosidade: %.2f%%
	Chuva: %.2f%%
	Probabilidade de Chuva: %d%%
	Logo da Fase da Lua: %s
	`, forecast.Date, forecast.Max, forecast.Min, forecast.Description, logoCondicao, forecast.Humidity,
			forecast.Cloudiness, forecast.Rain, forecast.RainProbability, logoFaseLua)
	}
	logoCondicao := strings.Replace(assetsConditionsUrl, conditionSlugPlaceholder, result.ConditionSlug, 1)
	logoFaseLua := strings.Replace(assetsMoonUrl, moonPhasePlaceholder, result.MoonPhase, 1)
	dadosAtuais := fmt.Sprintf(`
	Data Atual: %s
	Cidade: %s
	Temperatura Atual: %d°C
	Velocidade do Vento: %s
	Descrição da Condição Atual: %s
	Logo da Condição Atual: %s
	Humidade: %d%%
	Nebulosidade: %.2f%%
	Chuva: %.2f%%
	Logo da Fase da Lua: %s
	`, result.Date, result.City, result.Temp, result.WindSpeedy, result.Description,
		logoCondicao, result.Humidity, result.Cloudiness, result.Rain, logoFaseLua)
	return dadosAtuais + forecastStr
}

func getForecastLatLon(ctx context.Context, req *mcp.CallToolRequest, input ForecastLatLonInput) (*mcp.CallToolResult, any, error) {
	log.Println("MCP get_forecast_lat_lon:", req)

	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%g", input.Latitude))
	params.Add("lon", fmt.Sprintf("%g", input.Longitude))
	params.Add("key", apiKey)
	forecastUrl := fmt.Sprintf("%s?%s", apiUrl, params.Encode())

	data, err := makeHgBrasilRequest[WeatherResponse](ctx, forecastUrl)
	if err != nil {
		log.Println("Error getting forecast:", err)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Não foi possível obter dados de previsão para a localização informada."},
			},
		}, nil, nil
	}

	result := formatForecast(&data.Results)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func getForecastCity(ctx context.Context, req *mcp.CallToolRequest, input ForecastCityInput) (*mcp.CallToolResult, any, error) {
	log.Println("MCP get_forecast_city:", req)

	params := url.Values{}
	params.Add("city_name", input.City)
	params.Add("key", apiKey)
	forecastUrl := fmt.Sprintf("%s?%s", apiUrl, params.Encode())

	data, err := makeHgBrasilRequest[WeatherResponse](ctx, forecastUrl)
	if err != nil {
		log.Println("Error getting forecast:", err)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Não foi possível obter dados de previsão para a cidade informada."},
			},
		}, nil, nil
	}

	result := formatForecast(&data.Results)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func main() {
	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "weather",
		Version: "1.0.0",
	}, nil)

	// Add get_forecast_lat_lon tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_forecast_lat_lon",
		Description: "Retorna a previsão do tempo para uma determinada latitude e longitude",
	}, getForecastLatLon)

	// Add get_forecast_city tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_forecast_city",
		Description: "Retorna a previsão do tempo para uma determinada cidade.",
	}, getForecastCity)

	// Run server on stdio transport
	log.Println("MCP server running on stdio")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Error running MCP server: %v", err)
	}
}
