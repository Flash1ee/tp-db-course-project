package models

const (
	DOWN = iota - 1
	UP   = iota + 1
)

type Vote struct {
	Nickname string
	ThreadID int64
	Voice    int64
}

func NewVote(isUp bool) *Vote {
	if isUp {
		return &Vote{
			Voice: UP,
		}
	}
	return &Vote{
		Voice: DOWN,
	}
}

func (v Vote) isValid() bool {
	return v.Voice == UP || v.Voice == DOWN
}
