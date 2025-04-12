package routes_test

import (
	"backend/db"
	"backend/models"
	"backend/routes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStockHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/stock?search=BSBR&limit=10&currentPage=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	db.Db = db.InitMockDb()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.GetStockHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response struct {
		Stocks     []models.Stock `json:"stocks"`
		TotalCount int64          `json:"totalCount"`
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Stocks)
	assert.Greater(t, response.TotalCount, int64(0))
}

func TestGetStockRecomendationHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/stock-recomendation", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(routes.GetStockRecomendationHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var response struct {
		Stocks []models.StockRecomendation `json:"stocks"`
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Stocks)
}
