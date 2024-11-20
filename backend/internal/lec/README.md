# LEC - Layered Event Counter

LEC is compute & storage efficient data structure for creating real time histograms for time series which the cutoff between intervals shift to match presentation time. LEC compansates its disadvantages as the ratio of interval period/refresh duration grows. It is designed for the case 1 hour interval and 5 seconds refresh (720x).

**How it works**

LEV aggregates data in multiple layers, each layer store the sum of events happened in longer period of time. This reduces the number of additions needed to perform to refresh presentation.

**Example**

Giving that the length is 1 day and the minimum resolution is 1 second:

Note: `86400 = 60 * 60 * 24`

| layer | number of virtual cells | cell resolution      |
| ----- | ----------------------- | -------------------- |
| 17th  | `ceil(86400 / 65536)`   | 65536s (18h 12m 16s) |
| 16th  | `ceil(86400 / 32768)`   | 32768s (9h 6m 8s)    |
| 15th  | `ceil(86400 / 16384)`   | 16384s (4h 33m 4s)   |
| 14th  | `ceil(86400 / 8192)`    | 8192s (2h 16m 32s)   |
| 13th  | `ceil(86400 / 4096)`    | 4096s (1h 8m 16s)    |
| 12th  | `ceil(86400 / 2048)`    | 2048s (34m 8s)       |
| 11th  | `ceil(86400 / 1024)`    | 1024s (17m 4s)       |
| 10th  | `ceil(86400 / 512)`     | 512s (8m 32s)        |
| 9th   | `ceil(86400 / 256)`     | 256s (4m 16s)        |
| 8th   | `ceil(86400 / 128)`     | 128s (2m 8s)         |
| 7th   | `ceil(86400 / 64)`      | 64s (1m 4s)          |
| 6th   | `ceil(86400 / 32)`      | 32s                  |
| 5th   | `ceil(86400 / 16)`      | 16s                  |
| 4th   | `ceil(86400 / 8)`       | 8s                   |
| 3rd   | `ceil(86400 / 4)`       | 4s                   |
| 2nd   | `ceil(86400 / 2)`       | 2s                   |
| 1st   | `ceil(86400 / 1)`       | 1s                   |

**Storage optimization**

LEC only stores the second in each two cells in each layer, since the value of other cell can be calculated by subtracting the second one from the cell in layer above.
