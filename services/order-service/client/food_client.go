package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/dto"
	"os"

	"github.com/google/uuid"
)

type FoodClient interface {
	GetFoodById(id uuid.UUID) (*dto.FoodResponse, error)
}

type FoodClientImpl struct {
	baseUrl    string
	httpClient *http.Client
}

func NewFoodClientImpl() FoodClient {
	baseUrl := os.Getenv("FOOD_SERVICE_URL")
	return &FoodClientImpl{
		baseUrl:    baseUrl,
		httpClient: &http.Client{},
	}
}

func (c *FoodClientImpl) GetFoodById(id uuid.UUID) (*dto.FoodResponse, error) {
	url := fmt.Sprintf("%s/food/%s", c.baseUrl, id.String())

	response, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get food by id: %s", response.Status)
	}

	var food dto.FoodResponse
	if err := json.NewDecoder(response.Body).Decode(&food); err != nil {
		return nil, err
	}

	return &food, nil
}
