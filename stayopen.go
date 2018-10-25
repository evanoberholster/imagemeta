package exiftool

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	//"time"
	//"reflect"

	"github.com/pkg/errors"
)

// Stayopen abstracts running exiftool with `-stay_open` to greatly improve
// performance. Remember to call Stayopen.Stop() to signal exiftool to shutdown
// to avoid zombie perl processes
type Stayopen struct {
	l   sync.Mutex
	cmd *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser

	// default flags to pass to every extract call
	defaultFlags string

	scanner *bufio.Scanner
}

// Extract calls exiftool on the supplied filename
func (e *Stayopen) Extract(filename string) ([]byte, error) {
	return e.ExtractFlags(filename)
}

func (e *Stayopen) ExtractFlags(filename string, flags ...string) ([]byte, error) {
	e.l.Lock()
	defer e.l.Unlock()

	if e.cmd == nil {
		return nil, errors.New("Stopped")
	}

	if !strconv.CanBackquote(filename) {
		return nil, ErrFilenameInvalid
	}

	// send the request
	fmt.Fprintln(e.stdin, e.defaultFlags)
	if len(flags) > 0 {
		fmt.Fprintln(e.stdin, strings.Join(flags, "\n"))
	}

	fmt.Fprintln(e.stdin, filename)
	fmt.Fprintln(e.stdin, "-execute")

	if !e.scanner.Scan() {
		return nil, errors.New("Failed to read output")
	} else {
		
		results := e.scanner.Bytes()
		sendResults := make([]byte, len(results), len(results))
		copy(sendResults, results)
		return sendResults, nil
	}

}

func (e *Stayopen) ExtractReader(source io.Reader, flags ...string) ([]byte, error) {
	e.l.Lock()
	defer e.l.Unlock()

	if e.cmd == nil {
		return nil, errors.New("Stopped")
	}
	// send the request
	fmt.Fprintln(e.stdin, e.defaultFlags)
	if len(flags) > 0 {
		//fmt.Fprintln(e.stdin, strings.Join(flags, "\n"))
	}
	fmt.Fprintln(e.stdin, "-")
	//fmt.Println(io.Copy(e.stdin, source))
	fmt.Fprintln(e.stdin, source)
	//fmt.Fprintln(e.stdin, "\n")
	//fmt.Fprintln(e.stdin, "../testImages/image1mini.jpg")
	fmt.Fprintln(e.stdin, "-execute")

	if !e.scanner.Scan() {
		return nil, errors.New("Failed to read output")
	} else {
		results := e.scanner.Bytes()
		sendResults := make([]byte, len(results), len(results))
		copy(sendResults, results)
		return sendResults, nil
	}
}

func (e *Stayopen) Stop() {
	e.l.Lock()
	defer e.l.Unlock()

	// write message telling it to close
	// but don't actually wait for the command to stop
	fmt.Fprintln(e.stdin, "-stay_open")
	fmt.Fprintln(e.stdin, "False")
	fmt.Fprintln(e.stdin, "-execute")
	e.cmd = nil
}

func NewStayOpen(exiftool string, flags ...string) (*Stayopen, error) {

	var defaultFlags string
	if len(flags) > 0 {
		defaultFlags = strings.Join(flags, "\n")
	}
	stayopen := &Stayopen{
		defaultFlags: defaultFlags,
	}

	stayopen.cmd = exec.Command(exiftool, "-stay_open", "True", "-@", "-")

	stdin, err := stayopen.cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting stdin pipe")
	}

	stdout, err := stayopen.cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting stdout pipe")
	}

	stayopen.stdin = stdin
	stayopen.stdout = stdout
	stayopen.scanner = bufio.NewScanner(stdout)
	stayopen.scanner.Split(splitReadyToken)

	if err := stayopen.cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "Failed starting exiftool in stay_open mode")
	}

	// wait for both go-routines to startup
	return stayopen, nil
}

func splitReadyToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.Index(data, []byte("\n{ready}\n")); i >= 0 {
		if atEOF && len(data) == (i+9) { // nothing left to scan
			return i + 9, data[:i], bufio.ErrFinalToken
		} else {
			return i + 9, data[:i], nil
		}
	}

	if atEOF {
		return 0, data, io.EOF
	}

	return 0, nil, nil
}
