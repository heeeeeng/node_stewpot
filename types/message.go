package types

const (
	DefualtMsgSize = 256 * Byte
)

type Message struct {
	Source     Node
	Difficulty int
	ID         int64
	Size       FileSize
}

func NewMessage(source Node, d int, id int64, size FileSize) Message {
	return Message{
		Source:     source,
		Difficulty: d,
		ID:         id,
		Size:       size,
	}
}
