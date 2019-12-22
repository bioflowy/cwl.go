package cwl

import (
	"io"
	"os"
	"path/filepath"
)

// Entry represents fs entry, it means [File|Directory|Dirent]
type Entry struct {
	Class    string
	Location string
	Path     string
	Basename string
	File
	Directory
	Dirent
}

// File represents file entry.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#File
type File struct {
	Dirname string
	Size    int64
	Format  string
}

// Directory represents direcotry entry.
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#Directory
type Directory struct {
	Listing []Entry
}

// Dirent represents ?
// @see http://www.commonwl.org/v1.0/CommandLineTool.html#Dirent
type Dirent struct {
	Entry     string
	EntryName string
	Writable  bool
}

// NewList constructs a list of Entry from interface
func (_ Entry) NewList(i interface{}) []Entry {
	dest := []Entry{}
	switch x := i.(type) {
	case string:
		dest = append(dest, Entry{}.New(x))
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Entry{}.New(v))
		}
	}
	return dest
}

// New constructs an Entry from interface
func (_ Entry) New(i interface{}) Entry {
	dest := Entry{}
	switch x := i.(type) {
	case string:
		dest.Location = x
	case map[string]interface{}:
		for key, v := range x {
			switch key {
			case "entryname":
				dest.EntryName = v.(string)
			case "entry":
				dest.Entry = v.(string)
			case "writable":
				dest.Writable = v.(bool)
			}
		}
	}
	return dest
}
func copyFileContents(src, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, in); err != nil {
        return
    }
    err = out.Sync()
    return
}
func linkOrCopy(src string,dist string) error {
	err := os.Link(src, dist)
	if err != nil {
		return copyFileContents(src,dist)
	}
	return nil
}
// LinkTo creates hardlink of this entry under destdir.
func (entry *Entry) LinkTo(destdir, srcdir string) error {
	destpath := filepath.Join(destdir, filepath.Base(entry.Location))
	if filepath.IsAbs(entry.Location) {
		return linkOrCopy(entry.Location, destpath)
	}
	return linkOrCopy(filepath.Join(srcdir, entry.Location), destpath)
}
