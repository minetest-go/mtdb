package types

import "archive/zip"

// provides a backup interface for imports and exports
type Backup interface {
	Export(z *zip.Writer) error
	Import(z *zip.Reader) error
}
