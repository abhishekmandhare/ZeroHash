package stream

type Streamer interface {
	StreamData(data chan StreamData)
}

func NewStreamer() Streamer {
	return &ConsoleStreamer{}
}
