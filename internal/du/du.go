package du

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

var ComputeHashes = false

// File is the object that contains the info and path of the file
type File struct {
	Path         string
	HighDir      string
	Name         string
	Size         int64
	ApparentSize int64
	HumanSize    string
	Mode         os.FileMode
	ModTime      time.Time
	Hash         uint64 `hash:"ignore"`
}

type Folder struct {
	Path         string
	HighDir      string
	Name         string
	Size         int64
	ApparentSize int64
	HumanSize    string
	Mode         os.FileMode
	ModTime      time.Time
	Hash         uint64 `hash:"ignore"`
	Files        []File
	Folders      []Folder
}

type NameSorter []File

func (a NameSorter) Len() int           { return len(a) }
func (a NameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

type TimeSorter []File

func (a TimeSorter) Len() int           { return len(a) }
func (a TimeSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TimeSorter) Less(i, j int) bool { return a[i].ModTime.Before(a[j].ModTime) }

func PrettyPrintSize(size int64) string {
	switch {
	case size > 1024*1024*1024:
		return fmt.Sprintf("%.1fG", float64(size)/(1024*1024*1024))
	case size > 1024*1024:
		return fmt.Sprintf("%.1fM", float64(size)/(1024*1024))
	case size > 1024:
		return fmt.Sprintf("%.1fK", float64(size)/1024)
	default:
		return fmt.Sprintf("%d", size)
	}
}

// TODO(david): properly handle errors

func CreateFileTree(dir string) (root Folder, err error) {
	dir = filepath.Clean(dir)
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	info, err := f.Stat()
	if err != nil {
		return
	}

	// Prestep things before creating struct
	size := info.Size()
	fls, err := ioutil.ReadDir(dir)
	files := make([]File, 0)
	folders := make([]Folder, 0)
	if err != nil {
		return
	}
	for _, f := range fls {
		if f.IsDir() {
			folders = append(folders, createFileTreeHelper(path.Join(dir, f.Name())))
		} else {
			files = append(files, File{
				Path:      f.Name(),
				HighDir:   dir,
				Name:      f.Name(),
				Size:      f.Size(),
				HumanSize: PrettyPrintSize(f.Size()),
				Mode:      f.Mode(),
				ModTime:   f.ModTime(),
			})
		}
	}

	// Maybe not count directory as 4K?
	for _, file := range files {
		size += file.Size
	}
	for _, folder := range folders {
		size += folder.Size
	}

	root = Folder{
		Path:      dir,
		HighDir:   dir,
		Name:      info.Name(),
		Size:      size,
		HumanSize: PrettyPrintSize(size),
		Mode:      info.Mode(),
		ModTime:   info.ModTime(),
		Files:     files,
		Folders:   folders,
	}
	f.Close()
	return
}

func createFileTreeHelper(dir string) (root Folder) {
	dir = filepath.Clean(dir)
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	info, err := f.Stat()
	if err != nil {
		return
	}

	// Prestep things before creating struct
	size := info.Size()
	fls, err := ioutil.ReadDir(dir)
	files := make([]File, 0)
	folders := make([]Folder, 0)
	if err != nil {
		return
	}
	for _, f := range fls {
		if f.IsDir() {
			folders = append(folders, createFileTreeHelper(path.Join(dir, f.Name())))
		} else {
			files = append(files, File{
				Path:      f.Name(),
				HighDir:   dir,
				Name:      f.Name(),
				Size:      f.Size(),
				HumanSize: PrettyPrintSize(f.Size()),
				Mode:      f.Mode(),
				ModTime:   f.ModTime(),
			})
		}
	}

	// Maybe not count directory as 4K?
	for _, file := range files {
		size += file.Size
	}
	for _, folder := range folders {
		size += folder.Size
	}

	root = Folder{
		Path:      dir,
		HighDir:   dir,
		Name:      info.Name(),
		Size:      size,
		HumanSize: PrettyPrintSize(size),
		Mode:      info.Mode(),
		ModTime:   info.ModTime(),
		Files:     files,
		Folders:   folders,
	}
	f.Close()
	return
}
