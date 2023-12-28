package migrate

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4/source"
)

var _ source.Driver = (*EmbedDriver)(nil)

func NewFsDriver() *EmbedDriver {
	e := &EmbedDriver{
		sorts: make([]uint, 0, 20),
		files: make(map[string]file),
	}
	return e
}

type EmbedDriver struct {
	sorts []uint
	files map[string]file
}

func (e *EmbedDriver) Add(f fs.FS) {
	dirEntries, err := fs.ReadDir(f, ".")
	if err != nil {
		panic(err)
	}
	for _, entry := range dirEntries {
		// 取version
		tmp, err := strconv.ParseUint(strings.Split(entry.Name(), "_")[0], 10, 64)
		if err != nil {
			panic(fmt.Errorf("file name invalid: %w", err))
		}
		version := uint(tmp)
		has := false
		for _, v := range e.sorts {
			if v == version {
				has = true
				break
			}
		}
		if !has {
			e.sorts = append(e.sorts, version)
		}

		//取file
		f := file{
			fs:   f,
			name: entry.Name(),
		}
		key := fmt.Sprintf("%d_%s", version, strings.Split(entry.Name(), ".")[1])
		e.files[key] = f
	}
	sort.Slice(e.sorts, func(i, j int) bool {
		return e.sorts[i] < e.sorts[j]
	})
}

func (e *EmbedDriver) Open(url string) (source.Driver, error) {
	return e, nil
}

func (e *EmbedDriver) Close() error {
	return nil
}

func (e *EmbedDriver) First() (version uint, err error) {
	if len(e.sorts) < 1 {
		return 0, os.ErrNotExist
	}
	return e.sorts[0], nil
}

func (e *EmbedDriver) find(version uint) (index int, err error) {
	var ver uint
	for index, ver = range e.sorts {
		if ver == version {
			return
		}
	}
	return 0, os.ErrNotExist
}

func (e *EmbedDriver) Prev(version uint) (prevVersion uint, err error) {
	index, err := e.find(version)
	if err != nil {
		return
	}
	if index-1 >= 0 {
		return e.sorts[index-1], nil
	}
	return 0, os.ErrNotExist
}

func (e *EmbedDriver) Next(version uint) (nextVersion uint, err error) {
	index, err := e.find(version)
	if err != nil {
		return
	}
	if index+1 < len(e.sorts) {
		return e.sorts[index+1], nil
	}
	return 0, os.ErrNotExist
}

func (e *EmbedDriver) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	file, ok := e.files[fmt.Sprintf("%d_%s", version, "up")]
	if !ok {
		return nil, "", os.ErrNotExist
	}
	identifier = file.name
	r, err = file.Open()
	return
}

func (e *EmbedDriver) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	file, ok := e.files[fmt.Sprintf("%d_%s", version, "down")]
	if !ok {
		return nil, "", os.ErrNotExist
	}
	identifier = file.name
	r, err = file.Open()
	return
}

type file struct {
	fs   fs.FS
	name string
}

func (f file) Open() (io.ReadCloser, error) {
	return f.fs.Open(f.name)
}
