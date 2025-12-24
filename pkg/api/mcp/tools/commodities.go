package tools

import (
	"context"
	"market_data_mcp_server/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type CommoditiesService interface {
	GetCommodityTimeSeries(commodity domain.Commodity) (domain.CommodityTimeSeries, error)
}

type CommodityTimeSeriesEntrySchema struct {
	Date  string `json:"date" jsonschema_description:"Date of the data entry"`
	Value string `json:"value" jsonschema_description:"Value of the data entry"`
}

type GetCommodityTimeSeriesResponse struct {
	CommodityName string                           `json:"commodity_name" jsonschema_description:"Name of the commodity"`
	Interval      string                           `json:"interval" jsonschema_description:"Interval of the commodity"`
	Unit          string                           `json:"unit" jsonschema_description:"Unit of the commodity"`
	Data          []CommodityTimeSeriesEntrySchema `json:"data" jsonschema_description:"Data entries of the commodity"`
}

type GetCommodityTimeSeriesRequest struct {
	CommodityName string `json:"commodity_name" jsonschema_description:"Name of the commodity" jsonschema:"enum=CrudeOil,enum=NaturalGas,enum=Copper,enum=Aluminum,enum=Wheat,enum=Corn,enum=Sugar,enum=Coffee"`
	Limit         int    `json:"limit" jsonschema_description:"Number of data entries to return(Default is 100)"`
}

type GetCommodityTimeSeriesTool struct {
	commoditiesService CommoditiesService
}

func NewGetCommodityTimeSeriesTool(commoditiesService CommoditiesService) (*GetCommodityTimeSeriesTool, error) {
	return &GetCommodityTimeSeriesTool{
		commoditiesService: commoditiesService,
	}, nil
}

func (t *GetCommodityTimeSeriesTool) HandleGetCommodityTimeSeries(ctx context.Context, req mcp.CallToolRequest, args GetCommodityTimeSeriesRequest) (GetCommodityTimeSeriesResponse, error) {
	var err error

	timeSeries, err := t.commoditiesService.GetCommodityTimeSeries(domain.Commodity(args.CommodityName))
	if err != nil {
		return GetCommodityTimeSeriesResponse{}, err
	}

	response := GetCommodityTimeSeriesResponse{
		CommodityName: string(timeSeries.Name),
		Interval:      string(timeSeries.Interval),
		Unit:          string(timeSeries.Unit),
		Data:          make([]CommodityTimeSeriesEntrySchema, 0, len(timeSeries.Data)),
	}

	limit := 100
	if args.Limit > 0 {
		limit = args.Limit
	}

	for _, entry := range timeSeries.Data[:limit] {
		response.Data = append(response.Data, CommodityTimeSeriesEntrySchema{
			Date:  entry.Date,
			Value: entry.Value,
		})
	}

	return response, nil
}

func (t *GetCommodityTimeSeriesTool) GetTool() mcp.Tool {
	return mcp.NewTool("getCommodityTimeSeries",
		mcp.WithDescription("Get the time series of the given commodity."),
		mcp.WithInputSchema[GetCommodityTimeSeriesRequest](),
		mcp.WithOutputSchema[GetCommodityTimeSeriesResponse](),
	)
}
