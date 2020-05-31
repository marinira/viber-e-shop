package viberapi

import (
	"log"
	"time"
)

// myMsgReceivedFunc will be called everytime when user send us a message
func MyMsgReceivedFunc(v *Viber, u User, m Message, token uint64, t time.Time) {
	switch m.(type) {

	case *TextMessage:
		v.SendTextMessage(u.ID, "Thank you for your message")
		txt := m.(*TextMessage).Text
		v.SendTextMessage(u.ID, "This is the text you have sent to me "+txt)

	case *URLMessage:
		url := m.(*URLMessage).Media
		v.SendTextMessage(u.ID, "You have sent me an interesting link "+url)

	case *PictureMessage:
		v.SendTextMessage(u.ID, "Nice pic!")

	}
}

func MyDeliveredFunc(v *Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "delivered to user ID", userID)
}

func MySeenFunc(v *Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "seen by user ID", userID)
}

func MyConversationStarted(v *Viber, u User, conversationType, context string, subscribed bool, token uint64, t time.Time) Message {
	log.Println("Message ID", token, "seen by user ID", u.ID)
	m := v.NewTextMessage("hi")
	return m
}

func MyFailed(v *Viber, userID string, token uint64, descr string, t time.Time) {
	log.Println("Message ID", token, "seen by user ID", userID, descr)
}

// All events that you can assign your function, declarations must match
// ConversationStarted func(v *Viber, u User, conversationType, context string, subscribed bool, token uint64, t time.Time) Message
// Message             func(v *Viber, u User, m Message, token uint64, t time.Time)
// Subscribed          func(v *Viber, u User, token uint64, t time.Time)
// Unsubscribed        func(v *Viber, userID string, token uint64, t time.Time)
// Delivered           func(v *Viber, userID string, token uint64, t time.Time)
// Seen                func(v *Viber, userID string, token uint64, t time.Time)
// Failed              func(v *Viber, userID string, token uint64, descr string, t time.Time)
