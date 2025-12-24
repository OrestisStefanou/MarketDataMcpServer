package tools

import (
	"context"
	"fmt"
	"market_data_mcp_server/pkg/domain"

	"github.com/mark3labs/mcp-go/mcp"
)

type EconomicIndicatorsService interface {
	GetRealGdpTimeSeries(interval domain.EconomicIndicatorInterval) (domain.EconomicIndicatorTimeSeries, error)
	GetTreasuryYieldTimeSeries(maturity domain.TreasuryYieldMaturity) (domain.EconomicIndicatorTimeSeries, error)
	GetInterestRatesTimeSeries() (domain.EconomicIndicatorTimeSeries, error)
	GetInflationTimeSeries() (domain.EconomicIndicatorTimeSeries, error)
	GetUnemploymentRateTimeSeries() (domain.EconomicIndicatorTimeSeries, error)
}

type EconomicIndicatorTimeSeriesEntrySchema struct {
	Date  string `json:"date" jsonschema_description:"Date of the data entry"`
	Value string `json:"value" jsonschema_description:"Value of the data entry"`
}

type GetEconomicIndicatorTimeSeriesResponse struct {
	IndicatorName string                                   `json:"indicator_name" jsonschema_description:"Name of the economic indicator"`
	Interval      string                                   `json:"interval" jsonschema_description:"Interval of the economic indicator"`
	Unit          string                                   `json:"unit" jsonschema_description:"Unit of the economic indicator"`
	Data          []EconomicIndicatorTimeSeriesEntrySchema `json:"data" jsonschema_description:"Data entries of the economic indicator"`
}

type GetEconomicIndicatorTimeSeriesRequest struct {
	IndicatorName         string `json:"indicator_name" jsonschema_description:"Name of the economic indicator" jsonschema:"enum=RealGDP,enum=TreasuryYield,enum=InterestRate,enum=Inflation,enum=UnemploymentRate"`
	TreasuryYieldMaturity string `json:"treasury_yield_maturity" jsonschema_description:"Maturity of the treasury yield(Only for TreasuryYield indicator, 5Y is the default)" jsonschema:"enum=3m,enum=2Y,enum=5Y,enum=10Y,enum=30Y"`
	Limit                 int    `json:"limit" jsonschema_description:"Number of data entries to return(Default is 100)"`
}

type GetEconomicIndicatorTimeSeriesTool struct {
	economicIndicatorsService EconomicIndicatorsService
}

func NewGetEconomicIndicatorTimeSeriesTool(economicIndicatorsService EconomicIndicatorsService) (*GetEconomicIndicatorTimeSeriesTool, error) {
	return &GetEconomicIndicatorTimeSeriesTool{
		economicIndicatorsService: economicIndicatorsService,
	}, nil
}

func (t *GetEconomicIndicatorTimeSeriesTool) HandleGetEconomicIndicatorTimeSeries(ctx context.Context, req mcp.CallToolRequest, args GetEconomicIndicatorTimeSeriesRequest) (GetEconomicIndicatorTimeSeriesResponse, error) {
	var err error

	var timeSeries domain.EconomicIndicatorTimeSeries
	switch args.IndicatorName {
	case string(domain.RealGDP):
		timeSeries, err = t.economicIndicatorsService.GetRealGdpTimeSeries(domain.MonthlyEconomicIndicatorInterval)
		if err != nil {
			return GetEconomicIndicatorTimeSeriesResponse{}, err
		}
	case string(domain.TreasuryYield):
		if args.TreasuryYieldMaturity == "" {
			args.TreasuryYieldMaturity = string(domain.FiveYearTreasuryYieldMaturity)
		}
		timeSeries, err = t.economicIndicatorsService.GetTreasuryYieldTimeSeries(domain.TreasuryYieldMaturity(args.TreasuryYieldMaturity))
		if err != nil {
			return GetEconomicIndicatorTimeSeriesResponse{}, err
		}
	case string(domain.InterestRate):
		timeSeries, err = t.economicIndicatorsService.GetInterestRatesTimeSeries()
		if err != nil {
			return GetEconomicIndicatorTimeSeriesResponse{}, err
		}
	case string(domain.Inflation):
		timeSeries, err = t.economicIndicatorsService.GetInflationTimeSeries()
		if err != nil {
			return GetEconomicIndicatorTimeSeriesResponse{}, err
		}
	case string(domain.UnemploymentRate):
		timeSeries, err = t.economicIndicatorsService.GetUnemploymentRateTimeSeries()
		if err != nil {
			return GetEconomicIndicatorTimeSeriesResponse{}, err
		}
	default:
		return GetEconomicIndicatorTimeSeriesResponse{}, fmt.Errorf("invalid economic indicator name: %s", args.IndicatorName)
	}

	response := GetEconomicIndicatorTimeSeriesResponse{
		IndicatorName: string(timeSeries.Name),
		Interval:      string(timeSeries.Interval),
		Unit:          string(timeSeries.Unit),
		Data:          make([]EconomicIndicatorTimeSeriesEntrySchema, 0, len(timeSeries.Data)),
	}

	limit := 100
	if args.Limit > 0 {
		limit = args.Limit
	}

	for _, entry := range timeSeries.Data[:limit] {
		response.Data = append(response.Data, EconomicIndicatorTimeSeriesEntrySchema{
			Date:  entry.Date,
			Value: entry.Value,
		})
	}

	return response, nil
}

func (t *GetEconomicIndicatorTimeSeriesTool) GetTool() mcp.Tool {
	return mcp.NewTool("getEconomicIndicatorTimeSeries",
		mcp.WithDescription("Get the time series of the given economic indicator."),
		mcp.WithInputSchema[GetEconomicIndicatorTimeSeriesRequest](),
		mcp.WithOutputSchema[GetEconomicIndicatorTimeSeriesResponse](),
	)
}
