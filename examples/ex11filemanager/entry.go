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

	selectedChildIndex      int  // Index of the selected child in the file list
	isLinkToParentDirectory bool // True if this entry is a link to the parent directory
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

	if c.isLinkToParentDirectory {
		return ".."
	}

	return c.ServicePath[len(c.ServicePath)-1]
}
