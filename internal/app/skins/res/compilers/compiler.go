package compilers

type Compiler interface {
	ContentType() string
	Compile(src []byte) ([]byte, error)
}
