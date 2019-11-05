package main

type Stewpot struct {
	protocol Protocol
}

func NewStewpot() *Stewpot {
	s := &Stewpot{}

	return s
}

func (s *Stewpot) InitNetwork() {

}

// Stew start the simulating of the nodes network.
func (s *Stewpot) Stew() {

}

func (s *Stewpot) FlushSteps() {

}

type Protocol interface {
	GetPackages() [][]byte
}
