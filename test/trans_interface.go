package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
)

func trans() {
	type GobCodec struct {
		conn io.ReadWriteCloser
		buf  *bufio.Writer
		dec  *gob.Decoder
		enc  *gob.Encoder
	}

	type Codec interface {
		//PrintMain()
	}

	var _ Codec = (*GobCodec)(nil)
	type a string
	var _ a = "aaa"

	fmt.Println("")
}
func main() {
	trans()
}
