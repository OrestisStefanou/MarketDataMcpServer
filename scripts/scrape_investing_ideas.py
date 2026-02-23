"""
Scrapes investing ideas from simplywall.st
How to run:
1. Install dependencies:
pip install requests beautifulsoup4
2. Run the script:
python scripts/scrape_investing_ideas.py
"""

import csv
import json
import re
from urllib.parse import (
    urljoin, 
    urlparse,
)

import requests
from bs4 import BeautifulSoup


PAGE_URL = "https://simplywall.st/discover/investing-ideas"
BASE = "https://simplywall.st"
OUTPUT_FILE = "static_data/investing_ideas.json"

def clean_company_name(text):
    """
    Clean company name by removing everything after 'Market Cap' and other metadata
    """
    if not text:
        return None
    
    # Split on common separators and take first part (company name)
    # Remove everything from "Market Cap" onwards
    if "Market Cap" in text:
        text = text.split("Market Cap")[0]
    
    # Remove stock codes (numbers only or exchange codes like SZSE, NYSE, etc.)
    text = re.sub(r'\b\d{6}\b', '', text)  # 6-digit codes
    text = re.sub(r'\b[A-Z]{2,4}:\d+\b', '', text)  # Exchange codes like NYSE:1234
    
    # Remove currency amounts
    text = re.sub(r'[A-Z]{2}\$[\d,\.]+[bmk]?', '', text, flags=re.IGNORECASE)
    text = re.sub(r'[$¥€£₹][\d,\.]+[bmk]?', '', text)
    text = re.sub(r'CN¥[\d,\.]+[bmk]?', '', text)
    text = re.sub(r'NT\$[\d,\.]+[bmk]?', '', text)
    text = re.sub(r'JP¥[\d,\.]+[bmk]?', '', text)
    
    # Remove percentage patterns
    text = re.sub(r'\d+\.\d+%', '', text)
    text = re.sub(r'\d+%', '', text)
    
    # Remove time period indicators
    text = re.sub(r'\d+[DdYy]', '', text)
    text = re.sub(r'[17][YD]', '', text)
    
    # Remove standalone numbers
    text = re.sub(r'\b\d+\b', '', text)
    
    # Remove extra whitespace
    text = re.sub(r'\s+', ' ', text)
    
    # Clean up
    text = text.strip()
    
    # Filter out if it's too short or looks like metadata
    if len(text) < 3:
        return None
    
    # Filter out common non-company words
    exclude_words = ['Market', 'Cap', 'Engages', 'Develops', 'Designs', 'Manufactures', 
                     'New', 'Together', 'Research', 'Stocks']
    if text in exclude_words:
        return None
    
    return text


def scrape_companies(url) -> list[str]:    
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
    }
    
    try:
        response = requests.get(url, headers=headers, timeout=30)
        response.raise_for_status()
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        companies = []
        
        # Method A: Find all links to stock pages and extract company names
        stock_links = soup.find_all('a', href=re.compile(r'/stocks/'))
        
        for link in stock_links:
            # Get the full text from the link
            text = link.get_text(strip=True)
            
            # Clean the company name
            cleaned = clean_company_name(text)
            
            if cleaned:
                companies.append(cleaned)
        
        # Method B: Look for pattern where company name appears before "Market Cap:"
        text_content = soup.get_text()
        
        # More precise pattern: capture company name that ends right before "Market Cap"
        # This pattern looks for text followed immediately by "Market Cap"
        pattern = r'([A-Za-z][A-Za-z\s&\.\'-]+?)(?=Market Cap:)'
        matches = re.findall(pattern, text_content)
        
        for match in matches:
            cleaned = match.strip()
            # Additional cleaning to remove trailing words that aren't part of name
            cleaned = re.sub(r'\s+(Engages|Develops|Designs|Manufactures|Research|Together|Operates).*$', '', cleaned, flags=re.IGNORECASE)
            if len(cleaned) > 2:
                companies.append(cleaned)
        
        # Remove duplicates while preserving order
        seen = set()
        unique_companies = []
        for company in companies:
            if company not in seen and company:
                seen.add(company)
                unique_companies.append(company)
        
        return unique_companies
    
    except Exception as e:
        print(f"Error with requests method: {e}")
        return []


def clean_text(s: str) -> str:
    # Basic whitespace cleanup
    s = re.sub(r"\s+", " ", (s or "")).strip()
    # Remove "+123 companies" and any trailing text like "5 New"
    s = re.sub(r"\+\d+ companies.*", "", s).strip()
    return s


def main():
    headers = {
        "User-Agent": "Mozilla/5.0",
        "Accept": "text/html,application/xhtml+xml",
        "Accept-Language": "en-US,en;q=0.9",
    }

    resp = requests.get(PAGE_URL, headers=headers, timeout=30)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.text, "html.parser")

    items = []
    seen = set()

    # On this page, idea cards are normal links to /discover/investing-ideas/<id>/<slug>/...
    for a in soup.select('a[href^="/discover/investing-ideas/"], a[href^="https://simplywall.st/discover/investing-ideas/"]'):
        href = (a.get("href") or "").strip()
        
        # Try to find the specific title element (p.font-serif)
        title_el = a.select_one("p.font-serif")
        if title_el:
            title = clean_text(title_el.get_text(" ", strip=True))
        else:
            # Fallback to current behavior if not found
            title = clean_text(a.get_text(" ", strip=True))

        if not href or not title:
            continue

        full_url = urljoin(BASE, href)

        if not full_url.endswith("/global"):
            continue

        # Skip the landing page itself
        if full_url.rstrip("/") == PAGE_URL.rstrip("/"):
            continue

        # Keep only same-domain links (safety)
        if urlparse(full_url).netloc != urlparse(BASE).netloc:
            continue

        key = (title, full_url)
        if key in seen:
            continue
        seen.add(key)

        print(f"Scraping companies for {title}")
        companies = scrape_companies(full_url)

        items.append({"title": title, "link": full_url, "companies": companies})

    # Optional: sort for stable output
    items.sort(key=lambda x: x["title"].lower())

    print(f"Found {len(items)} investing ideas.")
    for x in items[:10]:
        print(f"- {x['title']} -> {x['link']}")

    # Save JSON
    with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
        json.dump(items, f, ensure_ascii=False, indent=2)


if __name__ == "__main__":
    main()