# Connect Pipeline

Store and forward standard processing using Redpanda Connect

1. Capture UDP 514 syslog packets in a buffer
2. Emit batches from buffer to both active backends
3. One each producer,
   1. round robin between N producers
   2. use a SQLlite db to store batches
   3. transmit to Redpanda services
   4. Remove from database once batch succeeds

```mermaid
flowchart TD
    A[Generator] -->|1 second| C[Socket]
    C --> H[UDP] -->|Batch - 3 seconds| K(Send Batch)
    K --> M{Fan Out}
    M --> N{Round Robin}
    M --> S{Round Robin}
    N -->|Sqlite DB| D[Cluster 1]
    N -->|Sqlite DB| J[Cluster 1]
    S -->|Sqlite DB| E[Cluster 2]
    S -->|Sqlite DB| L[Cluster 2]
```
