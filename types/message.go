package types

const (
	DefualtMsgSize = 256 * SizeByte
)

type Message struct {
	Source     Node
	Difficulty int64
	ID         int64
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
