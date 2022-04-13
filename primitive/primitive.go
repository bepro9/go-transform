package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Mode defines the shapes used when transforming images
type Mode int

// Mode supported by the primitive package
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

// Transform will take the provided image and apply the primitive tarnsformation on it and then
// return a reader to the resulting image
func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}
	
	in, err := tempFile("in_", ext)
	if err != nil {
		return nil, errors.New("primitive : unable to create the temp input file")
	}
	defer os.Remove(in.Name())

	out, err := tempFile("in_", ext)
	if err != nil {
		return nil, errors.New("primitive : unable to create the temp output file")
	}
	defer os.Remove(out.Name())

	// Read image into a file (in)
	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("primitive : failed to copy image into temp file ")
	}

	// Run the primitive w/ -i in.Name() -o out.Name()
	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, args...)
	if err != nil {
		return nil, errors.New("primitive : failed to Transform the imput image ")
	}
	fmt.Println(stdCombo)

	// Read out a reader, return reader, delete reader
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func primitive(inputFile, outputFile string, numShapes int,args ...string) (string, error) {
	argStr := fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, numShapes)
	args = append(strings.Fields(argStr), args...)
	cmd := exec.Command("primitive", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func tempFile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, errors.New("primitive : unable to create the temp  file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}
