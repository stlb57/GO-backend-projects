# Project 10 — Concurrent Web Crawler

A production-oriented Go concurrency project that starts from a URL, discovers internal links, and eventually crawls them concurrently while controlling duplicate work, cancellation, timeouts, and request rate.

The project is deliberately divided into phases so that each concept is introduced incrementally.

---

## Goal

Build a crawler that can:

1. Fetch web pages.
2. Parse their HTML.
3. Extract links.
4. Resolve and filter URLs.
5. Crawl multiple pages.
6. Avoid duplicate work.
7. Crawl pages concurrently using a worker pool.
8. Safely manage shared state.
9. Support cancellation and timeouts.
10. Control concurrency and request rate.
11. Measure performance.
12. Test the concurrent implementation for correctness and races.

---

# Phase 1 — HTTP Fundamentals ✅

Learn how Go communicates with a web server.

### Concepts

* `net/http`
* `http.Client`
* `http.NewRequest`
* `client.Do`
* `http.Response`
* HTTP status
* HTTP headers
* Response body
* `io.Reader`
* Reading/closing response bodies
* HTTP timeouts

### Output

Given a URL:

```text
URL
 ↓
HTTP request
 ↓
HTTP response
 ↓
status / headers / body
```

---

# Phase 2 — HTML + URL Processing ✅

Turn an HTTP response into useful URLs.

## HTML

Learn:

* `golang.org/x/net/html`
* HTML tokenizer/parser concept
* HTML documents as trees
* Nodes
* Element nodes
* Attributes
* `<a>` elements
* `href`
* Tree traversal

## URLs

Learn:

* `net/url`
* Parsing URLs
* Relative URLs
* Absolute URLs
* Resolving relative → absolute URLs
* Host/scheme checking
* Same-domain filtering

### Output

```text
https://example.com
        ↓
     fetch page
        ↓
    parse HTML
        ↓
   find <a href>
        ↓
   extract links
        ↓
 resolve URLs
        ↓
 keep same-domain URLs
```

**Phase 2 is complete.**

---

# Phase 3 — Sequential Crawler ← CURRENT

Turn the link extractor into an actual crawler.

## Concepts

* `[]string` as a queue
* Maps as a visited set
* BFS-style traversal
* Deduplication
* Maximum pages / crawl limits
* Handling cycles

### Basic architecture

```text
              Start URL
                  ↓
                Queue
                  ↓
              Take URL
                  ↓
                Fetch
                  ↓
            Extract links
                  ↓
          Check visited map
                  ↓
        Add new URLs to queue
                  ↑
                  └──── repeat
```

### Important

**No goroutines yet.**

The goal is simply:

> Starting from one URL, discover and crawl internal pages without repeatedly crawling the same page or getting stuck in cycles.

---

# Phase 4 — Concurrency

Convert the sequential crawler into a **Concurrent Web Crawler**.

## Goroutines

Multiple pages should be fetched simultaneously.

## Channels

Use channels to communicate jobs/results between goroutines.

## Worker Pool

Instead of creating unlimited goroutines:

```text
100 URLs
   ↓
100 goroutines
```

use a controlled number:

```text
              URL queue
                  ↓
        ┌────┬────┬────┬────┐
        ↓    ↓    ↓    ↓
       W1   W2   W3   W4
        ↓    ↓    ↓    ↓
      fetch fetch fetch fetch
```

The number of workers should be configurable.

## Concurrent State

Multiple workers will need to access shared state such as the visited set.

Example:

```text
Worker 1 ──┐
Worker 2 ──┼──→ visited
Worker 3 ──┤
Worker 4 ──┘
```

## Mutexes

Protect shared mutable state where necessary.

---

# Phase 5 — Bounded / Constrained Concurrency

Make concurrency configurable.

Example:

```bash
crawler -workers=8 -buffer=20
```

Experiment with:

* Worker count
* Channel buffer size
* Memory usage
* Throughput
* Backpressure
* Concurrency limits

### Main question

> Why doesn't increasing the number of goroutines always make the program faster?

