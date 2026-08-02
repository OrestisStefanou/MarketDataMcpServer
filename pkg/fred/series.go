package fred

import "market_data_mcp_server/pkg/domain"

// economicIndicatorSpec describes how one FRED series maps onto our domain model.
//
// lookbackYears bounds the CSV payload via FRED's cosd parameter. Daily series are the
// reason it exists: the full history of DGS10 is roughly 17,000 rows.
type economicIndicatorSpec struct {
	Id            string
	Interval      domain.EconomicIndicatorInterval
	Unit          domain.EconomicIndicatorUnit
	LookbackYears int
	// YearOverYear converts an index-level series into percent change against the same
	// period a year earlier. FRED publishes CPI as an index, but the tool has always
	// returned an inflation rate.
	YearOverYear bool
}

// realGdpQuarterlySpec is used for every interval except Annual. The tool currently always
// asks for Monthly, and there is no monthly real GDP series, so quarterly is the closest
// truthful answer. The response reports the interval actually returned, not the one asked for.
var realGdpQuarterlySpec = economicIndicatorSpec{
	Id:       "GDPC1",
	Interval: domain.QuarterlyEconomicIndicatorInterval,
	Unit:     domain.BillionsOfDollarsEconomicIndicatorUnit,
}

var realGdpAnnualSpec = economicIndicatorSpec{
	Id:       "GDPCA",
	Interval: domain.AnnualEconomicIndicatorInterval,
	Unit:     domain.BillionsOfDollarsEconomicIndicatorUnit,
}

var interestRateSpec = economicIndicatorSpec{
	Id:            "FEDFUNDS",
	Interval:      domain.MonthlyEconomicIndicatorInterval,
	Unit:          domain.PercentEconomicIndicatorUnit,
	LookbackYears: 40,
}

// CPIAUCSL is an index (1982-84 = 100), converted to a year over year rate below. The
// lookback is one year longer than the window we want, since the first 12 months of an
// index series produce no year over year value.
var inflationSpec = economicIndicatorSpec{
	Id:            "CPIAUCSL",
	Interval:      domain.MonthlyEconomicIndicatorInterval,
	Unit:          domain.PercentEconomicIndicatorUnit,
	LookbackYears: 41,
	YearOverYear:  true,
}

var unemploymentRateSpec = economicIndicatorSpec{
	Id:            "UNRATE",
	Interval:      domain.MonthlyEconomicIndicatorInterval,
	Unit:          domain.PercentEconomicIndicatorUnit,
	LookbackYears: 40,
}

var treasuryYieldSpecs = map[domain.TreasuryYieldMaturity]economicIndicatorSpec{
	domain.ThreeMonthTreasuryYieldMaturity: {Id: "DGS3MO", Interval: domain.DailyEconomicIndicatorInterval, Unit: domain.PercentEconomicIndicatorUnit, LookbackYears: 10},
	domain.TwoYearTreasuryYieldMaturity:    {Id: "DGS2", Interval: domain.DailyEconomicIndicatorInterval, Unit: domain.PercentEconomicIndicatorUnit, LookbackYears: 10},
	domain.FiveYearTreasuryYieldMaturity:   {Id: "DGS5", Interval: domain.DailyEconomicIndicatorInterval, Unit: domain.PercentEconomicIndicatorUnit, LookbackYears: 10},
	domain.TenYearTreasuryYieldMaturity:    {Id: "DGS10", Interval: domain.DailyEconomicIndicatorInterval, Unit: domain.PercentEconomicIndicatorUnit, LookbackYears: 10},
	domain.ThirtyYearTreasuryYieldMaturity: {Id: "DGS30", Interval: domain.DailyEconomicIndicatorInterval, Unit: domain.PercentEconomicIndicatorUnit, LookbackYears: 10},
}

// commoditySpec is the commodity equivalent of economicIndicatorSpec. Commodities have
// their own interval and unit enums in the domain package, hence the separate type.
type commoditySpec struct {
	Id            string
	Interval      domain.CommodityInterval
	Unit          domain.CommodityUnit
	LookbackYears int
}

// Oil and gas use FRED's daily spot series. The rest use the IMF global price series,
// which are monthly and publish with a one to two month lag.
var commoditySpecs = map[domain.Commodity]commoditySpec{
	domain.CrudeOil:   {Id: "DCOILWTICO", Interval: domain.DailyCommodityInterval, Unit: domain.DollarsPerBarrelCommodityUnit, LookbackYears: 10},
	domain.NaturalGas: {Id: "DHHNGSP", Interval: domain.DailyCommodityInterval, Unit: domain.DollarsPerMillionBTUCommodityUnit, LookbackYears: 10},
	domain.Copper:     {Id: "PCOPPUSDM", Interval: domain.MonthlyCommodityInterval, Unit: domain.DollarsPerMetricTonCommodityUnit, LookbackYears: 40},
	domain.Aluminum:   {Id: "PALUMUSDM", Interval: domain.MonthlyCommodityInterval, Unit: domain.DollarsPerMetricTonCommodityUnit, LookbackYears: 40},
	domain.Wheat:      {Id: "PWHEAMTUSDM", Interval: domain.MonthlyCommodityInterval, Unit: domain.DollarsPerMetricTonCommodityUnit, LookbackYears: 40},
	domain.Corn:       {Id: "PMAIZMTUSDM", Interval: domain.MonthlyCommodityInterval, Unit: domain.DollarsPerMetricTonCommodityUnit, LookbackYears: 40},
	// PSUGAISAUSDM is the International Sugar Agreement world price. PSUGAUSAUSDM is the
	// US price and runs roughly 2.5x higher; they are easy to confuse.
	domain.Sugar: {Id: "PSUGAISAUSDM", Interval: domain.MonthlyCommodityInterval, Unit: domain.CentsPerPoundCommodityUnit, LookbackYears: 40},
	// Other Mild Arabica. PCOFFROBUSDM is Robusta.
	domain.Coffee: {Id: "PCOFFOTMUSDM", Interval: domain.MonthlyCommodityInterval, Unit: domain.CentsPerPoundCommodityUnit, LookbackYears: 40},
}
