package scss

import (
	"bytes"

	"go.ybk.im/homepage/internal/app/res/compilers"

	"github.com/wellington/go-libsass"
)

type Compiler struct {
	basePath string
}

func NewCompiler(basePath string) compilers.Compiler {
	return &Compiler{
		basePath: basePath,
	}
}

func (*Compiler) ContentType() string {
	return "text/css"
}

func (c *Compiler) Compile(src []byte) ([]byte, error) {
	srcBuf := bytes.NewBuffer(src)
	dstBuf := bytes.NewBuffer([]byte{})

	comp, err := libsass.New(
		dstBuf, srcBuf,
		libsass.Comments(false),
		libsass.IncludePaths([]string{c.basePath}),
		libsass.OutputStyle(libsass.COMPRESSED_STYLE),
		libsass.SourceMap(false, "", ""),
		libsass.WithSyntax(libsass.SCSSSyntax),
	)
	if err != nil {
		return nil, err
	}

	err = comp.Run()
	if err != nil {
		return nil, err
	}

	return dstBuf.Bytes(), nil
}
