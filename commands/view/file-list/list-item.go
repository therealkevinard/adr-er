package file_list

import (
	"path"
	"time"

	"github.com/dustin/go-humanize"
)

// Item is a single item to render in the FileListModel.
type Item struct {
	name     string
	parent   string
	modified time.Time
}

// NewItem builds a new item from input.
func NewItem(name, parent string, modtime time.Time) Item {
	return Item{
		name:     name,
		parent:   parent,
		modified: modtime,
	}
}

// Title is used by list.DefaultDelegate.
func (i Item) Title() string { return i.name }

// Description is used by list.DefaultDelegate.
func (i Item) Description() string {
	return humanize.RelTime(i.modified, time.Now(), "ago", "from now")
}

// FilterValue returns the value to reference when the list is in filter mode.
func (i Item) FilterValue() string {
	return i.name
}

// FullPath returns the absolute path to the file item.
func (i Item) FullPath() string {
	// TODO: this is currently muy naive.
	return path.Join(i.parent, i.name)
}
