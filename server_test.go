package gres

import (
	"fmt"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	s := NewServer(ConnectionTimeoutOption(10*time.Second), DbnumOption(8))
	fmt.Printf("%+v\n", s)
	s.Start()
}
