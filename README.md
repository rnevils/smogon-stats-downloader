# smogon-stats-downloader

Scrapes https://www.smogon.com/stats/ and downloads usage stats as json.

## Requirements
- go 1.22

## Usage
```
go run main.go
```

## Output

`./data/` gets created if it doesn't exist. Sample from `./data/gen3ou-1760.json`:
```
[
  {
    "pokemon": "Tyranitar",
    "usage": 66.50269,
    "raw": 163556,
    "raw_pct": 43.917,
    "real": 138600,
    "real_pct": 44.707
  },
  {
    "pokemon": "Skarmory",
    "usage": 43.93202,
    "raw": 121810,
    "raw_pct": 32.708,
    "real": 112933,
    "real_pct": 36.428
  },
  {
    "pokemon": "Metagross",
    "usage": 42.28149,
    "raw": 131337,
    "raw_pct": 35.266,
    "real": 109642,
    "real_pct": 35.366
  }...
]
```
