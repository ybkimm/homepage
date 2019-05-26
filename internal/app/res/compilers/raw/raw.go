package raw

import "go.ybk.im/homepage/internal/app/res/compilers"

type Compiler struct {
	contentType string
}

func NewCompiler(contentType string) compilers.Compiler {
	return &Compiler{
		contentType: contentType,
	}
}

func (c *Compiler) ContentType() string {
	return c.contentType
}

func (*Compiler) Compile(src []byte) ([]byte, error) {
	return src, nil
}
