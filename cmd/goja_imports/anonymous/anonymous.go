package anonymous

import (
	"fmt"
	"image"
	im "image"

	"github.com/garfeng/gojatools/cmd/goja_imports/anonymous/impkg"
)

type sub1 struct {
}

type Sub2 struct {
}

type iSub interface {
	InterfaceFunc()
}

func (s *sub1) Sub1Func() {
	fmt.Println("sub1 func")
}

func (s *Sub2) Sub2Func() {
	fmt.Println("sub2 func")
}

type Top struct {
	*Sub2
	*sub1
	iSub
}

func UseImage(m impkg.Image) im.Image {
	return im.NewRGBA(im.Rect(0, 0, 10, 10))
}

func UseMyImage(m impkg.Image2) impkg.Image {
	return im.NewRGBA(image.Rect(0, 0, 10, 10))
}
