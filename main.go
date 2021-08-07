package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.prodigy9.co/stega/images"
)

var rootCmd = &cobra.Command{
	Use:   "stega",
	Short: "Quick-and-dirty steganography tool",
}

var encodeCmd = &cobra.Command{
	Use:   "encode filename text",
	Short: "Encodes hidden text into the image file",
	RunE:  runEncode,
}

var decodeCmd = &cobra.Command{
	Use:   "decode filename",
	Short: "Decodes hidden text from the image file",
	RunE:  runDecode,
}

func init() {
	rootCmd.AddCommand(encodeCmd, decodeCmd)
}

func runEncode(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		log.Fatalln("requires 2 arguments, filename and the text to encode")
		return nil
	}

	if img, err := load(args[0]); err != nil {
		return err
	} else if err = img.SetHiddenText(args[1]); err != nil {
		return err
	} else if err = save(args[0], img); err != nil {
		return err
	} else {
		return nil
	}
}

func runDecode(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		log.Fatalln("requires 1 argument, filename")
		return nil
	}

	if img, err := load(args[0]); err != nil {
		return err
	} else if text, err := img.GetHiddenText(); err != nil {
		return err
	} else {
		log.Println("recovered text: " + text)
		return nil
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func load(inname string) (*images.Image, error) {
	file, err := os.Open(inname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return images.New(file)
}

func save(outname string, img *images.Image) error {
	buf := &bytes.Buffer{}
	if _, err := img.WriteTo(buf); err != nil {
		return err
	}

	file, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}

	tmpname := file.Name()
	file.Close()
	if err = ioutil.WriteFile(tmpname, buf.Bytes(), 0644); err != nil {
		return err
	} else if err = os.Rename(tmpname, outname); err != nil {
		return err
	} else {
		os.Remove(tmpname)
		return nil
	}
}
