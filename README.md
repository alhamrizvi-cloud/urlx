<img width="748" height="310" alt="image" src="https://github.com/user-attachments/assets/08dc9c98-5a81-439d-8b94-3ba7d85fa7a1" />

# urlX — Multi‑Source URL Discovery Tool

[![Go Version](https://img.shields.io/badge/Go-1.19%2B-00ADD8?style=flat\&logo=go)](https://golang.org)
[![Stars](https://img.shields.io/github/stars/alhamrizvi-cloud/urlx?style=flat)](https://github.com/alhamrizvi-cloud/urlx/stargazers)
[![Forks](https://img.shields.io/github/forks/alhamrizvi-cloud/urlx?style=flat)](https://github.com/alhamrizvi-cloud/urlx/network)
[![Issues](https://img.shields.io/github/issues/alhamrizvi-cloud/urlx)](https://github.com/alhamrizvi-cloud/urlx/issues)
[![Last Commit](https://img.shields.io/github/last-commit/alhamrizvi-cloud/urlx)](https://github.com/alhamrizvi-cloud/urlx/commits/main)
[![Code Size](https://img.shields.io/github/languages/code-size/alhamrizvi-cloud/urlx)](https://github.com/alhamrizvi-cloud/urlx)
[![Top Language](https://img.shields.io/github/languages/top/alhamrizvi-cloud/urlx)](https://github.com/alhamrizvi-cloud/urlx)

## 🔍 Overview (Resume)

**urlX** is a high-performance reconnaissance tool designed for **bug bounty hunters**, **penetration testers**, and **security researchers**.

It combines:

* Passive URL discovery from **11+ intelligence sources**
* **Live host probing** (DNS, TCP, TLS, HTTP)
* **Active web crawling** for hidden endpoints
* Smart **file & extension filtering**
* Fast, concurrent processing using **Go routines**

urlX helps you quickly discover **attack surface**, **hidden assets**, and **live endpoints** with minimal setup.

> Built for speed, scale, and real-world recon workflows.

Created by **[Alham Rizvi](https://github.com/alhamrizvi-cloud)**

---

## 🚀 Installation Guide

### Requirements

* **Go 1.19 or higher**
* Linux / macOS (Windows supported with Go environment)
* Internet access for passive sources

Verify Go installation:

```bash
go version
```

---

### 🔹 Method 1: Go Install (Recommended)

```bash
go install github.com/alhamrizvi-cloud/urlx@latest
```

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verify:

```bash
urlx -h
```

---

### 🔹 Method 2: Build From Source

```bash
git clone https://github.com/alhamrizvi-cloud/urlx.git
cd urlx
go build -o urlx main.go
```

Move binary system-wide (optional):

```bash
sudo mv urlx /usr/local/bin/
```

Verify:

```bash
urlx -h
```

---

### 🔹 Method 3: Quick Install Script

```bash
curl -sSfL https://raw.githubusercontent.com/alhamrizvi-cloud/urlx/main/install.sh | sh
```

---

### 🔹 Method 4: Docker

```bash
docker build -t urlx .
docker run --rm urlx -h
```

---

## 🔑 Optional API Keys (Enhance Results)

urlX works **without API keys**, but adding them increases coverage and accuracy.

Supported providers:

* URLScan.io
* AlienVault OTX
* VirusTotal
* GitHub
* Censys
* SecurityTrails

Keys can be passed via CLI flags or environment variables.

---

## ⚠️ Disclaimer

This tool is intended **only for authorized security testing and legal research**.

You are responsible for ensuring you have **explicit permission** before scanning any target.

**Misuse may be illegal.**

---

## 📄 License

MIT License – see [LICENSE](LICENSE)

---

## 🤝 Contributing

Pull requests, feature ideas, and source additions are welcome.
See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## ⭐ Support the Project

If you find **urlX** useful, please consider giving it a ⭐ on GitHub.

👉 [https://github.com/alhamrizvi-cloud/urlx](https://github.com/alhamrizvi-cloud/urlx)

---

**Built with ❤️ by Alham Rizvi**
**Happy Hunting 🎯**
