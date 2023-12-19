package gen

import "sync"

type File struct {
	Name string
	base int64
	size int64
	set  *Files
}

type Files struct {
	files   []*File
	content []byte
	lock    sync.Mutex
}

func (f *Files) Find(name string) (file *File, err error) {
	for _, file = range f.files {
		if file.Name == name {
			return
		}
	}
	file = nil
	return
}

func (f *Files) Create(name string, size ...int) (err error) {
	s := 0
	if len(size) > 0 {
		s = size[0] - 1
	}
	if s <= 0 {
		f.files = append(f.files, &File{
			Name: name,
			base: -1,
		})
		return
	}

}

func (f *Files) write(base, content []byte) (err error) {

}
