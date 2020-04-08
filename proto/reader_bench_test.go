package proto_test

import (
	"bytes"
	"github.com/clovers4/gres/proto"
	"testing"
)

func BenchmarkReader_ParseReply_Status(b *testing.B) {
	benchmarkParseReply(b, "+OK\r\n", false)
}

func BenchmarkReader_ParseReply_Int(b *testing.B) {
	benchmarkParseReply(b, ":1\r\n", false)
}

func BenchmarkReader_ParseReply_Error(b *testing.B) {
	benchmarkParseReply(b, "-Error message\r\n", true)
}

func BenchmarkReader_ParseReply_String(b *testing.B) {
	benchmarkParseReply(b, "$5\r\nhello\r\n", false)
}

func BenchmarkReader_ParseReply_Slice(b *testing.B) {
	benchmarkParseReply(b, "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", false)
}

func benchmarkParseReply(b *testing.B, reply string, wanterr bool) {
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		buf.WriteString(reply)
	}
	p := proto.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := p.ReadReply()
		if !wanterr && err != nil {
			b.Fatal(err)
		}
	}
}

func multiBulkParse(p *proto.Reader, n int64) (interface{}, error) {
	vv := make([]interface{}, 0, n)
	for i := int64(0); i < n; i++ {
		v, err := p.ReadReply()
		if err != nil {
			return nil, err
		}
		vv = append(vv, v)
	}
	return vv, nil
}
