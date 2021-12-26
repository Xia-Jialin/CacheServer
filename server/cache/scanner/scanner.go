package scanner

type Scanner interface {
	Scan() bool
	Key() string
	Value() []byte
	Close()
}
