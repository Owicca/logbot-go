package buffer

import (
	"os"
	"path/filepath"
)

type Buffer struct {
	channel     string
	isConnected bool
	sink        *os.File
}

func NewBuffer(channel string) *Buffer {
	buf := &Buffer{}
	buf.isConnected = false
	buf.SetChannel(channel)

	return buf
}

func (b *Buffer) IsConnected() bool {
	return b.isConnected
}

func (b *Buffer) Channel() string {
	return b.channel
}

func (b *Buffer) SetChannel(channel string) (bool, error) {
	b.channel = channel

	sink_path := filepath.Join("./logs", channel)
	file, err := os.Open(sink_path)
	file.WriteString("test")
	b.sink = file

	if err != nil {
		return false, err
	}
	return true, err
}

func (b *Buffer) Write(line string) (bool, error) {
	written, err := b.sink.WriteString(line)

	if err != nil {
		return false, err
	}

	return written > 0, err
}

func (b *Buffer) Close() (bool, error) {
	if b.IsConnected() {
		err := b.sink.Close()
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
