# urlX - Complete Feature List

**Author**: [Alham Rizvi](https://github.com/alhamrizvi-cloud)  
**Repository**: https://github.com/alhamrizvi-cloud/urlx

## 🎯 Core Capabilities

urlX is a comprehensive reconnaissance tool that combines three powerful engines:

1. **Passive URL Discovery** - 11+ data sources
2. **Active Web Crawling** - HTML/JavaScript parsing
3. **Live Host Probing** - DNS, TCP, TLS, HTTP validation

## 📦 Passive Discovery Engine

### Supported Data Sources (11+)

#### Historical & Archive Sources
✅ **Wayback Machine** (archive.org)
- Historical snapshots of web pages
- URL history across years
- No API key required

✅ **Archive.today** (archive.is/archive.ph)
- Recent snapshots
- Bypasses some robots.txt restrictions
- No API key required

✅ **Common Crawl**
- Massive web crawl datasets
- Petabytes of web data
- Free access via S3

#### Certificate Intelligence
✅ **crt.sh**
- Certificate Transparency logs
- Subdomain discovery via SSL/TLS certs
- Finds staging/dev environments
- No API key required

✅ **Censys**
- SSL/TLS certificate search
- Historical certificate data
- API key required (free tier: 250 queries/month)

#### Threat Intelligence
✅ **AlienVault OTX**
- Open Threat Exchange
- Community-contributed URLs
- Malware-related endpoints
- API key optional (free)

✅ **VirusTotal**
- Historical domain/URL analysis
- Threat intelligence data
- API key required (free tier available)

✅ **ThreatCrowd**
- Community threat intelligence
- Subdomain enumeration
- No API key required

#### Security Platforms
✅ **URLScan.io**
- Website scanner submissions
- Recent scans and screenshots
- API key optional (free tier)

✅ **SecurityTrails**
- DNS history
- Subdomain enumeration
- API key required (limited free tier)

#### Code Repositories
✅ **GitHub Code Search**
- Find endpoints in config files
- API keys, URLs in code
- Personal access token required (free)

### Discovery Features
- **Concurrent querying** - All sources queried in parallel
- **Smart deduplication** - Removes duplicates across all sources
- **Domain normalization** - Handles various input formats
- **Automatic retries** - Handles transient failures
- **Rate limiting** - Respects API rate limits
- **Verbose logging** - See exactly what's being discovered

## 🕷️ Active Crawling Engine

### HTML Parser (golang.org/x/net/html)

#### Extraction Capabilities
✅ **Links** - All `<a href="">` tags
- Internal and external links
- Relative and absolute URLs
- Fragment identifiers handled

✅ **JavaScript Files** - `<script src="">`
- External scripts
- CDN references
- Inline script processing

✅ **CSS Files** - `<link rel="stylesheet">`
- Stylesheet references
- CSS imports
- Asset URLs

✅ **Images** - `<img src="">`
- Image sources
- srcset attributes
- Picture elements

✅ **Forms** - `<form action="">`
- Form submission endpoints
- Method detection
- Hidden form fields

✅ **Inline JavaScript** - `<script>` content
- Embedded code analysis
- URL extraction from JS

### JavaScript Parser (Regex-based)

#### Pattern Matching
✅ **URL Patterns**
```javascript
"https://api.example.com/users"
'/api/v1/posts'
`/endpoint/${id}`
```

✅ **Fetch Calls**
```javascript
fetch('/api/data')
fetch("https://api.example.com/users")
```

✅ **AJAX Requests**
```javascript
$.ajax('/api/users')
axios.get('/api/posts')
$.get('/data')
```

✅ **URL Assignments**
```javascript
src = "/assets/image.png"
href = "https://example.com/page"
window.location = "/redirect"
```

✅ **CSS url() Functions**
```javascript
url('/images/bg.png')
url("https://cdn.example.com/font.woff")
background: url('/pattern.svg')
```

✅ **API Endpoints**
```javascript
const API_URL = "https://api.example.com"
endpoint: '/v1/users'
baseURL: 'https://api.example.com'
```

### Crawling Features
- **Configurable depth** - Control crawl levels (1-5+)
- **Domain filtering** - Stay within target domain
- **Visited tracking** - Avoid duplicate crawling
- **Content-type detection** - Handle HTML, JS, CSS differently
- **Concurrent crawling** - Parallel page processing
- **URL resolution** - Handles relative URLs correctly
- **Breadth-first search** - Systematic exploration
- **Memory efficient** - Smart caching and cleanup

### Asset Categorization
Automatically categorizes discovered URLs:
- JavaScript files (`.js`)
- CSS files (`.css`)
- Images (`.jpg`, `.png`, `.gif`, `.svg`, `.webp`)
- Documents (`.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`)
- Data files (`.json`, `.xml`, `.yml`, `.yaml`)
- Archives (`.zip`, `.tar`, `.gz`, `.rar`)
- Config files (`.conf`, `.config`, `.ini`, `.env`)

## 📁 File Filtering System

### Include Filter (`--only`)

Extract only specific file types:

```bash
# Single extension
urlx -domain example.com -crawl -only ".pdf"

# Multiple extensions
urlx -domain example.com -crawl -only ".pdf,.doc,.docx,.xls,.xlsx"

# Code files
urlx -domain example.com -crawl -only ".js,.json,.xml,.yml"
```

**Use Cases:**
- PDF extraction: `-only ".pdf"`
- Documents: `-only ".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx"`
- Code: `-only ".js,.css,.html,.php,.py,.rb,.go"`
- Data: `-only ".json,.xml,.yml,.yaml,.csv"`
- Config: `-only ".conf,.config,.ini,.env,.properties"`
- Archives: `-only ".zip,.tar,.gz,.rar,.7z,.bak"`

### Exclude Filter (`--remove`)

Remove unwanted file types:

```bash
# Remove images
urlx -domain example.com -crawl -remove ".png,.jpg,.gif,.svg"

# Remove media
urlx -domain example.com -crawl -remove ".mp4,.mp3,.avi,.mov,.wmv"

# Remove common assets
urlx -domain example.com -crawl -remove ".png,.jpg,.css,.woff,.ttf,.eot"
```

**Use Cases:**
- No images: `-remove ".png,.jpg,.jpeg,.gif,.svg,.webp,.ico"`
- No media: `-remove ".mp4,.mp3,.avi,.mov,.wmv,.flv,.webm"`
- No fonts: `-remove ".woff,.woff2,.ttf,.eot,.otf"`
- No minified: `-remove ".min.js,.min.css"`

### Filter Features
- **Smart extension handling** - Auto-adds `.` if missing
- **Real-time filtering** - Applied during discovery/crawling
- **Case-insensitive** - Works with any case
- **Combinable** - Use `--only` and `--remove` together
- **Works with all modes** - Discovery, crawling, probing

## 🔍 Live Host Probing Engine

### Layer 1: DNS Resolution

✅ **IPv4 Support**
- A record lookups
- Standard DNS resolution

✅ **IPv6 Support**
- AAAA record lookups
- Dual-stack support

✅ **Features**
- Timeout handling
- Retry logic
- Custom resolvers support
- Validation of DNS responses

### Layer 2: TCP Socket Probing

✅ **Port Testing**
- Default ports: 80, 443, 8080, 8443
- Configurable port list
- Connection establishment verification

✅ **Metrics**
- Connection latency
- Port open/closed detection
- Timeout handling

### Layer 3: TLS Handshake (HTTPS)

✅ **Certificate Extraction**
- Common Name (CN)
- Subject Alternative Names (SANs)
- Certificate chain
- Expiration dates

✅ **TLS Information**
- TLS version detection (1.0, 1.1, 1.2, 1.3)
- Cipher suite information
- Certificate validation
- InsecureSkipVerify for testing

### Layer 4: HTTP Request Engine

✅ **Request Methods**
- GET requests (default)
- HEAD requests (optional)
- Custom headers support
- User-Agent customization

✅ **Redirect Handling**
- Follow redirects (configurable)
- Redirect chain tracking
- 301, 302, 307, 308 support

✅ **Custom Headers**
```go
-header "Authorization: Bearer token"
-header "X-API-Key: key"
```

### Layer 5: Response Inspection

✅ **Status Code Analysis**
- Configurable status code matching
- Default: 200, 201, 204, 301, 302, 307, 308, 401, 403
- Custom filters: `-mc 200,403,500`

✅ **Content Analysis**
- Page title extraction
- Response body length
- Content-Type detection
- Server header extraction

✅ **Technology Detection**
Automatically detects:
- **CMS**: WordPress, Joomla, Drupal
- **Frameworks**: React, Vue.js, Angular, Laravel, Django, Flask
- **Libraries**: jQuery, Bootstrap
- **Servers**: nginx, Apache, IIS
- **Languages**: Express (Node.js), PHP indicators

### Probing Features
- **Massive concurrency** - Goroutine-based parallel probing
- **Response time tracking** - Millisecond precision
- **Connection pooling** - Efficient resource usage
- **Timeout controls** - Configurable per-request
- **Error handling** - Detailed error reporting
- **Rate limiting** - Avoid overwhelming targets

## 🎨 Output Formats

### Standard Text Output
```
[200] https://example.com [Example Domain] [nginx/1.18.0] [React] [245ms]
[403] https://admin.example.com [Forbidden] [Apache] [312ms]
[301] https://www.example.com [Moved] [189ms]
```

### JSON Output
```json
{
  "URL": "https://example.com",
  "StatusCode": 200,
  "ContentLength": 1256,
  "Title": "Example Domain",
  "Server": "nginx/1.18.0",
  "ContentType": "text/html; charset=UTF-8",
  "Technologies": ["nginx", "React", "Bootstrap"],
  "TLSVersion": "TLS 1.3",
  "Certificate": "example.com",
  "ResponseTime": 245000000,
  "IsAlive": true,
  "DNSResolved": true,
  "TCPOpen": true,
  "TLSHandshake": true,
  "Error": ""
}
```

### Verbose Output
```
[*] Loaded 3 domains
[*] Active sources: 11
[*] Fetching from Wayback Machine for example.com...
[+] Found 1247 URLs from Wayback Machine
[*] Crawling [depth=0]: https://example.com
[+] Found 45 URLs, 12 JS, 3 CSS
[*] Resolving DNS for example.com
[+] DNS resolved: [93.184.216.34]
[*] Probing TCP example.com:443
[+] TCP port 443 open
[*] Performing TLS handshake
[+] TLS handshake successful, CN: example.com
[+] HTTP 200 | Example Domain | 245ms
```

## ⚙️ Configuration Options

### Input Methods
- Single domain: `-domain example.com`
- Domain file: `-d domains.txt`
- URL file: `-u urls.txt`
- Stdin support: `cat domains.txt | urlx -d -`

### Performance Tuning
- Concurrency: `-c 20` (1-100+)
- Timeout: `-timeout 10s` (1s-60s)
- Depth: `-depth 3` (1-10)
- Rate limiting: Built-in per source

### Output Control
- Output file: `-o results.txt`
- JSON format: `-json`
- Verbose mode: `-v`
- Silent mode: `-no-banner`

### Filtering
- Status codes: `-mc 200,403,500`
- Include extensions: `-only ".pdf,.doc"`
- Exclude extensions: `-remove ".png,.jpg"`

## 🔧 Integration Features

### Works With
- **subfinder** - Subdomain enumeration
- **httpx** - HTTP probing
- **nuclei** - Vulnerability scanning
- **ffuf** - Fuzzing
- **waybackurls** - Additional URL discovery
- **gau** - GetAllUrls
- **aquatone** - Screenshots
- **meg** - Fetch many paths

### Scriptable
- Shell scripts
- Python automation
- CI/CD pipelines
- Cron jobs
- Docker containers

### Data Export
- Plain text (default)
- JSON (structured)
- CSV (with jq)
- Markdown reports
- Custom formats

## 📊 Performance Characteristics

### Speed
- **Discovery**: 11 sources in parallel
- **Crawling**: Concurrent page processing
- **Probing**: Goroutine-based parallelism
- **Typical**: 100-500 URLs/second with `-c 20`

### Resource Usage
- **Memory**: ~50-200MB typical
- **CPU**: Multi-core utilization
- **Network**: Configurable concurrency
- **Disk**: Minimal (results only)

### Scalability
- **Small targets**: <100 URLs - seconds
- **Medium targets**: 1K-10K URLs - minutes
- **Large targets**: 100K+ URLs - optimized batching

## 🛡️ Security Features

### Safe Defaults
- TLS verification (can be disabled)
- Rate limiting per source
- Timeout protections
- Domain filtering

### Privacy
- No data sent to third parties
- All processing local
- API keys user-controlled
- No telemetry

## 📈 Use Case Matrix

| Use Case | Discovery | Crawling | Probing | Filtering |
|----------|-----------|----------|---------|-----------|
| Bug Bounty Recon | ✅ | ✅ | ✅ | Optional |
| Subdomain Enum | ✅ | ❌ | ✅ | ❌ |
| Asset Discovery | ✅ | ✅ | ✅ | ✅ |
| PDF Extraction | ✅ | ✅ | ✅ | ✅ `.pdf` |
| JS Analysis | ✅ | ✅ | Optional | ✅ `.js` |
| API Discovery | ✅ | ✅ | ✅ | ✅ `.json` |
| Config Hunt | ✅ | ✅ | ✅ | ✅ `.conf,.env` |
| Tech Stack | ✅ | ✅ | ✅ | ❌ |

---

**Created by [Alham Rizvi](https://github.com/alhamrizvi-cloud)**  
**Repository**: https://github.com/alhamrizvi-cloud/urlx  
**License**: MIT
