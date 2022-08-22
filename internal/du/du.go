package du

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/mitchellh/hashstructure"
)

var ComputeHashes = true

// File is the object that contains the info and path of the file
type File struct {
	Path         string
	HighDir      string
	Name         string
	Size         int64
	ApparentSize int64
	HumanSize    string // bytesize.ByteSize
	Mode         os.FileMode
	ModTime      time.Time
	IsDir        bool
	Hash         uint64 `hash:"ignore"`
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

// ListFilesRecursivelyInParallel uses goroutines to list all the files
func ListFilesRecursivelyInParallel(dir string) (files []File, err error) {
	dir = filepath.Clean(dir)
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	info, err := f.Stat()
	if err != nil {
		return
	}
	files = []File{
		{
			Path:      dir,
			HighDir:   dir,
			Name:      dir,
			Size:      info.Size(),
			HumanSize: PrettyPrintSize(info.Size()),
			Mode:      info.Mode(),
			ModTime:   info.ModTime(),
			IsDir:     info.IsDir(),
		},
	}
	f.Close()

	if ComputeHashes {
		h, err := hashstructure.Hash(files[0], nil)
		if err != nil {
			panic(err)
		}
		files[0].Hash = h
	}

	fileChan := make(chan File)
	startedDirectories := make(chan bool)
	go listFilesInParallel(dir, startedDirectories, fileChan)

	runningCount := 1
	for {
		select {
		case file := <-fileChan:
			files = append(files, file)
		case newDir := <-startedDirectories:
			if newDir {
				runningCount++
			} else {
				runningCount--
			}
		default:
		}
		if runningCount == 0 {
			break
		}
	}
	return
}

func listFilesInParallel(dir string, startedDirectories chan bool, fileChan chan File) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fileStruct := File{
			Path:      path.Join(dir, f.Name()),
			HighDir:   dir,
			Name:      f.Name(),
			Size:      f.Size(),
			HumanSize: PrettyPrintSize(f.Size()),
			Mode:      f.Mode(),
			ModTime:   f.ModTime(),
			IsDir:     f.IsDir(),
		}
		if ComputeHashes {
			h, err := hashstructure.Hash(fileStruct, nil)
			if err != nil {
				panic(err)
			}
			fileStruct.Hash = h
		}
		fileChan <- fileStruct
		if f.IsDir() {
			startedDirectories <- true
			go listFilesInParallel(path.Join(dir, f.Name()), startedDirectories, fileChan)
		}
	}
	startedDirectories <- false
	return
}
