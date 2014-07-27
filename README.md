# ljstats #

ljstats is a utility program I wrote to parse Lockjaw game data into a more usable format. It turns lj-scores.txt into CSV-formatted data. It can also be used to create heatmaps from game data.

## Usage ##

```bash
> go build stats.go heatmap.go
> ./stats INPUT [OUTPUT]
```
OUTPUT must be a .csv or .png file.
If OUTPUT is not specified, csv data is written to stdout.

## Notes ##

There is a lot of variation possible when creating a heatmap -- too many variables to specify on the command line. Instead, you can generate different heatmaps by editing the source directly. In the future, I may add support for specifying at least the X and Y values on the command line.

For now, this program only cares about one type of game. It throws out any game that:
- Was not played in 40-line mode
- Was not completed
- Was not played using the default bag of pieces

These settings could easily be adjusted, but I have little reason to do so unless someone specifically asks.

Also note that game data is normalized after parsing in order to remove outliers. You'll probably want to tweak the normalization parameters, since they are currently tailored to my own stats.
