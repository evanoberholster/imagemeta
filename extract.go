// Package exiftool provides golang bindings for calling exiftool and
// working with the metadata it is able to extract from a media file
package exiftool

import (
	"bytes"
	"io"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
)

var ErrFilenameInvalid = errors.New("Filename contains control characters")

// Extract calls a specific exiftool with specific CLI flags
func Extract(exiftool, filename string, flags ...string) ([]byte, error) {

	if !strconv.CanBackquote(filename) {
		return nil, ErrFilenameInvalid
	}

	flags = append(flags, filename)
	cmd := exec.Command(exiftool, flags...)
	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// exiftool will exit and print valid output to stdout
	// if it exits with an unrecognized filetype, don't process
	// that situtation here
	if err != nil && stdout.Len() == 0 {
		return nil, errors.Errorf("%s", stderr.String())
	}

	// no exit error but also no output
	if stdout.Len() == 0 {
		return nil, errors.New("No output")
	}

	return stdout.Bytes(), nil
}

// ExtractReader extracts EXIF/metadata from an io.Reader, passing data to
// exiftool via stdin
func ExtractReader(exiftool string, source io.Reader, flags ...string) ([]byte, error) {
	flags = append(flags, "-")
	cmd := exec.Command(exiftool, flags...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = source

	err := cmd.Run()

	// exiftool will exit and print valid output to stdout
	// if it exits with an unrecognized filetype, don't process
	// that situtation here
	if err != nil && stdout.Len() == 0 {
		return nil, errors.Errorf("%s", stderr.String())
	}

	// no exit error but also no output
	if stdout.Len() == 0 {
		return nil, errors.New("No output")
	}

	return stdout.Bytes(), nil
}
