# urlX Quick Start Guide

Get started with urlX in minutes!

**Created by [Alham Rizvi](https://github.com/alhamrizvi-cloud)**

## 🚀 Installation

### Option 1: Go Install (Recommended)
```bash
go install github.com/alhamrizvi-cloud/urlx@latest
```

### Option 2: Quick Install Script
```bash
curl -sSfL https://raw.githubusercontent.com/alhamrizvi-cloud/urlx/main/install.sh | sh
```

### Option 3: Build from Source
```bash
git clone https://github.com/alhamrizvi-cloud/urlx.git
cd urlx
go mod download
go build -o urlx main.go
sudo mv urlx /usr/local/bin/
```

### Option 4: Docker
```bash
docker pull alhamrizvi/urlx:latest
docker run --rm alhamrizvi/urlx -domain example.com -v
```

## ⚡ Quick Examples

### 1. Basic URL Discovery
```bash
# Discover URLs for a single domain
urlx -domain example.com

# From subdomains file
urlx -d subdomains.txt -o urls.txt

# With verbose output
urlx -domain example.com -v
```

### 2. Active Crawling
```bash
# Crawl a domain
urlx -domain example.com -crawl

# Deep crawl
urlx -domain example.com -crawl -depth 3 -v

# Crawl and save results
urlx -u urls.txt -crawl -depth 2 -o crawled.txt
```

### 3. File Extraction
```bash
# Extract only PDFs
urlx -domain example.com -crawl -only ".pdf" -o pdfs.txt

# Extract documents
urlx -domain example.com -crawl -only ".pdf,.doc,.docx"

# Extract JavaScript (no minified)
urlx -domain example.com -crawl -only ".js" -remove ".min.js"

# Get everything except images
urlx -domain example.com -crawl -remove ".png,.jpg,.gif"
```

### 4. Live Host Probing
```bash
# Probe for live hosts
urlx -d subdomains.txt -probe -v

# Custom status codes
urlx -u urls.txt -probe -mc 200,403

# Fast probing with high concurrency
urlx -d domains.txt -probe -c 20 -timeout 5s

# JSON output
urlx -d subdomains.txt -probe -json -o results.json
```

### 5. Combined Workflows
```bash
# Discover + Crawl + Probe (all in one)
urlx -domain example.com -crawl -depth 2 -probe -v

# Extract PDFs and verify accessibility
urlx -domain example.com -crawl -only ".pdf" -probe -mc 200 -o live-pdfs.txt

# Complete recon with filtering
urlx -d subs.txt -crawl -depth 3 -probe -mc 200,403 -remove ".png,.jpg" -o results.txt
```

## 🎯 Common Use Cases

### Bug Bounty Recon
```bash
# Complete workflow
subfinder -d target.com -silent > subs.txt
urlx -d subs.txt -crawl -depth 2 -probe -v -mc 200,403 -o alive.txt
cat alive.txt | httpx -silent | nuclei -t cves/
```

### Find Hidden Files
```bash
# Extract all documents
urlx -domain target.com -crawl -only ".pdf,.doc,.xls,.zip,.bak" -o files.txt

# Probe to verify they're accessible
urlx -u files.txt -probe -mc 200 -o accessible.txt
```

### JavaScript Analysis
```bash
# Get all non-minified JS files
urlx -domain target.com -crawl -depth 3 -only ".js" -remove ".min.js" -o js.txt

# Download and search for secrets
cat js.txt | while read js; do
    curl -s "$js" | grep -E "api_key|secret|token"
done
```

### API Discovery
```bash
# Find API endpoints
urlx -domain api.target.com -crawl -depth 3 | grep -E '/api/|/v[0-9]+/' > apis.txt

# Test them
urlx -u apis.txt -probe -mc 200,401,403 -json > api-results.json

# Analyze
jq -r 'select(.StatusCode == 200) | .URL' api-results.json > public-apis.txt
```

## 📊 Understanding Output

### Standard Format
```
[200] https://example.com [Example Domain] [nginx] [React] [245ms]
[403] https://admin.example.com [Forbidden] [Apache] [312ms]
```
- `[200]` - HTTP status code
- URL
- `[Example Domain]` - Page title
- `[nginx]` - Server
- `[React]` - Detected technologies
- `[245ms]` - Response time

### JSON Format
```bash
urlx -domain example.com -probe -json | jq '.'
```
```json
{
  "URL": "https://example.com",
  "StatusCode": 200,
  "Title": "Example Domain",
  "Server": "nginx",
  "Technologies": ["React", "nginx"],
  "ResponseTime": 245000000,
  "IsAlive": true
}
```

## 🔧 Useful Flags

| Flag | What it does | Example |
|------|-------------|---------|
| `-d` | Input file with domains | `-d subs.txt` |
| `-u` | Input file with URLs | `-u urls.txt` |
| `-domain` | Single domain | `-domain example.com` |
| `-crawl` | Enable crawling | `-crawl` |
| `-depth` | Crawl depth | `-depth 3` |
| `-probe` | Enable probing | `-probe` |
| `-only` | Include only these | `-only ".pdf,.doc"` |
| `-remove` | Exclude these | `-remove ".png,.jpg"` |
| `-mc` | Match status codes | `-mc 200,403` |
| `-o` | Output file | `-o results.txt` |
| `-v` | Verbose mode | `-v` |
| `-c` | Concurrency | `-c 20` |
| `-timeout` | Probe timeout | `-timeout 10s` |
| `-json` | JSON output | `-json` |

## 💡 Pro Tips

### 1. Start Simple, Add Complexity
```bash
# Start with basic discovery
urlx -domain example.com -o urls.txt

# Then crawl
urlx -u urls.txt -crawl -depth 2 -o crawled.txt

# Finally probe
urlx -u crawled.txt -probe -v
```

### 2. Use File Filters Wisely
```bash
# For code review
urlx -domain example.com -crawl -only ".js,.json,.xml"

# For document extraction
urlx -domain example.com -crawl -only ".pdf,.doc,.xls"

# To reduce noise
urlx -domain example.com -crawl -remove ".png,.jpg,.gif,.css,.woff"
```

### 3. Adjust Concurrency
```bash
# Fast scanning
urlx -d domains.txt -crawl -probe -c 50

# Gentle (avoid rate limits)
urlx -d domains.txt -crawl -probe -c 3

# Balanced (default)
urlx -d domains.txt -crawl -probe -c 10
```

### 4. Combine with Other Tools
```bash
# With subfinder
subfinder -d target.com -silent | urlx -d - -crawl -probe

# With httpx
urlx -d domains.txt -crawl -o urls.txt
cat urls.txt | httpx -silent -mc 200

# With nuclei
urlx -d subs.txt -probe -mc 200 -o alive.txt
nuclei -l alive.txt -t exposures/
```

## 🔑 API Keys (Optional but Recommended)

Set environment variables or use flags:

```bash
# Using flags
urlx -domain example.com \
  -urlscan-key "YOUR_KEY" \
  -otx-key "YOUR_KEY" \
  -vt-key "YOUR_KEY"

# Or set environment variables
export URLSCAN_KEY="your_key"
export OTX_KEY="your_key"
export VT_KEY="your_key"
```

Get free API keys:
- URLScan.io: https://urlscan.io/user/signup
- AlienVault OTX: https://otx.alienvault.com/api
- VirusTotal: https://www.virustotal.com/gui/my-apikey
- GitHub: https://github.com/settings/tokens

## 📚 Next Steps

1. **Read the full README**: [README.md](README.md)
2. **Check probing guide**: [PROBING-GUIDE.md](PROBING-GUIDE.md)
3. **Learn crawling**: [CRAWLING-GUIDE.md](CRAWLING-GUIDE.md)
4. **See examples**: [examples/](examples/)

## ❓ Common Issues

### "Binary not found after install"
```bash
# Add to PATH
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

### "Too many open files"
```bash
# Increase limit
ulimit -n 10000
```

### "Connection timeout"
```bash
# Increase timeout
urlx -u urls.txt -probe -timeout 30s
```

### "Rate limited"
```bash
# Reduce concurrency
urlx -d domains.txt -crawl -probe -c 3
```

## 🆘 Getting Help

```bash
# Show help
urlx -h

# Verbose mode shows what's happening
urlx -domain example.com -crawl -probe -v
```

## 📬 Support

- **Issues**: [github.com/alhamrizvi-cloud/urlx/issues](https://github.com/alhamrizvi-cloud/urlx/issues)
- **Discussions**: [github.com/alhamrizvi-cloud/urlx/discussions](https://github.com/alhamrizvi-cloud/urlx/discussions)
- **Author**: [@alhamrizvi-cloud](https://github.com/alhamrizvi-cloud)

---

**Created by [Alham Rizvi](https://github.com/alhamrizvi-cloud)**

**Happy Hunting! 🎯**
