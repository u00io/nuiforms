package ex11filemanager

import (
	"runtime"
	"time"
)

const (
	DriverTypeLocal   = "local"
	DriverTypeNetwork = "network"
	DriverTypeUnknown = "unknown"
)

type Entry struct {
	ServicePath []string  // Path separated, e.g., ["C:", "path", "to", "file.txt"]
	Size        int64     // Size in bytes, 0 for directories
	Created     time.Time // Creation time
	Modified    time.Time // Last modified time
	IsDir       bool      // True if this entry is a directory
	DriverType  string    // Type of the driver (e.g., "local", "network")
	Error       error     // Error if any occurred while reading this entry
}

func NewEntry() *Entry {
	var c Entry
	return &c
}

func (c *Entry) FullPath() string {
	result := ""
	if c.DriverType == DriverTypeLocal {
		if runtime.GOOS == "windows" {
			for i, part := range c.ServicePath {
				if i < len(c.ServicePath)-1 {
					result += part + "\\"
				} else {
					result += part
				}
			}
			if len(c.ServicePath) == 1 {
				result += "\\"
			}
		}
	}
	return result
}

func (c *Entry) DisplayName() string {
	if len(c.ServicePath) == 0 {
		return ""
	}
	return c.ServicePath[len(c.ServicePath)-1]
}

func (c *Entry) CreateChildEntry(name string) *Entry {
	child := NewEntry()
	child.ServicePath = append(c.ServicePath, name)
	child.IsDir = true
	child.DriverType = c.DriverType
	return child
}

func (c *Entry) CreateParentEntry() *Entry {
	parent := NewEntry()
	if len(c.ServicePath) > 1 {
		parent.ServicePath = make([]string, len(c.ServicePath)-1)
		copy(parent.ServicePath, c.ServicePath[:len(c.ServicePath)-1])
	}
	parent.IsDir = true
	parent.DriverType = c.DriverType
	return parent
}
