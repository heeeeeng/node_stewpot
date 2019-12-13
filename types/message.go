package types

type Message struct {
	Source     Node
	Difficulty int
	ID         int64
}

func NewMessage(source Node, d int, id int64) Message {
	return Message{
		Source:     source,
		Difficulty: d,
		ID:         id,
	}
}
