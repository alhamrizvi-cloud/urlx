package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const banner = `

██  ██ █████▄  ██     ██  ██ 
██  ██ ██▄▄██▄ ██      ████  
▀████▀ ██   ██ ██████ ██  ██ 
                             
    Multi-Source URL Discovery & Advanced Web Crawler
    ==================================================
`

type ProbeResult struct {
	URL            string
	StatusCode     int
	ContentLength  int64
	Title          string
	Server         string
	ContentType    string
	Technologies   []string
	TLSVersion     string
	Certificate    string
	ResponseTime   time.Duration
	IsAlive        bool
	DNSResolved    bool
	TCPOpen        bool
	TLSHandshake   bool
	Error          string
}

type CrawlResult struct {
	URL           string
	SourceURL     string
	ExtractedURLs []string
	JSFiles       []string
	CSSFiles      []string
	Images        []string
	Forms         []string
	Links         []string
	Scripts       []string
	Depth         int
}

type URLSource interface {
	Fetch(domain string) ([]string, error)
	Name() string
}

// Wayback Machine
type WaybackMachine struct {
	client *http.Client
}

func (w *WaybackMachine) Name() string { return "Wayback Machine" }

func (w *WaybackMachine) Fetch(domain string) ([]string, error) {
	apiURL := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=json&collapse=urlkey", domain)
	resp, err := w.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var results [][]string
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for i, record := range results {
		if i == 0 {
			continue
		}
		if len(record) > 2 {
			urls = append(urls, record[2])
		}
	}
	return urls, nil
}

// Archive.today
type ArchiveToday struct {
	client *http.Client
}

func (a *ArchiveToday) Name() string { return "Archive.today" }

func (a *ArchiveToday) Fetch(domain string) ([]string, error) {
	apiURL := fmt.Sprintf("https://archive.ph/newest/%s", domain)
	resp, err := a.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	content := string(body)
	urls := make([]string, 0)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "http") && strings.Contains(line, domain) {
			if u := extractURL(line); u != "" && strings.Contains(u, domain) {
				urls = append(urls, u)
			}
		}
	}
	return urls, nil
}

// Common Crawl
type CommonCrawl struct {
	client *http.Client
}

func (c *CommonCrawl) Name() string { return "Common Crawl" }

func (c *CommonCrawl) Fetch(domain string) ([]string, error) {
	apiURL := fmt.Sprintf("http://index.commoncrawl.org/CC-MAIN-2024-10-index?url=*.%s&output=json", domain)
	resp, err := c.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	urls := make([]string, 0)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var record map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			continue
		}
		if urlStr, ok := record["url"].(string); ok {
			urls = append(urls, urlStr)
		}
	}
	return urls, scanner.Err()
}

// URLScan.io
type URLScan struct {
	client *http.Client
	apiKey string
}

func (u *URLScan) Name() string { return "URLScan.io" }

func (u *URLScan) Fetch(domain string) ([]string, error) {
	apiURL := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	if u.apiKey != "" {
		req.Header.Set("API-Key", u.apiKey)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Page struct {
				URL string `json:"url"`
			} `json:"page"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, r := range result.Results {
		urls = append(urls, r.Page.URL)
	}
	return urls, nil
}

// AlienVault OTX
type AlienVaultOTX struct {
	client *http.Client
	apiKey string
}

func (a *AlienVaultOTX) Name() string { return "AlienVault OTX" }

func (a *AlienVaultOTX) Fetch(domain string) ([]string, error) {
	apiURL := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/url_list", domain)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	if a.apiKey != "" {
		req.Header.Set("X-OTX-API-KEY", a.apiKey)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		URLList []struct {
			URL string `json:"url"`
		} `json:"url_list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, u := range result.URLList {
		urls = append(urls, u.URL)
	}
	return urls, nil
}

// crt.sh
type CrtSh struct {
	client *http.Client
}

func (c *CrtSh) Name() string { return "crt.sh" }

