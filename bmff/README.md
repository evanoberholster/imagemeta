# BMFF Package

## Description

Golang library for parsing ISOBMFF Image files specifically Heif/Heic, AV1/AVIF and CR3.

### Parsing File Type Support

- Heif/Heic (Partial)
- AV1/Avif (WIP)
- CR3 (WIP)

### Usage 

To be continued...

### Benchmarks

Benchmark for reading a Heif Metabox and extracting Exif:

```
name                   old time/op    new time/op    delta
HeicExif100/1             202µs ± 4%      29µs ±17%  -85.67%  (p=0.000 n=8+20)
HeicExif100/2            21.8µs ±11%    10.1µs ± 8%  -53.88%  (p=0.000 n=10+19)
HeicExif100/3             301µs ± 7%      40µs ± 9%  -86.72%  (p=0.000 n=10+18)
HeicExif100/10            215µs ± 8%      31µs ±15%  -85.66%  (p=0.000 n=10+19)
HeicExif100/d            1.02ms ± 4%    0.12ms ± 5%  -88.60%  (p=0.000 n=9+19)
HeicExif100/Canon_R6     83.9µs ±11%    11.5µs ± 7%  -86.24%  (p=0.000 n=10+19)
HeicExif100/iPhone_12     207µs ± 7%      28µs ± 9%  -86.37%  (p=0.000 n=10+18)

name                   old alloc/op   new alloc/op   delta
HeicExif100/1             479kB ± 0%       7kB ± 0%  -98.49%  (p=0.000 n=10+20)
HeicExif100/2            41.2kB ± 0%     1.7kB ± 0%  -95.77%  (p=0.000 n=10+20)
HeicExif100/3             737kB ± 0%      12kB ± 0%  -98.35%  (p=0.000 n=10+20)
HeicExif100/10            473kB ± 0%       7kB ± 0%  -98.45%  (p=0.000 n=10+20)
HeicExif100/d            2.22MB ± 0%    0.04MB ± 0%  -98.25%  (p=0.000 n=8+20)
HeicExif100/Canon_R6      178kB ± 0%       2kB ± 0%  -98.90%  (p=0.000 n=10+20)
HeicExif100/iPhone_12     481kB ± 0%       8kB ± 0%  -98.41%  (p=0.000 n=10+20)

name                   old allocs/op  new allocs/op  delta
HeicExif100/1             1.23k ± 0%     0.06k ± 0%  -94.79%  (p=0.000 n=10+20)
HeicExif100/2               116 ± 0%        29 ± 0%  -75.00%  (p=0.000 n=10+20)
HeicExif100/3             1.92k ± 0%     0.10k ± 0%  -94.58%  (p=0.000 n=10+20)
HeicExif100/10            1.29k ± 0%     0.07k ± 0%  -94.74%  (p=0.000 n=10+20)
HeicExif100/d             5.87k ± 0%     0.31k ± 0%  -94.68%  (p=0.000 n=10+20)
HeicExif100/Canon_R6        549 ± 0%        33 ± 0%  -93.99%  (p=0.000 n=10+20)
HeicExif100/iPhone_12     1.33k ± 0%     0.07k ± 0%  -94.59%  (p=0.000 n=10+20)
```

### Testing

To be continued...

### Special Thanks to:

- Laurent Clévy (@Lorenzo2472) (https://github.com/lclevy/canon_cr3) for Canon CR3 structure
- The go4 Authors (https://github.com/go4org/go4) for their work on BMFF parser and HEIF structure
- Lasse Heikkilä (https://trepo.tuni.fi/bitstream/handle/123456789/24147/heikkila.pdf?sequence=3) for his thesis for ideas for the parser.

