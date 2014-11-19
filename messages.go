package socketio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type Message struct {
	Type int
	ID   int
	Body []byte
	Ack  func(message []byte)
}

var currentIndex = 1

func incrementMessageIndex() {
	if currentIndex == 128 {
		currentIndex = 1
	} else {
		currentIndex++
	}
}

type Event struct {
	Name string            `json:"name"`
	Args []json.RawMessage `json:"args"`
}

func CreateMessageEvent(message string, ack func(message []byte)) Message {

	var temp json.RawMessage
	json.Unmarshal([]byte(message), &temp)
	tempArray := []json.RawMessage{temp}

	messageBody := Event{
		Name: "message",
		Args: tempArray,
	}

	messageEvent := Message{
		Type: 5,
		ID:   currentIndex,
		Ack:  ack,
	}
	incrementMessageIndex()

	tempMessage, err := json.Marshal(messageBody)
	if err != nil {
		fmt.Println("error on marshal: ", err)
		return Message{}
	}

	messageEvent.Body = tempMessage

	return messageEvent

}

func CreateMessageHeartbeat() Message {
	message := Message{
		Type: 2,
	}
	return message
}

func (message Message) PrintMessage() string {
	switch message.Type {
	case 2:
		return "2::"
	case 5:
		return "5:" + strconv.Itoa(message.ID) + "+::" + string(message.Body)
	default:
		return ""
	}
}

func parseEvent(buffer []byte) (string, []byte) {
	var event Event
	index := bytes.Index(buffer, []byte("{"))
	json.Unmarshal(buffer[index:], &event)
	return event.Name, event.Args[0]
}

func parseMessage(buffer []byte) []byte {
	splitChunks := bytes.Split(buffer, []byte(":"))
	if len(splitChunks) < 4 {
		return []byte("")
	}
	return splitChunks[3]
}

func parseAck(buffer []byte) []byte {
	return []byte("")
}
