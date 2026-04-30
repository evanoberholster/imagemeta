package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta/logging"
)

func init() {
	logging.Logger = logging.New(os.Stdout, slog.LevelDebug)
}

func main() {
	path := "1.cr3"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open %q: %v\n", path, err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && !errors.Is(closeErr, io.EOF) {
			fmt.Fprintf(os.Stderr, "close %q: %v\n", path, closeErr)
		}
	}()

	e, err := imagemeta.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode %q: %v\n", path, err)
		os.Exit(1)
	}
	buf, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "json %q: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Println(string(buf))
}
