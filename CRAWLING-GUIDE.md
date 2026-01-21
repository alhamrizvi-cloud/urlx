# urlX Active Crawling & File Extraction Guide

Complete guide for using urlX's active crawling, HTML/JS parsing, and file filtering capabilities.

## Table of Contents

- [Active Crawling](#active-crawling)
- [File Filtering](#file-filtering)
- [HTML Parsing](#html-parsing)
- [JavaScript Parsing](#javascript-parsing)
- [Workflow Examples](#workflow-examples)
- [Advanced Use Cases](#advanced-use-cases)

## Active Crawling

### Basic Crawling

```bash
# Crawl a single domain
urlx -domain example.com -crawl -v

# Crawl with custom depth
urlx -domain example.com -crawl -depth 3 -v

# Crawl URLs from file
urlx -u urls.txt -crawl -depth 2 -o crawled.txt
```

### What Gets Extracted

The crawler automatically extracts:

1. **Links** (`<a href="">`)
2. **JavaScript files** (`<script src="">` and inline scripts)
3. **CSS files** (`<link rel="stylesheet">`)
4. **Images** (`<img src="">`)
5. **Forms** (`<form action="">`)
6. **AJAX endpoints** (from JS code)
7. **API calls** (from fetch/axios in JS)

### Crawling Output

```bash
urlx -domain example.com -crawl -v

# Output:
# [*] Crawling [depth=0]: https://example.com
# [+] Found 45 URLs, 12 JS, 3 CSS from https://example.com
# [*] Crawling [depth=1]: https://example.com/about
# [+] Found 23 URLs, 5 JS, 2 CSS from https://example.com/about
```

## File Filtering

### Only Include Specific Extensions

Extract only specific file types:

```bash
# Only PDF files
urlx -domain example.com -crawl -only ".pdf" -o pdfs.txt

# Multiple extensions
urlx -domain example.com -crawl -only ".pdf,.doc,.docx" -o documents.txt

# JavaScript files only
urlx -d subdomains.txt -crawl -only ".js" -o javascript-files.txt

# Images only
urlx -u urls.txt -crawl -only ".jpg,.png,.gif,.svg" -o images.txt
```

### Remove Specific Extensions

Exclude unwanted file types:

```bash
# Remove images
urlx -domain example.com -crawl -remove ".png,.jpg,.gif" -o no-images.txt

# Remove media files
urlx -d domains.txt -crawl -remove ".mp4,.mp3,.avi,.mov" -o no-media.txt

# Remove common assets
urlx -u urls.txt -crawl -remove ".png,.jpg,.css,.woff,.ttf" -o endpoints-only.txt
```

### Combining Filters

```bash
# Get PDFs but exclude scanned images
urlx -domain example.com -crawl -only ".pdf" -remove ".png,.jpg"

# Get source code files, exclude minified
urlx -domain example.com -crawl -only ".js,.css,.html" -remove ".min.js,.min.css"
```

## HTML Parsing

### What the HTML Parser Extracts

The HTML parser uses `golang.org/x/net/html` to extract:

1. **All Links**
   ```html
   <a href="/about">About</a>
   <a href="https://api.example.com/v1">API</a>
   ```

2. **JavaScript Sources**
   ```html
   <script src="/assets/app.js"></script>
   <script src="https://cdn.example.com/lib.js"></script>
   ```

3. **CSS Files**
   ```html
   <link rel="stylesheet" href="/styles.css">
   ```

4. **Images**
   ```html
   <img src="/logo.png">
   <img src="https://cdn.example.com/banner.jpg">
   ```

5. **Forms**
   ```html
   <form action="/login" method="POST">
   <form action="https://api.example.com/submit">
   ```

6. **Inline JavaScript**
   ```html
   <script>
   fetch('/api/users').then(...)
   </script>
   ```

### Example: Extract All Links

```bash
# Crawl and extract all links
urlx -domain example.com -crawl -depth 2 -v

# Filter only internal links
urlx -domain example.com -crawl -v | grep "example.com"

# Find admin panels
urlx -domain example.com -crawl -v | grep -iE "admin|panel|dashboard"
```

## JavaScript Parsing

### What the JS Parser Extracts

The JavaScript parser uses regex patterns to find:

1. **URL Patterns**
   ```javascript
   "https://api.example.com/users"
   '/api/v1/posts'
   ```

2. **Fetch Calls**
   ```javascript
   fetch('/api/data')
   fetch("https://api.example.com/users")
   ```

3. **AJAX Requests**
   ```javascript
   $.ajax('/api/users')
   axios.get('/api/posts')
   ```

4. **URL Assignments**
   ```javascript
   src = "/assets/image.png"
   href = "https://example.com/page"
   ```

5. **CSS url() Functions**
   ```javascript
   url('/images/bg.png')
   url("https://cdn.example.com/font.woff")
   ```

### Example: Find API Endpoints

```bash
# Extract JavaScript files
urlx -domain example.com -crawl -only ".js" -o js-files.txt

# Crawl and find API endpoints
urlx -domain api.example.com -crawl -v | grep -E '/api/|/v[0-9]+'

# Find all fetch/ajax calls
urlx -domain example.com -crawl -depth 3 -v 2>&1 | grep -i "fetch\|ajax"
```

## Workflow Examples

### Workflow 1: Extract All PDFs from a Domain

```bash
#!/bin/bash
TARGET="example.com"

echo "[*] Discovering URLs..."
urlx -domain $TARGET -o all-urls.txt

echo "[*] Crawling for PDFs..."
urlx -u all-urls.txt -crawl -depth 2 -only ".pdf" -o pdfs.txt

echo "[*] Downloading PDFs..."
while read pdf; do
    wget "$pdf" -P downloads/
done < pdfs.txt

echo "[+] Found $(wc -l < pdfs.txt) PDF files"
```

### Workflow 2: JavaScript File Analysis

```bash
#!/bin/bash
TARGET="target.com"

echo "[*] Extracting all JavaScript files..."
urlx -domain $TARGET -crawl -depth 3 -only ".js" -o js-files.txt

echo "[*] Downloading JS files..."
mkdir -p js-downloads
while read jsfile; do
    filename=$(basename "$jsfile")
    curl -s "$jsfile" -o "js-downloads/$filename"
done < js-files.txt

echo "[*] Searching for sensitive data..."
grep -r "api_key\|password\|secret\|token" js-downloads/

echo "[*] Finding API endpoints..."
grep -rE "https?://[^\"']+" js-downloads/ | grep -E "/api/|/v[0-9]+"
```

### Workflow 3: Complete Asset Extraction

```bash
#!/bin/bash
TARGET="example.com"

# Extract different file types
urlx -domain $TARGET -crawl -only ".pdf" -o pdfs.txt
urlx -domain $TARGET -crawl -only ".doc,.docx" -o documents.txt
urlx -domain $TARGET -crawl -only ".xls,.xlsx" -o spreadsheets.txt
urlx -domain $TARGET -crawl -only ".js" -o javascript.txt
urlx -domain $TARGET -crawl -only ".css" -o stylesheets.txt
urlx -domain $TARGET -crawl -only ".json,.xml" -o data-files.txt

# Summary
echo "PDF files: $(wc -l < pdfs.txt)"
echo "Documents: $(wc -l < documents.txt)"
echo "Spreadsheets: $(wc -l < spreadsheets.txt)"
echo "JavaScript: $(wc -l < javascript.txt)"
echo "CSS: $(wc -l < stylesheets.txt)"
echo "Data files: $(wc -l < data-files.txt)"
```

### Workflow 4: Find Hidden Endpoints

```bash
#!/bin/bash
TARGET="target.com"

echo "[*] Stage 1: Passive discovery"
urlx -domain $TARGET -o passive-urls.txt

echo "[*] Stage 2: Active crawling"
urlx -u passive-urls.txt -crawl -depth 3 -remove ".png,.jpg,.gif,.css,.woff" -o crawled.txt

echo "[*] Stage 3: Extract interesting paths"
cat crawled.txt | unfurl paths | sort -u > paths.txt

echo "[*] Stage 4: Find admin/API endpoints"
grep -iE 'admin|api|v[0-9]|dashboard|panel|config|backup|test|dev|stage' paths.txt > interesting.txt

echo "[+] Found $(wc -l < interesting.txt) interesting endpoints"
```

### Workflow 5: Technology Stack Discovery

```bash
#!/bin/bash
TARGET="example.com"

echo "[*] Crawling for all assets..."
urlx -domain $TARGET -crawl -depth 2 -v -o all-assets.txt

echo "[*] JavaScript frameworks:"
cat all-assets.txt | grep -E "react|vue|angular|jquery" | head -10

echo "[*] CSS frameworks:"
cat all-assets.txt | grep -E "bootstrap|tailwind|bulma|foundation" | head -10

echo "[*] APIs found:"
cat all-assets.txt | grep -E "/api/|/graphql" | head -10

echo "[*] Backend technologies:"
urlx -domain $TARGET -probe -json | jq -r '.Technologies[]' | sort -u
```

## Advanced Use Cases

### Case 1: Extract Only Unminified JavaScript

```bash
# Get JS files but exclude minified versions
urlx -domain example.com -crawl -only ".js" -remove ".min.js" -o clean-js.txt

# Analyze for secrets
cat clean-js.txt | while read js; do
    echo "Checking: $js"
    curl -s "$js" | grep -E "api_key|secret|password|token"
done
```

### Case 2: Find Configuration Files

```bash
# Look for config files
urlx -domain example.com -crawl -depth 3 -only ".json,.xml,.yml,.yaml,.conf,.config,.env" -o configs.txt

# Probe if they're accessible
urlx -u configs.txt -probe -mc 200 -o accessible-configs.txt
```

### Case 3: Download All Documents

```bash
#!/bin/bash

# Extract all document types
urlx -domain example.com -crawl -only ".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx" -o documents.txt

# Create directory structure
mkdir -p downloads/{pdf,word,excel,powerpoint}

# Download with categorization
cat documents.txt | while read doc; do
    ext="${doc##*.}"
    case $ext in
        pdf)
            wget -q "$doc" -P downloads/pdf/
            ;;
        doc|docx)
            wget -q "$doc" -P downloads/word/
            ;;
        xls|xlsx)
            wget -q "$doc" -P downloads/excel/
            ;;
        ppt|pptx)
            wget -q "$doc" -P downloads/powerpoint/
            ;;
    esac
done

# Summary
find downloads/ -type f | wc -l
```

### Case 4: API Endpoint Discovery and Testing

```bash
#!/bin/bash
TARGET="api.example.com"

echo "[*] Discovering API endpoints..."
urlx -domain $TARGET -crawl -depth 3 -v | grep -E "/api/|/v[0-9]+" > api-endpoints.txt

echo "[*] Testing endpoints..."
urlx -u api-endpoints.txt -probe -mc 200,401,403 -json > api-results.json

echo "[*] Analyzing results..."
# Public endpoints (200)
jq -r 'select(.StatusCode == 200) | .URL' api-results.json > public-apis.txt

# Protected endpoints (401/403)
jq -r 'select(.StatusCode == 401 or .StatusCode == 403) | .URL' api-results.json > protected-apis.txt

echo "Public APIs: $(wc -l < public-apis.txt)"
echo "Protected APIs: $(wc -l < protected-apis.txt)"
```

### Case 5: Comprehensive Recon Pipeline

```bash
#!/bin/bash
TARGET="$1"

if [ -z "$TARGET" ]; then
    echo "Usage: $0 <domain>"
    exit 1
fi

echo "=== urlX Comprehensive Recon ==="
echo "Target: $TARGET"
echo ""

# Stage 1: Subdomain enumeration
echo "[1/6] Subdomain enumeration..."
subfinder -d $TARGET -silent > subs.txt
echo "  Found: $(wc -l < subs.txt) subdomains"

# Stage 2: Passive URL discovery
echo "[2/6] Passive URL discovery..."
urlx -d subs.txt -o passive-urls.txt
echo "  Found: $(wc -l < passive-urls.txt) URLs"

# Stage 3: Active crawling
echo "[3/6] Active crawling..."
urlx -u passive-urls.txt -crawl -depth 2 -remove ".png,.jpg,.gif,.svg,.woff,.ttf" -o crawled.txt
echo "  Found: $(wc -l < crawled.txt) URLs after crawling"

# Stage 4: Live probing
echo "[4/6] Live host probing..."
urlx -u crawled.txt -probe -mc 200,401,403 -o alive.txt
echo "  Alive: $(wc -l < alive.txt) hosts"

# Stage 5: File extraction
echo "[5/6] Extracting interesting files..."
urlx -u crawled.txt -only ".js" -o js-files.txt
urlx -u crawled.txt -only ".pdf,.doc,.docx" -o documents.txt
urlx -u crawled.txt -only ".json,.xml" -o data-files.txt
echo "  JS files: $(wc -l < js-files.txt)"
echo "  Documents: $(wc -l < documents.txt)"
echo "  Data files: $(wc -l < data-files.txt)"

# Stage 6: Summary
echo "[6/6] Generating report..."
cat << EOF > report.txt
=== Recon Report for $TARGET ===
Date: $(date)

Subdomains found: $(wc -l < subs.txt)
Total URLs discovered: $(wc -l < crawled.txt)
Live hosts: $(wc -l < alive.txt)

File breakdown:
- JavaScript: $(wc -l < js-files.txt)
- Documents: $(wc -l < documents.txt)
- Data files: $(wc -l < data-files.txt)

Files saved:
- subs.txt: All subdomains
- passive-urls.txt: Passive URL discovery
- crawled.txt: All crawled URLs
- alive.txt: Live hosts
- js-files.txt: JavaScript files
- documents.txt: Document files
- data-files.txt: JSON/XML files
EOF

cat report.txt
echo ""
echo "[+] Recon complete! Check report.txt for summary."
```

## Tips and Best Practices

### Performance Optimization

```bash
# Adjust concurrency based on target
urlx -domain example.com -crawl -c 10  # Default
urlx -domain example.com -crawl -c 50  # Aggressive
urlx -domain example.com -crawl -c 3   # Gentle

# Set appropriate depth
urlx -domain example.com -crawl -depth 1  # Fast, surface level
urlx -domain example.com -crawl -depth 3  # Thorough
urlx -domain example.com -crawl -depth 5  # Very deep (slow)
```

### Filtering Best Practices

```bash
# For bug bounties - focus on code/config
urlx -domain target.com -crawl -only ".js,.json,.xml,.yml,.yaml" -remove ".min.js"

# For data gathering - documents only
urlx -domain target.com -crawl -only ".pdf,.doc,.docx,.xls,.xlsx"

# For asset discovery - no media
urlx -domain target.com -crawl -remove ".png,.jpg,.gif,.mp4,.mp3,.svg,.woff"
```

### Domain-Specific Crawling

```bash
# Crawl only within target domain
urlx -domain example.com -crawl -depth 3 -v

# The crawler automatically filters to stay within the domain
# This prevents crawling external resources
```

## Common Patterns

### Extract File Type

```bash
# Images
urlx -domain example.com -crawl -only ".jpg,.jpeg,.png,.gif,.svg,.webp"

# Documents  
urlx -domain example.com -crawl -only ".pdf,.doc,.docx,.odt"

# Spreadsheets
urlx -domain example.com -crawl -only ".xls,.xlsx,.csv"

# Archives
urlx -domain example.com -crawl -only ".zip,.rar,.tar,.gz,.7z"

# Code
urlx -domain example.com -crawl -only ".js,.css,.html,.php,.py,.rb"

# Data
urlx -domain example.com -crawl -only ".json,.xml,.yml,.yaml"

# Config
urlx -domain example.com -crawl -only ".conf,.config,.ini,.env"
```

## Troubleshooting

### Slow Crawling

```bash
# Reduce depth
urlx -domain example.com -crawl -depth 1

# Increase concurrency
urlx -domain example.com -crawl -c 20

# Filter out media files
urlx -domain example.com -crawl -remove ".png,.jpg,.mp4"
```

### Too Many Results

```bash
# Use stricter filters
urlx -domain example.com -crawl -only ".js,.json"

# Reduce depth
urlx -domain example.com -crawl -depth 1

# Filter output
urlx -domain example.com -crawl | grep -v "cdn\|static"
```

### Memory Issues

```bash
# Process in batches
split -l 100 urls.txt batch_
for batch in batch_*; do
    urlx -u $batch -crawl -o "output-$batch.txt"
done
```

---

For more examples, see [README.md](README.md) and [PROBING-GUIDE.md](PROBING-GUIDE.md).
