package ex11filemanager

import (
	"errors"
	"os"
	"runtime"
)

func ReadEntry(entry *Entry) ([]*Entry, error) {
	if entry == nil {
		return nil, errors.New("entry is nil")
	}

	if !entry.IsDir {
		return nil, errors.New("entry is not a directory")
	}

	switch entry.DriverType {
	case DriverTypeLocal:
		return readLocalDirectory(entry)
	case DriverTypeNetwork:
		return readNetworkDirectory(entry)
	default:
		return nil, errors.New("unknown driver type")
	}
}

func readRootEntries() []*Entry {
	entries := make([]*Entry, 0)
	if runtime.GOOS == "windows" {
		// get drives on Windows
		for _, drive := range []string{"C:", "D:", "E:", "F:"} {
			if _, err := os.Stat(drive); err == nil {
				entry := NewEntry()
				entry.ServicePath = []string{drive}
				entry.IsDir = true
				entry.DriverType = DriverTypeLocal
				entries = append(entries, entry)
			}
		}

	}

	return entries
}

func readLocalDirectory(entry *Entry) ([]*Entry, error) {
	fullPath := entry.FullPath()
	dirEntries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	result := make([]*Entry, 0)

	for _, e := range dirEntries {
		entry := NewEntry()
		entry.ServicePath = make([]string, len(entry.ServicePath)+1)
		copy(entry.ServicePath, entry.ServicePath)
		entry.ServicePath[len(entry.ServicePath)-1] = e.Name()
		entry.IsDir = e.IsDir()
		entry.DriverType = DriverTypeLocal

		fileInfo, err := e.Info()
		if err != nil {
			entry.Error = err
			continue
		}

		entry.Size = fileInfo.Size()
		entry.Created = fileInfo.ModTime()
		entry.Modified = fileInfo.ModTime()
		result = append(result, entry)
	}

	return result, nil
}

func readNetworkDirectory(_ *Entry) ([]*Entry, error) {
	return nil, errors.New("not implemented")
}