func (c *CrtSh) Fetch(domain string) ([]string, error) {
	apiURL := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := c.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []struct {
		NameValue string `json:"name_value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	seen := make(map[string]bool)
	for _, r := range results {
		subdomains := strings.Split(r.NameValue, "\n")
		for _, sub := range subdomains {
			sub = strings.TrimSpace(sub)
			if sub != "" && !seen[sub] {
				urls = append(urls, "https://"+sub)
				seen[sub] = true
			}
		}
	}
	return urls, nil
}

// VirusTotal
type VirusTotal struct {
	client *http.Client
	apiKey string
}

func (v *VirusTotal) Name() string { return "VirusTotal" }

func (v *VirusTotal) Fetch(domain string) ([]string, error) {
	if v.apiKey == "" {
		return nil, fmt.Errorf("API key required")
	}

	apiURL := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s/urls", domain)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apikey", v.apiKey)

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			Attributes struct {
				URL string `json:"url"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, item := range result.Data {
		urls = append(urls, item.Attributes.URL)
	}
	return urls, nil
}

// SecurityTrails
type SecurityTrails struct {
	client *http.Client
	apiKey string
}

func (s *SecurityTrails) Name() string { return "SecurityTrails" }

func (s *SecurityTrails) Fetch(domain string) ([]string, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("API key required")
	}

	apiURL := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains", domain)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("APIKEY", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Subdomains []string `json:"subdomains"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, sub := range result.Subdomains {
		urls = append(urls, fmt.Sprintf("https://%s.%s", sub, domain))
	}
	return urls, nil
}

// GitHub
type GitHub struct {
	client *http.Client
	apiKey string
}

func (g *GitHub) Name() string { return "GitHub" }

func (g *GitHub) Fetch(domain string) ([]string, error) {
	query := url.QueryEscape(fmt.Sprintf("%s extension:json OR extension:txt OR extension:config", domain))
	apiURL := fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=100", query)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	if g.apiKey != "" {
		req.Header.Set("Authorization", "token "+g.apiKey)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			HTMLURL string `json:"html_url"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, item := range result.Items {
		urls = append(urls, item.HTMLURL)
	}
	return urls, nil
}

// ThreatCrowd
type ThreatCrowd struct {
	client *http.Client
}

func (t *ThreatCrowd) Name() string { return "ThreatCrowd" }

func (t *ThreatCrowd) Fetch(domain string) ([]string, error) {
	apiURL := fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain)
	resp, err := t.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Subdomains []string `json:"subdomains"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, sub := range result.Subdomains {
		urls = append(urls, "https://"+sub)
	}
	return urls, nil
}

// Censys
type Censys struct {
	client    *http.Client
	apiID     string
	apiSecret string
}

func (c *Censys) Name() string { return "Censys" }

func (c *Censys) Fetch(domain string) ([]string, error) {
	if c.apiID == "" || c.apiSecret == "" {
		return nil, fmt.Errorf("API credentials required")
	}

	apiURL := "https://search.censys.io/api/v2/certificates/search"
	query := fmt.Sprintf("names: %s", domain)
	
	reqBody := strings.NewReader(fmt.Sprintf(`{"q":"%s","per_page":100}`, query))
	req, err := http.NewRequest("POST", apiURL, reqBody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.apiID, c.apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			Hits []struct {
				Names []string `json:"names"`
			} `json:"hits"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	seen := make(map[string]bool)
	for _, hit := range result.Result.Hits {
		for _, name := range hit.Names {
			if !seen[name] && strings.Contains(name, domain) {
				urls = append(urls, "https://"+name)
				seen[name] = true
			}
		}
	}
	return urls, nil
}

func extractURL(text string) string {
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "http") {
			return strings.Trim(word, `"'<>()[]{}`)
		}
	}
	return ""
}

// ActiveCrawler for HTML/JS parsing and URL extraction
type ActiveCrawler struct {
	client        *http.Client
	visited       map[string]bool
	mu            sync.Mutex
	maxDepth      int
	verbose       bool
	allowedDomain string
}

func NewActiveCrawler(maxDepth int, timeout time.Duration, verbose bool, domain string) *ActiveCrawler {
	return &ActiveCrawler{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		visited:       make(map[string]bool),
		maxDepth:      maxDepth,
		verbose:       verbose,
		allowedDomain: domain,
	}
}

