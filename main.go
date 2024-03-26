package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const openWeatherAPIKey = "5aa3ea249062e431351b6efadc64157d"

type WeatherCondition struct {
	Description string `json:"description"`
}

type MainWeather struct {
	Temp float64 `json:"temp"`
}

type WeatherResponse struct {
	Weather []WeatherCondition `json:"weather"`
	Main    MainWeather        `json:"main"`
}

func main() {
	fmt.Println("Jack Henry Weather Server")
	http.HandleFunc("/weather", weatherHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	weatherResp, err := getWeather(lat, lon)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		log.Println("Failed to fetch weather data:", err)
		return
	}

	if len(weatherResp.Weather) == 0 {
		http.Error(w, "Weather data not available", http.StatusInternalServerError)
		log.Println("Weather data not available")
		return
	}

	weatherCondition := weatherResp.Weather[0].Description

	var tempCondition string
	switch {
	case weatherResp.Main.Temp < 55:
		tempCondition = "cold"
	case weatherResp.Main.Temp > 90:
		tempCondition = "hot"
	default:
		tempCondition = "moderate"
	}

	response := fmt.Sprintf("Weather: %s\nTemperature: %s", weatherCondition, tempCondition)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, response)
}

func getWeather(lat, lon float64) (*WeatherResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/3.0/weather?lat=%f&lon=%f&appid=%s&units=imperial", lat, lon, openWeatherAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weatherResp WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return nil, err
	}

	return &weatherResp, nil
}

/*
JSON Response from OpenWeathermap.org:

{
  "coord": {
    "lon": 10.99,
    "lat": 44.34
  },
  "weather": [
    {
      "id": 501,
      "main": "Rain",
      "description": "moderate rain",
      "icon": "10d"
    }
  ],
  "base": "stations",
  "main": {
    "temp": 298.48,
    "feels_like": 298.74,
    "temp_min": 297.56,
    "temp_max": 300.05,
    "pressure": 1015,
    "humidity": 64,
    "sea_level": 1015,
    "grnd_level": 933
  },
  "visibility": 10000,
  "wind": {
    "speed": 0.62,
    "deg": 349,
    "gust": 1.18
  },
  "rain": {
    "1h": 3.16
  },
  "clouds": {
    "all": 100
  },
  "dt": 1661870592,
  "sys": {
    "type": 2,
    "id": 2075663,
    "country": "IT",
    "sunrise": 1661834187,
    "sunset": 1661882248
  },
  "timezone": 7200,
  "id": 3163858,
  "name": "Zocca",
  "cod": 200
}

*/
