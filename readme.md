# gobrot
Rewritten [Mandelbrot](https://github.com/rikcarve/mandelbrot) using GO.

Mainly for learning GO and a little bit to benchmark Java vs GO.

## UI
Used [Gio UI](https://gioui.org/) for no specific reason (just needed something working on Windows)

## Benchmark
5 runs 800x600 with zooming:

|run|Java|GO
|---|----|--
|first|2.4s|2.2s
|next 4|2.0s|2.1s

As more or less expected, Java is a bit slower in the first run. Overall no real difference...