func (ac *ActiveCrawler) Crawl(targetURL string, depth int) *CrawlResult {
	ac.mu.Lock()
	if ac.visited[targetURL] || depth > ac.maxDepth {
		ac.mu.Unlock()
		return nil
	}
	ac.visited[targetURL] = true
	ac.mu.Unlock()

	if ac.verbose {
		fmt.Fprintf(os.Stderr, "[*] Crawling [depth=%d]: %s\n", depth, targetURL)
	}

	result := &CrawlResult{
		URL:   targetURL,
		Depth: depth,
	}

	resp, err := ac.client.Get(targetURL)
	if err != nil {
		if ac.verbose {
			fmt.Fprintf(os.Stderr, "[!] Error crawling %s: %v\n", targetURL, err)
		}
		return result
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	
	// Parse HTML
	if strings.Contains(contentType, "text/html") {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			ac.parseHTML(string(body), targetURL, result)
		}
	}

	// Parse JavaScript
	if strings.Contains(contentType, "javascript") || strings.Contains(contentType, "application/json") {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			ac.parseJavaScript(string(body), targetURL, result)
		}
	}

	if ac.verbose {
		fmt.Fprintf(os.Stderr, "[+] Found %d URLs, %d JS, %d CSS from %s\n", 
			len(result.ExtractedURLs), len(result.JSFiles), len(result.CSSFiles), targetURL)
	}

	return result
}

