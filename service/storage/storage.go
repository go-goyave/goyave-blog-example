package storage

import (
	"fmt"
	"io/fs"

	"github.com/go-goyave/goyave-blog-example/service"
	"goyave.dev/goyave/v5/util/fsutil"
)

type FS interface {
	fsutil.FS
	fsutil.WritableFS
	fsutil.RemoveFS
}

type Service struct {
	FS FS
}

func NewService(fs FS) *Service {
	return &Service{
		FS: fs,
	}
}

func (s *Service) GetFS() fs.StatFS {
	return s.FS
}

func (s *Service) SaveAvatar(file fsutil.File) (string, error) {
	return file.Save(s.FS, "", fmt.Sprintf("user_avatar_%s", file.Header.Filename))
}

func (s *Service) Delete(name string) error {
	return s.FS.Remove(name)
}

func (*Service) Name() string {
	return service.Storage
}
