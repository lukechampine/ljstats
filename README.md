# ljstats #

ljstats is a utility program I wrote to parse Lockjaw game data into a more usable format. It turns lj-scores.txt into CSV-formatted data. It can also be used to create heatmaps from game data.

## Usage #

```bash
> go build stats.go heatmap.go
> ./stats INPUT [OUTPUT]
```
OUTPUT must be a .csv or .png file.
If OUTPUT is not specified, csv data is written to stdout.
