package domain

type Commodity string

const (
	CrudeOil   Commodity = "CrudeOil"
	NaturalGas Commodity = "NaturalGas"
	Copper     Commodity = "Copper"
	Aluminum   Commodity = "Aluminum"
	Wheat      Commodity = "Wheat"
	Corn       Commodity = "Corn"
	Sugar      Commodity = "Sugar"
	Coffee     Commodity = "Coffee"
)

type CommodityInterval string

const (
	DailyCommodityInterval     CommodityInterval = "Daily"
	WeeklyCommodityInterval    CommodityInterval = "Weekly"
	MonthlyCommodityInterval   CommodityInterval = "Monthly"
	QuarterlyCommodityInterval CommodityInterval = "Quarterly"
	AnnualCommodityInterval    CommodityInterval = "Annual"
)

type CommodityUnit string

const (
	DollarsPerBarrelCommodityUnit     CommodityUnit = "DollarsPerBarrel"
	DollarsPerMillionBTUCommodityUnit CommodityUnit = "DollarsPerMillionBTU"
	DollarsPerMetricTonCommodityUnit  CommodityUnit = "DollarsPerMetricTon"
	CentsPerPoundCommodityUnit        CommodityUnit = "CentsPerPound"
)

type CommodityTimeSeriesEntry struct {
	Date  string
	Value string
}

type CommodityTimeSeries struct {
	Name     Commodity
	Interval CommodityInterval
	Unit     CommodityUnit
	Data     []CommodityTimeSeriesEntry
}
