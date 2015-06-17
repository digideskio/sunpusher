package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Config struct {
	WeatherService  string
	PushbulletTitle string
	PushbulletUrl   string
	PushbulletToken string
}

type WeatherResponse struct {
	Hourly Hourly `json:"hourly"`
}

type Hourly struct {
	Summary string `json:"summary"`
}

func ReadConfig() (Config, error) {
	var config Config

	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	config, err := ReadConfig()
	if err != nil {
		panic(err.Error())
	}

	var weatherResponse WeatherResponse
	weatherRes, err := http.Get(config.WeatherService)
	if err != nil {
		panic(err.Error())
	}

	defer weatherRes.Body.Close()

	decoder := json.NewDecoder(weatherRes.Body)

	err = decoder.Decode(&weatherResponse)

	if err != nil {
		panic(err.Error())
	}

	client := &http.Client{}
	data := url.Values{}
	data.Set("type", "note")
	data.Add("title", config.PushbulletTitle)
	data.Add("body", weatherResponse.Hourly.Summary)
	body := bytes.NewBufferString(data.Encode())
	req, err := http.NewRequest("POST", config.PushbulletUrl, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(config.PushbulletToken, "")
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()
}
