package writer

type Writer interface {
	WriteData(data chan WriterData)
}

func NewWriter() Writer {
	return &ConsoleWriter{}
}