Compare different concurrency levels and understand where the bottleneck actually is.

---

# Phase 6 — Context + Cancellation

Introduce:

```go
context.Context
```

The crawler should be able to receive a global:

> **STOP.**

Cancellation should propagate through the system:

```text
main
 ↓
context
 ↓
workers
 ↓
HTTP requests
```

## Concepts

* `context.Background()`
* `context.WithCancel()`
* `ctx.Done()`
* Cancellation propagation
* Cancelling HTTP requests

The crawler should stop starting new work once cancellation occurs.

---

# Phase 7 — Timeouts + Rate Limiting

## Timeouts

Prevent a slow or dead server from occupying a worker indefinitely.

```text
request
   ↓
waiting...
   ↓
timeout
   ↓
worker becomes available
```

Learn how HTTP timeouts interact with the crawler.

## Rate Limiting

Avoid sending too many requests too quickly.

The crawler should be able to control how aggressively it makes requests.

This introduces an important real-world networking concern:

> Being able to send requests concurrently doesn't mean you should send unlimited requests concurrently.

---

# Phase 8 — Benchmarking

Measure the crawler rather than assuming concurrency improves performance.

Test different worker counts:

```text
Workers    Time
----------------
1          ?
2          ?
4          ?
8          ?
16         ?
```

Measure:

* Total execution time
* Pages/second
* Effect of worker count
* Network bottlenecks
* Concurrency overhead
* Diminishing returns

### Goal

Understand experimentally:

> Where does additional concurrency stop providing meaningful performance improvements?

---

# Phase 9 — Testing + Race Detection

Test the crawler for correctness under concurrent execution.

Test cases should include:

* Duplicate URLs
* Cyclic links
* Broken links
* Empty pages
* Multiple workers discovering the same URL
* Cancellation
* Worker termination
* Request failures

Use Go's race detector:

```bash
go test -race
```

The goal is to detect unsafe concurrent access to shared state.

---

# Final Architecture

By the end, the crawler should roughly look like this:

```text
                         START URL
                             │
                             ▼
                      ┌────────────┐
                      │ URL Queue  │
                      └─────┬──────┘
                            │
                    ┌───────┴───────┐
                    ▼       ▼       ▼
                   W1      W2      W3 ... WN
                    │       │       │
                    ▼       ▼       ▼
                  HTTP    HTTP    HTTP
                    │       │       │
                    └───────┬───────┘
                            ▼
                       Parse HTML
                            │
                            ▼
                      Extract URLs
                            │
                            ▼
                     Normalize URLs
                            │
                            ▼
                      Deduplicate
                            │
                            ▼
                        URL Queue
```

Cross-cutting controls:

```text
Context
Timeouts
Rate limiting
Bounded workers
Channels
Mutexes
Concurrent state
```

And finally:

```text
Benchmarks
Tests
Race detector
```

---

# Project 10 Checklist

| Original Topic      | Phase   |
| ------------------- | ------- |
| Goroutines          | Phase 4 |
| Channels            | Phase 4 |
| Mutexes             | Phase 4 |
| HTTP client         | Phase 1 |
| Context             | Phase 6 |
| Cancellation        | Phase 6 |
| Timeouts            | Phase 7 |
| Deduplication       | Phase 3 |
| Rate limiting       | Phase 7 |
| Concurrent state    | Phase 4 |
| Benchmarking        | Phase 8 |
| Concurrency testing | Phase 9 |

---

## Current Progress

```text
Phase 1 — HTTP fundamentals       ✅
Phase 2 — HTML + URL processing   ✅
Phase 3 — Sequential crawler      🔨 NEXT
Phase 4 — Concurrency             ⏳
Phase 5 — Bounded concurrency     ⏳
Phase 6 — Context/cancellation    ⏳
Phase 7 — Timeouts/rate limiting  ⏳
Phase 8 — Benchmarking            ⏳
Phase 9 — Testing/race detection  ⏳
```

**Rule for this project:** don't skip ahead. Each phase exists to give us a reason to introduce the next concurrency concept.
