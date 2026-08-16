package preview

import (
	"io"
	"log/slog"

	"github.com/evanoberholster/imagemeta/meta"
)

type previewReader struct {
	logger *slog.Logger

	PreviewImage []byte
}

func NewPreviewReader(l *slog.Logger) previewReader {
	ir := previewReader{
		logger: l,
	}
	return ir
}

func (pr *previewReader) RenderPreview(r io.Reader, h meta.PreviewHeader) error {
	img := make([]byte, h.Size)
	offset := uint32(0)
	maxSize := uint32(2048)
	for {
		maxOffset := offset + maxSize
		if h.Size < maxOffset {
			maxOffset = h.Size
		}

		readLength, err := r.Read(img[offset:maxOffset])
		if err != nil {
			if err == io.EOF {
				break
			}
			pr.logError(err).
				Uint32("offset", offset).
				Uint32("maxOffset", maxOffset).
				Msg("error read preview image")
			return err
		}
		if readLength == 0 {
			break
		}

		delta, ok := meta.SafecastIntToUint32(readLength)
		if !ok {
			return io.ErrUnexpectedEOF
		}
		offset += delta
	}

	pr.PreviewImage = img

	return nil
}
