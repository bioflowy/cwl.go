package cwl

import (
	"fmt"
	"io/ioutil"
)

// Provided represents the provided input value
// by parameter files.
type Provided struct {
	ID    string
	Raw   interface{}
	Entry *Entry // In most cases, it's "File" if "Entry" exists
	Error error

	// TODO: Refactor
	Int int
}

// New constructs new "Provided" struct.
func (provided Provided) New(id string, i interface{}) *Provided {
	dest := &Provided{ID: id, Raw: i}
	switch v := i.(type) {
	case nil:
		return nil
	case int:
		dest.Int = v
	case map[interface{}]interface{}: // It's "File" in most cases
		dest.Entry, dest.Error = dest.EntryFromDictionary(v)
	}
	return dest
}

// EntryFromDictionary ...
func (provided Provided) EntryFromDictionary(dict map[interface{}]interface{}) (*Entry, error) {
	if dict == nil {
		return nil, nil
	}
	class := dict["class"].(string)
	location := dict["location"]
	path := dict["path"]
	contents := dict["contents"]
	if class == "" && location == nil && contents == nil {
		return nil, nil
	}
	switch class {
	case "File":
		// Use location if specified
		if location != nil {
			return &Entry{
				Class:    class,
				Location: fmt.Sprintf("%v", location),
				File:     File{},
			}, nil
		}
		// If the path field is provided but the location field is not, 
		// an implementation may assign the value of the path field to location
		if path != nil {
			return &Entry{
				Class:    class,
				Location: fmt.Sprintf("%v", path),
				File:     File{},
			}, nil
		}
		// Use contents if specified
		if contentsstring, ok := contents.(string); ok {
			tmpfile, err := ioutil.TempFile("/tmp", provided.ID)
			if err != nil {
				return nil, err
			}
			defer tmpfile.Close()
			if _, err := tmpfile.WriteString(contentsstring); err != nil {
				return nil, err
			}
			return &Entry{
				Class:    class,
				Location: tmpfile.Name(),
				File:     File{},
			}, nil
		}
	}
	return nil, nil
}
