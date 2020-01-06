package types

const (
	DefualtMsgSize = 256 * Byte
)

type Message struct {
	Source     Node
	Difficulty int
	ID         int64
	Size       int64
}

func NewMessage(source Node, d int, id int64, size int64) Message {
	return Message{
		Source:     source,
		Difficulty: d,
		ID:         id,
		Size:       size,
	}
}
