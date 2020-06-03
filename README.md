# webchat-golang-socket.io
It shows a chat application in which socket.io is used to connect server and client.
Also, there is no difference for 1-to-1 and groupchat since it is based on channel communication.
In this project, PUBSUB structure was deployed to support Say, Join and Leave.

The server has a strangth to support massive traffics where the number of members is huge in a chatroom.
It is different with mobile text application since delivery and display notifications are not required since all members basically show a chatroom together as well as Slack.

### RUN

```c
$ go get github.com/nkovacs/go-socket.io 
$ go run main.go
```

Docker Build

```c
$ docker build -t webchat-golang:v1 .
```

### Result
- socket.io provide stable connection between server and client

- Participant lists are listed on the top

- User name is updated automatically without duplication

![image](https://user-images.githubusercontent.com/52392004/82513003-b255ab00-9b4c-11ea-8ef0-5f22cf872c11.png)


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

#### Channel Management
```go
server.On("connection", func(so socketio.Socket) {
		log.D("connected... %v", so.Id())

		newMessages := make(chan string)

		s := Subscribe()

		so.On("join", func(user string) {
			log.D("Join...%v (%v)", user, so.Id())

			Join(user) // Join notification
			userMap[so.Id()] = user

			// if there are archived events
			for _, event := range s.Archive {
				so.Emit("chat", event)
			}
		})

		so.On("chat", func(msg string) {
			newMessages <- msg
		})

		so.On("disconnection", func() {
			log.D("disconnected... %v", so.Id())

			user := userMap[so.Id()]
			delete(userMap, so.Id())

			Leave(user) // left notifcation
			s.Cancel()

			// update participant lists
			str := getParticipantList(userMap)
			userStr = str
			log.D("Update Participantlist: %v", userStr)
			so.Emit("participant", userStr)
		})

		go func() {
			for {
				select {
				case event := <-s.New: // send event to browser
					so.Emit("chat", event)

					// update participant lists
					if event.EvtType == "join" || event.EvtType == "leave" {
						str := getParticipantList(userMap)
						userStr = str
						log.D("Update Participantlist: %v", userStr)
						so.Emit("participant", userStr)
					}

				case msg := <-newMessages: // received message from browser
					var newMSG Message
					json.Unmarshal([]byte(msg), &newMSG)

					Say(newMSG)
				}
			}
		}()
	})
```

#### Join
```go
func Join(user string) {
	timestamp := time.Now().Unix()
	publish <- NewEvent("join", user, int(timestamp), "")
}
```

#### Say
```go
func Join(user string) {
	timestamp := time.Now().Unix()
	publish <- NewEvent("join", user, int(timestamp), "")
}
```

#### Leave
```go
func Leave(user string) {
	timestamp := time.Now().Unix()
	publish <- NewEvent("leave", user, int(timestamp), "")
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