func (ac *ActiveCrawler) parseHTML(body, baseURL string, result *CrawlResult) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						resolvedURL := ac.resolveURL(baseURL, attr.Val)
						if resolvedURL != "" {
							result.Links = append(result.Links, resolvedURL)
							result.ExtractedURLs = append(result.ExtractedURLs, resolvedURL)
						}
					}
				}
			case "script":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resolvedURL := ac.resolveURL(baseURL, attr.Val)
						if resolvedURL != "" {
							result.JSFiles = append(result.JSFiles, resolvedURL)
							result.ExtractedURLs = append(result.ExtractedURLs, resolvedURL)
						}
					}
				}
				// Extract inline JavaScript
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					ac.parseJavaScript(n.FirstChild.Data, baseURL, result)
				}
			case "link":
				var rel, href string
				for _, attr := range n.Attr {
					if attr.Key == "rel" {
						rel = attr.Val
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if rel == "stylesheet" && href != "" {
					resolvedURL := ac.resolveURL(baseURL, href)
					if resolvedURL != "" {
						result.CSSFiles = append(result.CSSFiles, resolvedURL)
						result.ExtractedURLs = append(result.ExtractedURLs, resolvedURL)
					}
				}
			case "img":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resolvedURL := ac.resolveURL(baseURL, attr.Val)
						if resolvedURL != "" {
							result.Images = append(result.Images, resolvedURL)
							result.ExtractedURLs = append(result.ExtractedURLs, resolvedURL)
						}
					}
				}
			case "form":
				for _, attr := range n.Attr {
					if attr.Key == "action" {
						resolvedURL := ac.resolveURL(baseURL, attr.Val)
						if resolvedURL != "" {
							result.Forms = append(result.Forms, resolvedURL)
							result.ExtractedURLs = append(result.ExtractedURLs, resolvedURL)
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
}

func (ac *ActiveCrawler) parseJavaScript(jsCode, baseURL string, result *CrawlResult) {
	// Extract URLs from JavaScript using regex patterns
	patterns := []string{
		`["'](https?://[^"']+)["']`,
		`["']([/][^"']+)["']`,
		`url\s*\(\s*["']?([^"')]+)["']?\)`,
		`fetch\s*\(\s*["']([^"']+)["']`,
		`ajax\s*\(\s*["']([^"']+)["']`,
		`src\s*[:=]\s*["']([^"']+)["']`,
		`href\s*[:=]\s*["']([^"']+)["']`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(jsCode, -1)
		for _, match := range matches {
			if len(match) > 1 {
				resolvedURL := ac.resolveURL(baseURL, match[1])
				if resolvedURL != "" && !contains(result.ExtractedURLs, resolvedURL) {
					result.ExtractedURLs = append(result.ExtractedURLs, resolvedURL)
					
					// Categorize by extension
					ext := strings.ToLower(path.Ext(resolvedURL))
					switch ext {
					case ".js":
						result.JSFiles = append(result.JSFiles, resolvedURL)
					case ".css":
						result.CSSFiles = append(result.CSSFiles, resolvedURL)
					case ".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp":
						result.Images = append(result.Images, resolvedURL)
					}
				}
			}
		}
	}
}

func (ac *ActiveCrawler) resolveURL(baseURL, relativeURL string) string {
	relativeURL = strings.TrimSpace(relativeURL)
	
	// Skip invalid URLs
	if relativeURL == "" || relativeURL == "#" || strings.HasPrefix(relativeURL, "javascript:") ||
		strings.HasPrefix(relativeURL, "mailto:") || strings.HasPrefix(relativeURL, "tel:") ||
		strings.HasPrefix(relativeURL, "data:") {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	ref, err := url.Parse(relativeURL)
	if err != nil {
		return ""
	}

	resolved := base.ResolveReference(ref).String()

	// Only return URLs from the same domain if domain filter is set
	if ac.allowedDomain != "" {
		parsedURL, err := url.Parse(resolved)
		if err != nil || !strings.Contains(parsedURL.Hostname(), ac.allowedDomain) {
			return ""
		}
	}

	return resolved
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// LiveProbe performs comprehensive host probing
type LiveProbe struct {
	Timeout        time.Duration
	Ports          []int
	CustomHeaders  map[string]string
	FollowRedirect bool
	StatusCodes    []int
	Verbose        bool
}

func NewLiveProbe(timeout time.Duration, verbose bool) *LiveProbe {
	return &LiveProbe{
		Timeout:        timeout,
		Ports:          []int{80, 443, 8080, 8443},
		CustomHeaders:  make(map[string]string),
		FollowRedirect: true,
		StatusCodes:    []int{200, 201, 204, 301, 302, 307, 308, 401, 403},
		Verbose:        verbose,
	}
}

func (lp *LiveProbe) Probe(target string) *ProbeResult {
	result := &ProbeResult{
		URL: target,
	}

	parsedURL, err := url.Parse(target)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// DNS Resolution
	ctx, cancel := context.WithTimeout(context.Background(), lp.Timeout)
	defer cancel()

	ips, err := net.DefaultResolver.LookupHost(ctx, parsedURL.Hostname())
	if err != nil {
		result.Error = fmt.Sprintf("DNS: %v", err)
		return result
	}
	result.DNSResolved = true

	// TCP Port Probing
	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	conn, err := net.DialTimeout("tcp", parsedURL.Hostname()+":"+port, lp.Timeout)
	if err != nil {
		result.Error = fmt.Sprintf("TCP: %v", err)
		return result
	}
	conn.Close()
	result.TCPOpen = true

	// TLS Handshake (if HTTPS)
	if parsedURL.Scheme == "https" {
		tlsConn, err := tls.DialWithDialer(&net.Dialer{Timeout: lp.Timeout}, "tcp", parsedURL.Hostname()+":"+port, &tls.Config{
			InsecureSkipVerify: true,
		})
		if err == nil {
			result.TLSHandshake = true
			state := tlsConn.ConnectionState()
			result.TLSVersion = fmt.Sprintf("TLS %d.%d", state.Version>>8, state.Version&0xff)
			
			if len(state.PeerCertificates) > 0 {
				cert := state.PeerCertificates[0]
				result.Certificate = cert.Subject.CommonName
			}
			tlsConn.Close()
		}
	}

	// HTTP Request
	start := time.Now()
	
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: lp.Timeout,
		}).DialContext,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   lp.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !lp.FollowRedirect {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		result.Error = fmt.Sprintf("HTTP: %v", err)
		return result
	}

	req.Header.Set("User-Agent", "urlX/2.0")
	for k, v := range lp.CustomHeaders {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("HTTP: %v", err)
		return result
	}
	defer resp.Body.Close()

	result.ResponseTime = time.Since(start)
	result.StatusCode = resp.StatusCode
	result.ContentLength = resp.ContentLength
	result.Server = resp.Header.Get("Server")
	result.ContentType = resp.Header.Get("Content-Type")

	// Check if status code matches filter
	for _, code := range lp.StatusCodes {
		if result.StatusCode == code {
			result.IsAlive = true
			break
		}
	}

	// Response Inspection
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err == nil {
		if title := extractTitle(string(body)); title != "" {
			result.Title = title
		}
		result.Technologies = detectTechnologies(string(body), resp.Header)
	}

	return result
}

func extractTitle(htmlContent string) string {
	start := strings.Index(strings.ToLower(htmlContent), "<title")
	if start == -1 {
		return ""
	}
	start = strings.Index(htmlContent[start:], ">")
	if start == -1 {
		return ""
	}
	start += strings.Index(htmlContent, "<title")
	end := strings.Index(htmlContent[start:], "</title>")
	if end == -1 {
		return ""
	}
	title := htmlContent[start+1 : start+end]
	title = strings.TrimSpace(title)
	if len(title) > 100 {
		title = title[:100] + "..."
	}
	return title
}

func detectTechnologies(body string, headers http.Header) []string {
	techs := make([]string, 0)
	body = strings.ToLower(body)

	checks := map[string]string{
		"WordPress":  "wp-content",
		"Joomla":     "joomla",
		"Drupal":     "drupal",
		"React":      "react",
		"Vue.js":     "vue",
		"Angular":    "angular",
		"jQuery":     "jquery",
		"Bootstrap":  "bootstrap",
		"Laravel":    "laravel",
		"Django":     "csrfmiddlewaretoken",
		"Flask":      "werkzeug",
		"Express":    "x-powered-by: express",
	}

	for tech, indicator := range checks {
		if strings.Contains(body, indicator) || strings.Contains(strings.ToLower(headers.Get("X-Powered-By")), strings.ToLower(indicator)) {
			techs = append(techs, tech)
		}
	}

	return techs
}

// FileFilter handles file extension filtering
type FileFilter struct {
	OnlyExtensions   []string
	RemoveExtensions []string
}

func NewFileFilter(only, remove string) *FileFilter {
	ff := &FileFilter{
		OnlyExtensions:   make([]string, 0),
		RemoveExtensions: make([]string, 0),
	}

	if only != "" {
		for _, ext := range strings.Split(only, ",") {
			ext = strings.TrimSpace(ext)
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			ff.OnlyExtensions = append(ff.OnlyExtensions, strings.ToLower(ext))
		}
	}

	if remove != "" {
		for _, ext := range strings.Split(remove, ",") {
			ext = strings.TrimSpace(ext)
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			ff.RemoveExtensions = append(ff.RemoveExtensions, strings.ToLower(ext))
		}
	}

	return ff
}

func (ff *FileFilter) ShouldInclude(urlStr string) bool {
	ext := strings.ToLower(path.Ext(urlStr))
	
	// Remove filter takes precedence
	if len(ff.RemoveExtensions) > 0 {
		for _, removeExt := range ff.RemoveExtensions {
			if ext == removeExt {
				return false
			}
		}
	}

	// Only filter
	if len(ff.OnlyExtensions) > 0 {
		for _, onlyExt := range ff.OnlyExtensions {
			if ext == onlyExt {
				return true
			}
		}
		return false
	}

	return true
}

// URLCollector manages fetching from multiple sources
type URLCollector struct {
	sources []URLSource
	verbose bool
}

func NewURLCollector(verbose bool, config map[string]string) *URLCollector {
	client := &http.Client{Timeout: 30 * time.Second}

	sources := []URLSource{
		&WaybackMachine{client: client},
		&ArchiveToday{client: client},
		&CommonCrawl{client: client},
		&URLScan{client: client, apiKey: config["urlscan"]},
		&AlienVaultOTX{client: client, apiKey: config["otx"]},
		&CrtSh{client: client},
		&ThreatCrowd{client: client},
	}

	if config["virustotal"] != "" {
		sources = append(sources, &VirusTotal{client: client, apiKey: config["virustotal"]})
	}
	if config["securitytrails"] != "" {
		sources = append(sources, &SecurityTrails{client: client, apiKey: config["securitytrails"]})
	}
	if config["github"] != "" {
		sources = append(sources, &GitHub{client: client, apiKey: config["github"]})
	}
	if config["censys_id"] != "" && config["censys_secret"] != "" {
		sources = append(sources, &Censys{
			client:    client,
			apiID:     config["censys_id"],
			apiSecret: config["censys_secret"],
		})
	}

	return &URLCollector{sources: sources, verbose: verbose}
}

func (uc *URLCollector) FetchURLs(domain string) []string {
	var wg sync.WaitGroup
	urlChan := make(chan string, 1000)
	uniqueURLs := make(map[string]bool)
	var mu sync.Mutex

	for _, source := range uc.sources {
		wg.Add(1)
		go func(s URLSource) {
			defer wg.Done()

			if uc.verbose {
				fmt.Fprintf(os.Stderr, "[*] Fetching from %s for %s...\n", s.Name(), domain)
			}

			urls, err := s.Fetch(domain)
			if err != nil {
				if uc.verbose {
					fmt.Fprintf(os.Stderr, "[!] Error from %s: %v\n", s.Name(), err)
				}
				return
			}

			if uc.verbose {
				fmt.Fprintf(os.Stderr, "[+] Found %d URLs from %s\n", len(urls), s.Name())
			}

			for _, u := range urls {
				urlChan <- u
			}
		}(source)
	}

	go func() {
		wg.Wait()
		close(urlChan)
	}()

	for u := range urlChan {
		mu.Lock()
		if !uniqueURLs[u] {
			uniqueURLs[u] = true
		}
		mu.Unlock()
	}

	result := make([]string, 0, len(uniqueURLs))
	for u := range uniqueURLs {
		result = append(result, u)
	}
	return result
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func normalizeInput(input string) string {
	input = strings.TrimSpace(input)
	if strings.Contains(input, "://") {
		if u, err := url.Parse(input); err == nil {
			return u.String()
		}
	}
	return input
}

func main() {
	var (
		domainFile      = flag.String("d", "", "File with domains/subdomains")
		urlsFile        = flag.String("u", "", "File with URLs")
		domain          = flag.String("domain", "", "Single domain to query")
		output          = flag.String("o", "", "Output file (default: stdout)")
		verbose         = flag.Bool("v", false, "Verbose output")
		probe           = flag.Bool("probe", false, "Enable live host probing")
		crawl           = flag.Bool("crawl", false, "Enable active crawling")
		crawlDepth      = flag.Int("depth", 2, "Crawling depth")
		probeTimeout    = flag.Duration("timeout", 10*time.Second, "Probe timeout")
		statusCodes     = flag.String("mc", "200,201,204,301,302,307,308,401,403", "Match status codes")
		concurrency     = flag.Int("c", 5, "Concurrency level")
		onlyExt         = flag.String("only", "", "Only include these extensions (e.g., '.pdf,.doc')")
		removeExt       = flag.String("remove", "", "Remove these extensions (e.g., '.png,.jpg')")
		urlscanKey      = flag.String("urlscan-key", "", "URLScan.io API key")
		otxKey          = flag.String("otx-key", "", "AlienVault OTX API key")
		vtKey           = flag.String("vt-key", "", "VirusTotal API key")
		sectrailsKey    = flag.String("securitytrails-key", "", "SecurityTrails API key")
		githubKey       = flag.String("github-key", "", "GitHub API token")
		censysID        = flag.String("censys-id", "", "Censys API ID")
		censysSecret    = flag.String("censys-secret", "", "Censys API Secret")
		noBanner        = flag.Bool("no-banner", false, "Disable banner")
		jsonOutput      = flag.Bool("json", false, "Output in JSON format")
	)

	flag.Parse()

	if !*noBanner {
		fmt.Fprint(os.Stderr, banner)
	}

	if *domainFile == "" && *urlsFile == "" && *domain == "" {
		fmt.Fprintln(os.Stderr, "Error: Provide -d (domains), -u (URLs), or -domain (single domain)")
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  urlx -domain example.com -crawl -v")
		fmt.Fprintln(os.Stderr, "  urlx -d subdomains.txt -probe -only '.pdf,.doc'")
		fmt.Fprintln(os.Stderr, "  urlx -u urls.txt -crawl -depth 3 -remove '.png,.jpg'")
		flag.Usage()
		os.Exit(1)
	}

	config := map[string]string{
		"urlscan":        *urlscanKey,
		"otx":            *otxKey,
		"virustotal":     *vtKey,
		"securitytrails": *sectrailsKey,
		"github":         *githubKey,
		"censys_id":      *censysID,
		"censys_secret":  *censysSecret,
	}

	// Parse status codes
	allowedCodes := make([]int, 0)
	for _, code := range strings.Split(*statusCodes, ",") {
		var c int
		fmt.Sscanf(strings.TrimSpace(code), "%d", &c)
		if c > 0 {
			allowedCodes = append(allowedCodes, c)
		}
	}

	// File filter
	fileFilter := NewFileFilter(*onlyExt, *removeExt)

	var writer io.Writer = os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		writer = f
	}

	var allTargets []string

	// Process input files or single domain
	if *urlsFile != "" {
		urls, err := readLines(*urlsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading URLs file: %v\n", err)
			os.Exit(1)
		}
		allTargets = append(allTargets, urls...)
		
		if *verbose {
			fmt.Fprintf(os.Stderr, "[*] Loaded %d URLs from file\n", len(urls))
		}
	} else {
		collector := NewURLCollector(*verbose, config)
		
		var domains []string
		if *domain != "" {
			domains = []string{*domain}
		} else {
			var err error
			domains, err = readLines(*domainFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading domain file: %v\n", err)
				os.Exit(1)
			}
		}

		if *verbose {
			fmt.Fprintf(os.Stderr, "[*] Loaded %d domains\n", len(domains))
			fmt.Fprintf(os.Stderr, "[*] Active sources: %d\n\n", len(collector.sources))
		}

		// URL Discovery phase
		domainChan := make(chan string, len(domains))
		for _, d := range domains {
			if strings.Contains(d, "://") {
				if u, err := url.Parse(d); err == nil {
					d = u.Hostname()
				}
			}
			domainChan <- d
		}
		close(domainChan)

		var wg sync.WaitGroup
		var mu sync.Mutex
		discoveredURLs := make(map[string]bool)

		for i := 0; i < *concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for domain := range domainChan {
					urls := collector.FetchURLs(domain)
					urls = append(urls, "http://"+domain, "https://"+domain)
					
					mu.Lock()
					for _, u := range urls {
						discoveredURLs[u] = true
					}
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		for u := range discoveredURLs {
			allTargets = append(allTargets, u)
		}

		if *verbose {
			fmt.Fprintf(os.Stderr, "\n[+] Total URLs discovered: %d\n\n", len(allTargets))
		}
	}

	// Apply file filters
	filteredTargets := make([]string, 0)
	for _, target := range allTargets {
		if fileFilter.ShouldInclude(target) {
			filteredTargets = append(filteredTargets, target)
		}
	}

	if *verbose && len(filteredTargets) != len(allTargets) {
		fmt.Fprintf(os.Stderr, "[*] Filtered %d URLs based on extension filters\n", len(allTargets)-len(filteredTargets))
		fmt.Fprintf(os.Stderr, "[*] Remaining URLs: %d\n\n", len(filteredTargets))
	}

	allTargets = filteredTargets

	// Active Crawling phase
	if *crawl {
		if *verbose {
			fmt.Fprintf(os.Stderr, "[*] Starting active crawling (depth: %d)...\n\n", *crawlDepth)
		}

		var crawlDomain string
		if *domain != "" {
			crawlDomain = *domain
		}

		crawler := NewActiveCrawler(*crawlDepth, *probeTimeout, *verbose, crawlDomain)
		
		targetChan := make(chan string, len(allTargets))
		for _, t := range allTargets {
			targetChan <- normalizeInput(t)
		}
		close(targetChan)

		var wg sync.WaitGroup
		var mu sync.Mutex
		crawledURLs := make(map[string]bool)

		for i := 0; i < *concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for target := range targetChan {
					result := crawler.Crawl(target, 0)
					if result != nil {
						mu.Lock()
						for _, u := range result.ExtractedURLs {
							if fileFilter.ShouldInclude(u) {
								crawledURLs[u] = true
							}
						}
						mu.Unlock()
					}
				}
			}()
		}

		wg.Wait()

		// Add crawled URLs to targets
		for u := range crawledURLs {
			if !contains(allTargets, u) {
				allTargets = append(allTargets, u)
			}
		}

		if *verbose {
			fmt.Fprintf(os.Stderr, "\n[+] Total URLs after crawling: %d\n\n", len(allTargets))
		}
	}

	// Probing phase
	if *probe {
		if *verbose {
			fmt.Fprintf(os.Stderr, "[*] Starting live host probing...\n")
			fmt.Fprintf(os.Stderr, "[*] Timeout: %v | Status codes: %s\n\n", *probeTimeout, *statusCodes)
		}

		prober := NewLiveProbe(*probeTimeout, *verbose)
		prober.StatusCodes = allowedCodes

		targetChan := make(chan string, len(allTargets))
		for _, t := range allTargets {
			targetChan <- normalizeInput(t)
		}
		close(targetChan)

		var wg sync.WaitGroup
		var mu sync.Mutex
		results := make([]*ProbeResult, 0)

		for i := 0; i < *concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for target := range targetChan {
					result := prober.Probe(target)
					
					mu.Lock()
					results = append(results, result)
					mu.Unlock()

					if result.IsAlive {
						if *jsonOutput {
							mu.Lock()
							jsonData, _ := json.Marshal(result)
							fmt.Fprintln(writer, string(jsonData))
							mu.Unlock()
						} else {
							output := fmt.Sprintf("[%d] %s", result.StatusCode, result.URL)
							if result.Title != "" {
								output += fmt.Sprintf(" [%s]", result.Title)
							}
							if result.Server != "" {
								output += fmt.Sprintf(" [%s]", result.Server)
							}
							if len(result.Technologies) > 0 {
								output += fmt.Sprintf(" [%s]", strings.Join(result.Technologies, ","))
							}
							output += fmt.Sprintf(" [%dms]", result.ResponseTime.Milliseconds())
							
							mu.Lock()
							fmt.Fprintln(writer, output)
							mu.Unlock()
						}
					}
				}
			}()
		}

		wg.Wait()

		// Summary
		if *verbose {
			alive := 0
			dnsResolved := 0
			tcpOpen := 0
			tlsSuccess := 0

			for _, r := range results {
				if r.IsAlive {
					alive++
				}
				if r.DNSResolved {
					dnsResolved++
				}
				if r.TCPOpen {
					tcpOpen++
				}
				if r.TLSHandshake {
					tlsSuccess++
				}
			}

			fmt.Fprintf(os.Stderr, "\n[*] Probing Summary:\n")
			fmt.Fprintf(os.Stderr, "    Total targets: %d\n", len(results))
			fmt.Fprintf(os.Stderr, "    DNS resolved: %d\n", dnsResolved)
			fmt.Fprintf(os.Stderr, "    TCP open: %d\n", tcpOpen)
			fmt.Fprintf(os.Stderr, "    TLS success: %d\n", tlsSuccess)
			fmt.Fprintf(os.Stderr, "    Alive hosts: %d\n", alive)
		}
	} else {
		// No probing - just output URLs
		for _, target := range allTargets {
			fmt.Fprintln(writer, target)
		}

		if *verbose {
			fmt.Fprintf(os.Stderr, "\n[+] Total unique URLs: %d\n", len(allTargets))
		}
	}
}
