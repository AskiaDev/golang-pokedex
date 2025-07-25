package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


func (c *Client) ListLocations(pageURL *string) (MapResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL 
	}

	if cachedData, found := c.cache.Get(url); found {
		locationResp := MapResponse{}
		err := json.Unmarshal(cachedData, &locationResp)
		if err == nil {
			return locationResp, nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MapResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return MapResponse{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return MapResponse{}, err
	}

	c.cache.Add(url, data)

	locationResp := MapResponse{}
	err = json.Unmarshal(data, &locationResp)
	if err != nil {
		return MapResponse{}, err
	}

	return locationResp, nil
} 

func (c *Client) GetLocationByName(locationName string) (*MapResult, error) {
	var pageURL *string

	for {
		locationResp, err := c.ListLocations(pageURL)

		if err != nil {
			return nil, err
		}

		for _, location := range locationResp.Results {
			if location.Name == locationName {
				return &location, nil
			}
		}

		if locationResp.Next == "" {
			return nil, fmt.Errorf("location %s not found", locationName)
		}

		pageURL = &locationResp.Next
	}
}

func (c *Client) GetAreaDetails(areaURL string) (AreaResponse, error) {
	if cachedData, found := c.cache.Get(areaURL); found {
		areaResp := AreaResponse{}
		err := json.Unmarshal(cachedData, &areaResp)
		if err == nil {
			return areaResp, nil
		}
	}
	req, err := http.NewRequest("GET", areaURL, nil)
	if err != nil {
		return AreaResponse{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return AreaResponse{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return AreaResponse{}, err
	}

	c.cache.Add(areaURL, data)

	areaResp := AreaResponse{}
	err = json.Unmarshal(data, &areaResp)
	if err != nil {
		return AreaResponse{}, err
	}

	return areaResp, nil
}