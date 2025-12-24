package marketDataScraper

import (
	"encoding/json"
	"fmt"
	"io"
	"market_data_mcp_server/pkg/domain"
	"net/http"
)

func scrapeMarketNews() ([]domain.NewsArticle, error) {
	url := "https://stockanalysis.com/news/__data.json"
	resp, err := http.Get(url)
	if err != nil {
		return []domain.NewsArticle{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []domain.NewsArticle{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return []domain.NewsArticle{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 2 {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure in 'nodes'")
	}

	// Access the second element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[1].(map[string]interface{})
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure in 'data[0]'")
	}

	marketNewsDataIndex, ok := dataMap["data"].(float64)
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure for 'industries'")
	}

	marketNewsDataIndicesArray := data[int(marketNewsDataIndex)].([]interface{})

	marketNews := make([]domain.NewsArticle, 0, len(marketNewsDataIndicesArray))

	for i := 0; i < len(marketNewsDataIndicesArray); i++ {
		marketNewsDataIndex := int(marketNewsDataIndicesArray[i].(float64))
		marketNewsData := data[marketNewsDataIndex].(map[string]interface{})
		urlIndex := int(marketNewsData["url"].(float64))
		imageIndex := int(marketNewsData["img"].(float64))
		titleIndex := int(marketNewsData["title"].(float64))
		textIndex := int(marketNewsData["text"].(float64))
		sourceIndex := int(marketNewsData["source"].(float64))
		timeIndex := int(marketNewsData["time"].(float64))

		article := domain.NewsArticle{
			Url:    data[urlIndex].(string),
			Image:  data[imageIndex].(string),
			Title:  data[titleIndex].(string),
			Text:   data[textIndex].(string),
			Source: data[sourceIndex].(string),
			Time:   data[timeIndex].(string),
		}
		marketNews = append(marketNews, article)
	}

	return marketNews, nil
}

func scrapeStockNews(symbol string) ([]domain.NewsArticle, error) {
	url := fmt.Sprintf("https://stockanalysis.com/stocks/%s/__data.json", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return []domain.NewsArticle{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []domain.NewsArticle{}, err
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return []domain.NewsArticle{}, err
	}

	// Extract "nodes" from rawData
	nodes, ok := rawData["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure in 'nodes'")
	}

	// Access the third element in "nodes" which contains the data we are interested in
	nodeData, ok := nodes[2].(map[string]interface{})
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure in 'nodes[2]'")
	}

	data, ok := nodeData["data"].([]interface{})
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure in 'data'")
	}

	dataMap, ok := data[0].(map[string]interface{})
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure in 'data[0]'")
	}

	marketNewsDataMapIndex, ok := dataMap["news"].(float64)
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure for 'industries'")
	}

	marketNewsDataMap, ok := data[int(marketNewsDataMapIndex)].(map[string]interface{})
	if !ok {
		return []domain.NewsArticle{}, fmt.Errorf("unexpected structure for 'industries'")
	}

	marketNewsDataIndex := marketNewsDataMap["data"].(float64)

	marketNewsDataIndicesArray := data[int(marketNewsDataIndex)].([]interface{})

	marketNews := make([]domain.NewsArticle, 0, len(marketNewsDataIndicesArray))

	for i := 0; i < len(marketNewsDataIndicesArray); i++ {
		marketNewsDataIndex := int(marketNewsDataIndicesArray[i].(float64))
		marketNewsData := data[marketNewsDataIndex].(map[string]interface{})
		urlIndex := int(marketNewsData["url"].(float64))
		imageIndex := int(marketNewsData["img"].(float64))
		titleIndex := int(marketNewsData["title"].(float64))
		textIndex := int(marketNewsData["text"].(float64))
		sourceIndex := int(marketNewsData["source"].(float64))
		timeIndex := int(marketNewsData["time"].(float64))

		article := domain.NewsArticle{
			Url:    data[urlIndex].(string),
			Image:  data[imageIndex].(string),
			Title:  data[titleIndex].(string),
			Text:   data[textIndex].(string),
			Source: data[sourceIndex].(string),
			Time:   data[timeIndex].(string),
		}
		marketNews = append(marketNews, article)
	}

	return marketNews, nil
}
