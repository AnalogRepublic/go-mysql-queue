# go-mysql-queue ![https://travis-ci.org/AnalogRepublic/go-mysql-queue.svg?branch=master](https://travis-ci.org/AnalogRepublic/go-mysql-queue)
A Very Basic Queue / Job implementation which uses MySQL for underlying storage

## Example Usage

```
import (
    "fmt"
    "time"

    msq "github.com/AnalogRepublic/go-mysql-queue"
)

// Connect to our backend database
queue, err := msq.Connect(msq.ConnectionConfig{
    Type: "mysql", // or could use "sqlite", where the "database" field is the full path, e.g. "./test.db"
    Host: "localhost",
    Username: "root",
    Password: "root",
    Database: "queue",
})

if err != nil {
    panic(err)
}

queue.Configure(&msq.QueueConfig{
    Name: "my-queue", // The namespace for the Queue
    MaxRetries: 3, // The maximum number of times the message can be retried.
})

if err != nil {
    panic(err)
}

// Using a listener
listener := &Listener{
    Queue:  *queue,
    Config: listenerConfig,
}

ctx := listener.Context()

listener.Start(func(event Event) bool {
    fmt.Println("Received event " + event.UID)
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
    } else {
        // or say we're happy with it
        queue.Done(event)
    }
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
