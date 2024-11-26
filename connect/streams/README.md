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
    A[UDP] -->|Batch 3s| B(Send Batch)
    B --> C{Fan Out}
    C --> H{Round Robin}
    C --> S{Round Robin}
    H -->|Sqlite DB| D[Cluster 1]
    H -->|Sqlite DB| J[Cluster 1]
    S -->|Sqlite DB| E[Cluster 2]
    S -->|Sqlite DB| L[Cluster 2]
```
