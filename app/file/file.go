package file

import (
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"
)

type ArchiveWriter struct {
	file *os.File
	zip  *zip.Writer
}

func NewArchive(path, name string) (w *ArchiveWriter, err error) {
	os.MkdirAll(path, os.ModePerm)
	w = new(ArchiveWriter)
	w.file, err = os.Create(filepath.Join(path, name+".zip"))
	if err != nil {
		return
	}
	w.zip = zip.NewWriter(w.file)
	return
}

func (z *ArchiveWriter) Write(name string, data []byte) (err error) {
	f, err := z.zip.Create(name)
	if err != nil {
		return err
	}
	f.Write(data)
	return
}

func (z *ArchiveWriter) Save() (err error) {
	if err := z.zip.Close(); err != nil {
		return err
	}
	if err := z.file.Close(); err != nil {
		return err
	}
	return
}

type ArchiveReader struct {
	closer *zip.ReadCloser
}

func Load(name string) (r *ArchiveReader, err error) {
	f, err := zip.OpenReader(name)
	if err != nil {
		return
	}
	r = new(ArchiveReader)
	r.closer = f
	return
}

func (r *ArchiveReader) Open(name string) (file fs.File, err error) {
	file, err = r.closer.Open(name)
	return
}

func (r *ArchiveReader) Close() {
	r.closer.Close()
}
