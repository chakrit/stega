package images

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"io/ioutil"
	"strconv"

	"io"
	"unicode/utf8"

	"github.com/auyer/steganography"
	"go.prodigy9.co/stega/logger"
)

var log = logger.For("images")

const MessageLength = 256

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
	} else if buf := steganography.Decode(MessageLength, img); !utf8.Valid(buf) {
		return "", log.WrapErr("", errors.New("no hidden text or data was corrupted"))
	} else {
		return string(buf), nil
	}
}

func (i *Image) SetHiddenText(hiddenText string) error {
	textBuf := []byte(hiddenText)
	if len(textBuf) > MessageLength {
		return log.WrapErr("", errors.New("hidden text size cannot exceed "+strconv.Itoa(MessageLength)+" bytes"))
	}

	buf := &bytes.Buffer{}
	if img, _, err := image.Decode(bytes.NewBuffer(i.buf)); err != nil {
		return log.WrapErr("", err)
	} else if err := steganography.Encode(buf, img, textBuf); err != nil {
		return log.WrapErr("steganography", err)
	} else {
		i.buf = buf.Bytes()
		return nil
	}
}

func (i *Image) WriteTo(writer io.Writer) (int64, error) {
	return bytes.NewBuffer(i.buf).WriteTo(writer)
}
