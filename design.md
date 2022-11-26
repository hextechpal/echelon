
# Design
This document proposes a new design for the prio based on gossip protocol and consistent hashing

## Prio

Prio is a distributed fault-tolerant scalable priority queueing service

### Concepts

- `topics` : `topics` divides all the work into logical partitions. Priorities are honored with in a topics
```go
type Topic struct {
	Name        string         
	Description string 
	CreatedAt   int64          
	UpdatedAt   int64          
}
```

- `jobs` : `jobs` are logical units of work. Each job is associated with a `topic`

```go
type Status int

const (
    PENDING   Status = iota // Represents the job is just inserted in to the database
    CLAIMED                 // The job has been claimed by the consumer
    COMPLETED               // The job has been marked completed by the consumer
)

type Job struct {
    ID       int64  
    Topic    string 
    Payload  []byte 
    Priority int32  
    Status   Status 
    
    ClaimedAt int64          
    ClaimedBy string 
    
    CompletedAt int64 
    
    CreatedAt int64 
    UpdatedAt int64 
}

```

### API

```go
type API interface{
    // GetTopics :Creates a new topic
    GetTopics(ctx context.Context) ([]string, error)
    
    // RegisterTopic :Creates a new topic
    RegisterTopic(ctx context.Context, req RegisterTopicRequest) (RegisterTopicResponse, error)
    
    // Enqueue : Persist a job in to the storage engine
    Enqueue(ctx context.Context, req EnqueueRequest) (EnqueueResponse, error)
    
    // Dequeue : Picks the top priority job from the given topic and returns empty if no job is present
    Dequeue(ctx context.Context, req DequeueRequest) (DequeueResponse, error)
    
    // Ack : Acknowledge a claimed jon, job is marked as complete once the acknowledgement is received
    // If consumer do not ack the jobId after a fixed amount of time (10sec) for mysql engine
    // the job will be moved back to pending state and is available to deque again based on priority
    Ack(ctx context.Context, request AckRequest) (AckResponse, error)
    
    // ReQueue : Requeue operation runs periodically and move the unacked jobs to the pending queue again
    ReQueue(ctx context.Context, req RequeueRequest) (RequeueResponse, error)
}
```

### Architecture
Prio will consist a set of workers which will enqueue, dequeue and ack requests.

- **Partitioning** : The topics are partitioned among the workers based on consistent hashing. By default the workers start with a configuration of 256 partitions 

- **Discovery** : Every worker discover all the other workers based on a gossip protocol

- **Routing** : Requests can be routed to any worker in the cluster. If the request can be served from thr worker it will get served else it is forwarded to the appropriate worker which can serve the request   


