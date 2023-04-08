package preview

import (
	"io"

	"github.com/evanoberholster/imagemeta/meta"
	"github.com/rs/zerolog"
)

type previewReader struct {
	logger zerolog.Logger

	PreviewImage []byte
}

func NewPreviewReader(l zerolog.Logger) previewReader {
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
				Msgf("error read preview image")
			return err
		}
		if readLength == 0 {
			break
		}

		offset += uint32(readLength)
	}

	pr.PreviewImage = img

	return nil
}
