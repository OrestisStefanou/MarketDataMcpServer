package marketDataScraper

import (
	"fmt"
	"market_data_mcp_server/pkg/domain"
	"market_data_mcp_server/pkg/errors"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// scrapeSuperInvestorsAndPortfolioLinks returns a map with key the super investor name
// and value the link for the super investor portfolio
func scrapeSuperInvestorsAndPortfolioLinks() (map[string]string, error) {
	url := "https://www.dataroma.com/m/managers.php"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Call to https://www.dataroma.com/m/managers.php failed with status code: %d", rsp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	investorToPortfolioLinkMap := make(map[string]string)
	// Target the specific table containing super investors
	doc.Find("table#grid a[href^='/m/holdings.php']").Each(func(index int, item *goquery.Selection) {
		name := item.Text()
		link, exists := item.Attr("href")
		if exists && len(name) > 0 {
			// fmt.Printf("Investor: %s, Link: %s\n", name, "https://www.dataroma.com"+link)
			investorToPortfolioLinkMap[name] = "https://www.dataroma.com" + link
		}
	})

	return investorToPortfolioLinkMap, nil
}

func scrapeSuperInvestors() ([]domain.SuperInvestor, error) {
	investorToPortfolioLinkMap, err := scrapeSuperInvestorsAndPortfolioLinks()
	if err != nil {
		return nil, err
	}

	superInvestors := make([]domain.SuperInvestor, 0, len(investorToPortfolioLinkMap))
	for superInvestorName, _ := range investorToPortfolioLinkMap {
		superInvestors = append(superInvestors, domain.SuperInvestor{Name: superInvestorName})
	}

	return superInvestors, nil
}

func scrapeSuperInvestorPortfolio(superInvestorName string) (domain.SuperInvestorPortfolio, error) {
	investorToPortfolioLinkMap, err := scrapeSuperInvestorsAndPortfolioLinks()
	if err != nil {
		return domain.SuperInvestorPortfolio{}, err
	}

	portfolioLink, found := investorToPortfolioLinkMap[superInvestorName]
	if !found {
		return domain.SuperInvestorPortfolio{}, &errors.SuperInvestorPortfolioNotFoundError{Message: fmt.Sprintf("Portfolio for super investor: %s not found", superInvestorName)}
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", portfolioLink, nil)
	if err != nil {
		return domain.SuperInvestorPortfolio{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return domain.SuperInvestorPortfolio{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return domain.SuperInvestorPortfolio{}, fmt.Errorf("Call to %s failed with status code: %d", portfolioLink, resp.StatusCode)
	}

	// Parse HTML response with goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return domain.SuperInvestorPortfolio{}, err
	}

	var holdings []domain.SuperInvestorPortfolioHolding
	var sectorAnalysis []domain.SuperInvestorPortfolioSectorAnalysis

	// Process each row
	doc.Find("table tr").Each(func(rowIndex int, rowHtml *goquery.Selection) {
		cells := rowHtml.Find("td")
		if cells.Length() == 0 {
			return // Skip empty rows
		}

		if cells.Length() >= 7 { // Ensure it's a stock row (not a sector row)
			stock := strings.TrimSpace(cells.Eq(1).Text())          // Stock name
			portfolioPct := strings.TrimSpace(cells.Eq(2).Text())   // % of portfolio
			recentActivity := strings.TrimSpace(cells.Eq(3).Text()) // Recent activity
			shares := strings.TrimSpace(cells.Eq(4).Text())         // Shares held
			value := strings.TrimSpace(cells.Eq(6).Text())          // Value of holding

			// Ensure it's not a header row
			if stock != "Stock" {
				holdings = append(holdings, domain.SuperInvestorPortfolioHolding{
					Stock:          stock,
					PortfolioPct:   portfolioPct,
					RecentActivity: recentActivity,
					Shares:         shares,
					Value:          value,
				})
			}
		} else if cells.Length() == 3 { // Likely a sector row
			sector := strings.TrimSpace(cells.Eq(0).Text())
			portfolioPct := strings.TrimSpace(cells.Eq(1).Text())

			// Ensure it's not a header row
			if sector != "" && portfolioPct != "" {
				sectorAnalysis = append(sectorAnalysis, domain.SuperInvestorPortfolioSectorAnalysis{
					Sector:       sector,
					PortfolioPct: portfolioPct,
				})
			}
		}
	})

	return domain.SuperInvestorPortfolio{
		Holdings:       holdings,
		SectorAnalysis: sectorAnalysis,
	}, nil
}
