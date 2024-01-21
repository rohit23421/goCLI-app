package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

//converting the json response to struct
// for this we create a struct to store thetype of data
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	fmt.Println("CLI from Golangs")
	//setting the default location value to Delhi , until the user passed his own
	q := "Mumbai"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}
	//Hitting the Get route for the weather api
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=687b8fded3914cb19fc33111242101&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}

	//closing the get request after everthing is fetched
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// full json response from the api
	//fmt.Println(string(body))

	//crearing a varibale of weather struct type
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}
	//printing thewhole struct format that we stored
	//fmt.Println(weather)

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s: %.0f C, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {

		//creating a date first for getting the timeepoch using the time package
		date := time.Unix(hour.TimeEpoch, 0)

		//checking if hours printed in cli is for future only means greater than current time
		if date.Before(time.Now()) {
			continue
		} // we skip this if time is before the current time and only get the future values
		message := fmt.Sprintf("%s - %.0f C, %.0f, %s\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)

		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}

	}

}
