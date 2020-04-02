package container

type Serializable interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte)
}
