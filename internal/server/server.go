package server

import (
	"container/list"
	"encoding/json"
	"net/http"
	"time"
	"webchat-basedon-pubsub/internal/config"
	"webchat-basedon-pubsub/internal/logger"

	socketio "github.com/nkovacs/go-socket.io"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("server")
}

var (
	subscribe   = make(chan (chan<- Subscription), 10)
	unsubscribe = make(chan (<-chan Event), 10)
	publish     = make(chan Event, 10)
)

// Event is to define the event
type Event struct {
	EvtType   string
	User      string
	Timestamp int
	Text      string
}

// Subscription is to manage subscribe events
type Subscription struct {
	Archive []Event
	New     <-chan Event
}

// Message is the data structure of messages
type Message struct {
	User      string
	Timestamp int
	Message   string
}

// InitServer initializes the server
func InitServer(conf *config.AppConfig) error {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.E("%v", err)
	}

	go Chatroom()

	userMap := make(map[string]string) // hashmap to memorize the pair of socket id and user id

	server.On("connection", func(so socketio.Socket) {
		log.D("connected... %v", so.Id())

		newMessages := make(chan string)

		s := Subscribe()

		so.On("join", func(user string) {
			log.D("Join...%v (%v)", user, so.Id())

			Join(user) // Join notification
			userMap[so.Id()] = user

			// if there are archived events
			//	for _, event := range s.Archive {
			//	log.D("archived event: %v %v %v %v", event.EvtType, event.User, event.Timestamp, event.Text)
			//	so.Emit("chat", event)
			// }
		})

		so.On("chat", func(msg string) {
			newMessages <- msg
		})

		so.On("disconnection", func() {
			log.D("disconnected... %v", so.Id())

			user := userMap[so.Id()]

			Leave(user) // left notifcation
			s.Cancel()
		})

		go func() {
			for {
				select {
				case event := <-s.New: // send event to browser
					log.D("sending event to browsers: %v %v %v %v (%v)", event.EvtType, event.User, event.Timestamp, event.Text, so.Id())
					so.Emit("chat", event)

				case msg := <-newMessages: // received message from browser
					var newMSG Message
					json.Unmarshal([]byte(msg), &newMSG)

					log.D("receiving message from browser: %v %v %v (%v)", newMSG.User, newMSG.Timestamp, newMSG.Message, so.Id())
					Say(newMSG)
				}
			}
		}()
	})

	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		// origin to excape Cross-Origin Resource Sharing (CORS)
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// address
		r.RemoteAddr = "10.253.69.155"
		log.I("Address: %v", r.RemoteAddr)

		server.ServeHTTP(w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("./asset")))

	log.I("Serving at %v:%v", conf.ChatInfo.Host, conf.ChatInfo.Port)
	//port := ":" + conf.ChatInfo.Port
	log.E("%v", http.ListenAndServe(":4000", nil))

	return err
}

// Chatroom is to manage all events in a chatroom
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

// NewEvent is to create an new event
func NewEvent(evtType string, user string, timestamp int, msg string) Event {
	return Event{evtType, user, timestamp, msg}
}

// Subscribe is to add a subscriber
func Subscribe() Subscription {
	c := make(chan Subscription)
	subscribe <- c
	return <-c
}

// Join is to make a join event
func Join(user string) {
	timestamp := time.Now().Unix()
	publish <- NewEvent("join", user, int(timestamp), "")
}

// Say is to define the event of chatting
func Say(msg Message) {
	publish <- NewEvent("message", msg.User, int(msg.Timestamp), msg.Message)
}

// Leave is to make leave event
func Leave(user string) {
	timestamp := time.Now().Unix()
	publish <- NewEvent("leave", user, int(timestamp), "")
}

// Cancel is to define the action of unsubscription
func (s Subscription) Cancel() {
	unsubscribe <- s.New

	for { // infinite loop
		select {
		case _, ok := <-s.New:
			if !ok {
				return
			}
		default:
			return
		}
	}
}
