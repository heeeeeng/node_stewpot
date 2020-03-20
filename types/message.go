package types

import "encoding/json"

const (
	DefualtMsgSize = 256 * SizeByte
)

type Message struct {
	Source     Node
	ID         int64
	Difficulty int64
	Size       int64
	Content    string
}

func NewMessage(source Node, difficulty int64, id int64, size int64, content string) Message {
	return Message{
		Source:     source,
		Difficulty: difficulty,
		ID:         id,
		Size:       size,
		Content:    content,
	}
}

type MessageJson struct {
	Source     string `json:"source"`
	ID         int64  `json:"id"`
	Difficulty int64  `json:"difficulty"`
	Size       int64  `json:"size"`
	Content    string `json:"content"`
}

func (t *Message) MarshalJSON() ([]byte, error) {
	j := MessageJson{}
	j.Source = t.Source.IP()
	j.ID = t.ID
	j.Difficulty = t.Difficulty
	j.Size = t.Size
	j.Content = t.Content

	return json.Marshal(j)
}
