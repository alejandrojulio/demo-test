package routes

import (
	"backend/db"
	"backend/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetStockHandler(w http.ResponseWriter, r *http.Request) {
	var stockData []models.Stock
	var totalCount int64
	query := db.Db

	search := r.URL.Query().Get("search")
	if search != "" {
		query = query.Where("ticker ILIKE ?", "%"+search+"%").Or("company ILIKE ?", "%"+search+"%").Or("brokerage ILIKE ?", "%"+search+"%")
	}

	limit := r.URL.Query().Get("limit")
	currentPage := r.URL.Query().Get("currentPage")
	if limit != "" && currentPage != "" {
		limitInt, err1 := strconv.Atoi(limit)
		currentPageInt, err2 := strconv.Atoi(currentPage)
		if err1 == nil && err2 == nil && currentPageInt > 0 {
			offset := (currentPageInt - 1) * limitInt
			query = query.Offset(offset).Limit(limitInt)
		} else {
			query = query.Limit(10)
		}
	} else {
		query = query.Limit(10)
	}
	err := db.Db.Model(&models.Stock{}).Where(query).Count(&totalCount).Error
	if err != nil {
		http.Error(w, "Error al contar registros en la DB", http.StatusInternalServerError)
		log.Fatal("Error - GetStockHandler: ", err)
		return
	}

	if err := query.Find(&stockData).Error; err != nil {
		http.Error(w, "Error al consultar la DB", http.StatusInternalServerError)
		log.Fatal("Error - GetStockHandler: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Stocks    []models.Stock `json:"stocks"`
		TotalCount int64         `json:"totalCount"`
	}{
		Stocks:    stockData,
		TotalCount: totalCount,
	}
	json.NewEncoder(w).Encode(response)
}

func GetStockRecomendationHandler(w http.ResponseWriter, r *http.Request) {
	var stockRecomendation []models.StockRecomendation

	query := db.Db

	rawQuery := `
		SELECT DISTINCT ON (ticker)
			id,
			ticker,
			target_from,
			target_to,
			company,
			action,
			brokerage,
			rating_from,
			rating_to,
			time,
			ROUND(((target_to - target_from) / target_from) * 100, 2) AS potential_gain
		FROM stock
		WHERE target_to > target_from
		ORDER BY ticker, ((target_to - target_from) / target_from) DESC
		LIMIT 10;
	`

	err := query.Raw(rawQuery).Scan(&stockRecomendation).Error
	if err != nil {
		http.Error(w, "Error al consultar recomendaciones", http.StatusInternalServerError)
		log.Fatal("Error - GetStockRecomendationHandler: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Stocks []models.StockRecomendation `json:"stocks"`
	}{
		Stocks: stockRecomendation,
	}
	json.NewEncoder(w).Encode(response)
}

