package fred

import (
	"fmt"

	"market_data_mcp_server/pkg/domain"
)

// FredClient reads public time series from the St. Louis Fed's FRED database.
// FRED's CSV endpoint needs no API key.
type FredClient struct{}

func NewFredClient() (*FredClient, error) {
	return &FredClient{}, nil
}

func (c *FredClient) GetRealGdpTimeSeries(interval domain.EconomicIndicatorInterval) (domain.EconomicIndicatorTimeSeries, error) {
	spec := realGdpQuarterlySpec
	if interval == domain.AnnualEconomicIndicatorInterval {
		spec = realGdpAnnualSpec
	}
	return c.economicIndicator(domain.RealGDP, spec)
}

func (c *FredClient) GetTreasuryYieldTimeSeries(maturity domain.TreasuryYieldMaturity) (domain.EconomicIndicatorTimeSeries, error) {
	spec, ok := treasuryYieldSpecs[maturity]
	if !ok {
		return domain.EconomicIndicatorTimeSeries{}, fmt.Errorf("unsupported treasury yield maturity: %s", maturity)
	}
	return c.economicIndicator(domain.TreasuryYield, spec)
}

func (c *FredClient) GetInterestRatesTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	return c.economicIndicator(domain.InterestRate, interestRateSpec)
}

func (c *FredClient) GetInflationTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	return c.economicIndicator(domain.Inflation, inflationSpec)
}

func (c *FredClient) GetUnemploymentRateTimeSeries() (domain.EconomicIndicatorTimeSeries, error) {
	return c.economicIndicator(domain.UnemploymentRate, unemploymentRateSpec)
}

func (c *FredClient) GetCommodityTimeSeries(commodity domain.Commodity) (domain.CommodityTimeSeries, error) {
	spec, ok := commoditySpecs[commodity]
	if !ok {
		return domain.CommodityTimeSeries{}, fmt.Errorf("unsupported commodity: %s", commodity)
	}

	observations, err := fetchSeries(spec.Id, spec.LookbackYears)
	if err != nil {
		return domain.CommodityTimeSeries{}, err
	}

	observations = reverse(observations)
	data := make([]domain.CommodityTimeSeriesEntry, 0, len(observations))
	for _, obs := range observations {
		data = append(data, domain.CommodityTimeSeriesEntry{Date: obs.Date, Value: obs.Value})
	}

	return domain.CommodityTimeSeries{
		Name:     commodity,
		Interval: spec.Interval,
		Unit:     spec.Unit,
		Data:     data,
	}, nil
}

func (c *FredClient) economicIndicator(name domain.EconomicIndicator, spec economicIndicatorSpec) (domain.EconomicIndicatorTimeSeries, error) {
	observations, err := fetchSeries(spec.Id, spec.LookbackYears)
	if err != nil {
		return domain.EconomicIndicatorTimeSeries{}, err
	}

	// The conversion runs while the series is still ascending, since it reads each
	// observation's prior-year neighbour.
	if spec.YearOverYear {
		observations, err = toYearOverYearPercent(observations)
		if err != nil {
			return domain.EconomicIndicatorTimeSeries{}, err
		}
	}

	observations = reverse(observations)
	data := make([]domain.EconomicIndicatorTimeSeriesEntry, 0, len(observations))
	for _, obs := range observations {
		data = append(data, domain.EconomicIndicatorTimeSeriesEntry{Date: obs.Date, Value: obs.Value})
	}

	return domain.EconomicIndicatorTimeSeries{
		Name:     name,
		Interval: spec.Interval,
		Unit:     spec.Unit,
		Data:     data,
	}, nil
}
