# go-mysql-queue
A Very Basic Queue / Job implementation which uses MySQL for underlying storage

## Configuration

- Setup a MySQL database or a path for an SQLite database.
- 

## Example Usage

```
import (
    "fmt"
    "time"

    msq "github.com/AnalogRepublic/go-mysql-queue"
)

// Connect to our backend database
queue, err := msq.Connect(msq.ConnectionConfig{
    Type: "mysql", // or could use "sqlite"
    Host: "localhost",
    Username: "root",
    Password: "root",
    Database: "queue",
})

if err != nil {
    panic(err)
}

queue.Configure(&msq.QueueConfig{
    MaxRetries: 3, // The maximum number of times the message can be retried.
    MessageTTL: 5 * time.Second, // the longest time they'll live in the queue before being added to the dead-letter table.
    Name: "my-queue", // The namespace for the Queue
})

if err != nil {
    panic(err)
}

// Using a listener
listener := &Listener{
    Queue:  *queue,
    Config: listenerConfig,
}

ctx, _ := listener.Context()

listener.Start(func(event Event) bool {
    assert.Equal(t, queuedEvent.UID, event.UID)
    return true
})

fmt.Println("Listener started")

select {
case <-ctx.Done():
    fmt.Println("Listener stopped")
}

// or manually pull an item off the queue
event, err := queue.Pop()

if err == nil {
    err := doSomethingWithMessage(event)

    // If we have an error we can requeue it
    if err != nil {
        queue.ReQueue(event)
    }

    // or say we're happy with it
    queue.Done(event)
}

time.AfterFunc(5 * time.Second, func() {

    // Push a new item onto the Queue
    queue.Push(msq.Payload{
        "example": "data",
        "testing": []string{
            "a", 
            "b",
        },
        "oh-look": map[string]string{
            "maps": "here",
        },
    })
})

```
