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

// Make sure we have the table & schema that we need.
err = queue.SetupDatabase()

if err != nil {
    panic(err)
}

// Setup an automatic listener.
_, err := queue.Listen(func(message msq.Message) {
    payload := message.Payload
    
    fmt.Println(payload["example"].(string))

    err := library.DoSomething(payload["example"].(string))

    if err != nil {
        // if we have an error, tell it to
        // be-requeued.
        return false
    }
    
    // Otherwise remove it.
    return true
}, msq.ListenerConfig{
    // How often we're fetching a new message
    Interval: time.Millisecond,
    
    // The number of messages to fetch at a time
    BatchSize: 1,

    // After 10 seconds we'll cut off the processing of a message.
    // and mark it to be requeued.
    Timeout: 10 * time.Second,
})

if err != nil {
    panic(err)
}

// or manually pull an item off the queue
message, err := queue.Pop()

if err == nil {
    err := doSomethingWithMessage(message)

    // If we have an error we can requeue it
    if err != nil {
        queue.ReQueue(message)
    }

    // or say we're happy with it
    queue.Done(message)
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