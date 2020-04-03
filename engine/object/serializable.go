package object

import "io"

type Serializable interface {
	Marshal(w io.Writer) error
	Unmarshal(r io.Reader) error
}
