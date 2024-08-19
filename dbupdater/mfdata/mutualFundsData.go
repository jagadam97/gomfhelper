package mfdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Define a struct to map the JSON response
type Meta struct {
	FundHouse  string `json:"fund_house"`
	SchemeName string `json:"scheme_name"`
}

type Data struct {
	Date string `json:"date"`
	Nav  string `json:"nav"`
}

type MFResponse struct {
	Meta Meta   `json:"meta"`
	Data []Data `json:"data"`
}

func GetLatestNav(apiCode string) (date string, nav float64) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.mfapi.in/mf/%s/latest", apiCode)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", 0
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", 0
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", 0
	}

	var mfResponse MFResponse
	if err := json.Unmarshal(body, &mfResponse); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return "", 0
	}

	if len(mfResponse.Data) > 0 {
		date = mfResponse.Data[0].Date
		
		nav, err = strconv.ParseFloat(mfResponse.Data[0].Nav, 64)
		if err != nil {
			fmt.Println("Error converting NAV to float64:", err)
			return "", 0
		}
	}
	return date, nav
}

func GetMFData(apiCode string) ([]Data) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.mfapi.in/mf/%s", apiCode)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var mfResponse MFResponse
	if err := json.Unmarshal(body, &mfResponse); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}

	for i, data := range mfResponse.Data {
		nav, err := strconv.ParseFloat(data.Nav, 64)
		if err != nil {
			fmt.Printf("Error converting NAV to float64 for entry %d: %v\n", i, err)
			continue
		}
		fmt.Printf("Date: %s, NAV: %f\n", data.Date, nav)
	}

	return mfResponse.Data
}
