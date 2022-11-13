# uploadfile-test


```
vegeta attack -targets=./files/target.json -duration=120s -rate=0 -max-workers=3 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         174220, 1451.80, 1451.79
Duration      [total, attack, wait]             2m0s, 2m0s, 1.181ms
Latencies     [min, mean, 50, 90, 95, 99, max]  230.92µs, 2.042ms, 1.348ms, 4.199ms, 5.968ms, 11.381ms, 68.868ms
Bytes In      [total, mean]                     5749260, 33.00
Bytes Out     [total, mean]                     595832400, 3420.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:174220 
```