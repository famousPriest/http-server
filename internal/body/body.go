package body

// TODO: need to add tests
type Body struct {
	Data []byte
}

func (b *Body) Parse(data []byte, contentLength int) {
	b.Data = append(b.Data, data...)
}
