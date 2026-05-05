# logslice

A fast log file slicer and filter tool with structured output support.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

---

## Usage

```bash
# Slice logs between two timestamps
logslice --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00" app.log

# Filter by keyword and output as JSON
logslice --filter "ERROR" --format json app.log

# Read from stdin and slice by line range
cat app.log | logslice --lines 100:500

# Combine filters with structured output
logslice --from "2024-01-15 08:00:00" --filter "WARN" --format json app.log > warnings.json
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start timestamp for slicing |
| `--to` | End timestamp for slicing |
| `--filter` | Keyword or regex to filter lines |
| `--lines` | Line range in `start:end` format |
| `--format` | Output format: `text` (default) or `json` |

---

## Features

- Fast line-by-line streaming with minimal memory usage
- Timestamp-aware slicing for common log formats
- Regex and plain-text filtering
- Structured JSON output for downstream processing
- Supports gzipped log files (`.log.gz`)

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)