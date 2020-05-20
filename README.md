# webchat-basedon-pubsub
Web chat server and client based on PUBSUB

### RUN

```c
$ go get github.com/nkovacs/go-socket.io 
$ go run main.go
```

### Docker

Docker Build

```c
$ docker build -t webchat-golang:v1 .
```

### Result
- PUBSUB is easy to manage message sending and receiving

- Participant lists are listed on the top

- User name is updated automatically without duplication

![image](https://user-images.githubusercontent.com/52392004/82476740-764d2680-9b09-11ea-8b52-8a0de1359c01.png)

### Data Structure

#### Event 
```go
type Event struct {
	EvtType   string
	User      string
	Timestamp int
	Text      string
}
```

#### Subscription
```go
type Subscription struct {
	Archive []Event
	New     <-chan Event
}
```

#### Message Data Structure
```go
type Message struct {
	User      string
	Timestamp int
	Message   string
}
```


#### Chatroom Management
```go
func Chatroom() {
	archive := list.New()
	subscribers := list.New() // participants

	for {
		select {
		case c := <-subscribe:
			var events []Event

			// If there are archived events
			for e := archive.Front(); e != nil; e = e.Next() {
				events = append(events, e.Value.(Event))
			}

			subscriber := make(chan Event, 10)
			subscribers.PushBack(subscriber)

			c <- Subscription{events, subscriber}

		case event := <-publish:
			for e := subscribers.Front(); e != nil; e = e.Next() {
				subscriber := e.Value.(chan Event)
				subscriber <- event
			}

			// at least 5 events were stored
			if archive.Len() >= 5 {
				archive.Remove(archive.Front())
			}

			archive.PushBack(event)

		case c := <-unsubscribe:
			for e := subscribers.Front(); e != nil; e = e.Next() {
				subscriber := e.Value.(chan Event)

				if subscriber == c {
					subscribers.Remove(e)
					break
				}
			}
		}
	}
}
```

#### Troubleshooting - CORS
In order to excape CORS, the header of Access-Control-Allow-Origin was appended as bellow.

```go
    http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		// origin to excape Cross-Origin Resource Sharing (CORS)
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		server.ServeHTTP(w, r)
	})
```

### Reference

https://github.com/socketio/socket.io

https://github.com/iamshaunjp/websockets-playlist

https://github.com/nkovacs/go-socket.io

https://github.com/pyrasis/golangbook/blob/master/Unit%2067/chat.go
