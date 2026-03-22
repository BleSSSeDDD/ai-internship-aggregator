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

### 1. Spin up the infrastructure
We use a **Monorepo** approach. Everything you need is in the root folder.
```bash
task infra


TODO:

1) Кафка не работат почему-то пока ---> 
---> 1.1) теперь работает, но ui не парсит протобаф нормально, хуй знает,
полный рандом тут, надо просто пробовать
(короче тебе уже можно писать код, считай, кафка работает, на ui пофиг)
-------------------------------------------------------- ГОТОВО
2) AI слой пока мокается двумя сообщениями (уже нет)
-------------------------------------------------------- ГОТОВО
3) Caddy серт самоподписанный генерит, браузер на это ругается, в принципе пофиг
-------------------------------------------------------- ГОТОВО


4) .env пока просто существует


5) дописать в гитигнор шляпы для жавы
-------------------------------------------------------- ГОТОВО
6) Написать бэк + фронт))0))
-------------------------------------------------------- ГОТОВО
7) придумать схему бд чтоб нестыдно было, полей дофига, возможно,
стоит жсоны хранить, тогда можно не посгрес а nosql бдшку какую-то
-------------------------------------------------------- ГОТОВО


8) настроить кэдди


9) этапы надо сделать строкой, а не массивом строк емае
-------------------------------------------------------- ГОТОВО


10) если на одной странице несколько стажировок, все ломается
---> 10.1) частично пофикшено, но всё еще неиделаьно, промпт инжиниринг это реально сложно


11) нейронка галюцинирует ссылки, надо ей дать их шфывопгтгшфптг мем
-------------------------------------------------------- ГОТОВО

12) время внутри контейнеров настроить

13) нужна админ-панель, метрики, возможность добавлять вручную карточки
---> 13.1) чтоб вручную добавлять я могу написать микросервис, в котором простенький фронт там анкету
заполняешь и он отправляет это в кафку, нооо хз насколько это правильно вообще с учетом того,
что там ttl в бд стоит, а парсер не поймает эту карточку по-новой, раз раньше её не ловил

вывод: надо думать головой

