package encoder

type Encoder interface {
	Encode(n int64) string
	Decode(token string) int64
}
