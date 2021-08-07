package images

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"io/ioutil"
	"strconv"
	"strings"

	"io"
	"unicode/utf8"

	"github.com/auyer/steganography"
	"go.prodigy9.co/stega/logger"
)

var log = logger.For("images")

type Image struct {
	buf []byte
}

func New(reader io.Reader) (*Image, error) {
	if buf, err := ioutil.ReadAll(reader); err != nil {
		return nil, log.WrapErr("", err)
	} else {
		return &Image{buf}, nil
	}
}

func (i *Image) GetHiddenText() (string, error) {
	if img, _, err := image.Decode(bytes.NewBuffer(i.buf)); err != nil {
		return "", log.WrapErr("", err)
	} else if lenBuf := steganography.Decode(4, img); !utf8.Valid(lenBuf) {
		log.Println(string(lenBuf))
		return "", log.WrapErr("", errors.New("no hidden text or data was corrupted"))
	} else if textLen, err := strconv.Atoi(strings.TrimSpace(string(lenBuf))); err != nil {
		return "", log.WrapErr("", errors.New("text length field is corrupted"))
	} else if textBuf := steganography.Decode(4 /*len*/ +1 /*sep*/ +uint32(textLen), img); !utf8.Valid(textBuf) {
		return "", log.WrapErr("", errors.New("text data is corrupted"))
	} else {
		return string(textBuf), nil
	}
}

func (i *Image) SetHiddenText(hiddenText string) error {
	textBuf := []byte(hiddenText)
	lineBuf := &bytes.Buffer{}
	lineBuf.Write([]byte(fmt.Sprintf("%4d", len(textBuf))))
	lineBuf.WriteRune('|')
	lineBuf.Write(textBuf)

	log.Println(string(lineBuf.Bytes()))

	buf := &bytes.Buffer{}
	if img, _, err := image.Decode(bytes.NewBuffer(i.buf)); err != nil {
		return log.WrapErr("", err)
	} else if err := steganography.Encode(buf, img, lineBuf.Bytes()); err != nil {
		return log.WrapErr("steganography", err)
	} else {
		i.buf = buf.Bytes()
		return nil
	}
}

func (i *Image) WriteTo(writer io.Writer) (int64, error) {
	return bytes.NewBuffer(i.buf).WriteTo(writer)
}
