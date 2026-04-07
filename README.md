# ⚡️ AI-INTERNSHIP-AGGREGATOR

> **"Turning messy HTML into gold. AI-powered internship engine fueled by Go, Kafka, and pure enthusiasm. Overengineered? Maybe. Overpowered? Definitely. The cleanest AI-parser in the game."**

---

## 🧬 The Concept
This is not your typical regex-based scraper. This is a high-performance **distributed system** designed to solve the chaos of internship postings. We use LLMs to give structure to the unstructured web, delivering clean, strongly-typed data via **Protobuf** and **Kafka**.



## 🛠 Tech Stack & Architecture

| Layer | Technology | Purpose |
| :--- | :--- | :--- |
| **Scraper** | `Go 1.25+` | High-concurrency crawling & data extraction |
| **Brain** | `Ollama (Llama 3)` | Local LLM for semantic parsing of job descriptions |
| **Contract** | `Protobuf` | Strict, language-agnostic data schemas |
| **Transport** | `Apache Kafka` | Event-driven streaming between services |
| **Dashboard** | `Java 21 / Spring Boot` | Solid backend for data aggregation & UI |

### 🏗 Architecture: Hexagonal / Clean
The Go service is built with **Separation of Concerns** in mind:
- `internal/domain`: Pure business logic & interfaces.
- `internal/usecase`: Application orchestrators.
- `internal/infrastructure`: Swappable adapters (Kafka, Ollama, Colly).

---

## 🚀 Quick Start

### Setup Guide

1. Install Task
Follow the official instructions:
https://taskfile.dev/installation

Clone repository

```bash
git clone https://github.com/BleSSSeDDD/ai-internship-aggregator.git
cd ai-internship-aggregator
```
2. Run the Task command

```bash
task up-with-model
```
## Browser view

admin panel for manually adding internship cards:

<img width="1802" height="977" alt="image" src="https://github.com/user-attachments/assets/8e6e3d6c-2f11-41fd-b62a-5eb35c3e4498" />



