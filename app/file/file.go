package file

import (
	"VulpesEditor/app/util"
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"
)

type ArchiveWriter struct {
	file *os.File
	zip  *zip.Writer
}

type Project struct {
	Name string
	Path string
}

func GetAllProjects(kind string) (projects []Project) {
	projectsDir := filepath.Join(util.AppDir, "projects", kind)

	if files, err := os.ReadDir(projectsDir); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				var p Project
				p.Name = file.Name()
				p.Path = filepath.Join(projectsDir, file.Name())
				projects = append(projects, p)
			}
		}
	}

	return
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

func Load(path string) (r *ArchiveReader, err error) {
	f, err := zip.OpenReader(path)
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
