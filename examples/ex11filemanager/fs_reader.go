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

func readLocalDirectory(entryToRead *Entry) ([]*Entry, error) {
	fullPath := entryToRead.FullPath()
	dirEntries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	result := make([]*Entry, 0)

	if len(entryToRead.ServicePath) > 1 {
		parentEntry := NewEntry()
		parentEntry.ServicePath = make([]string, len(entryToRead.ServicePath)-1)
		copy(parentEntry.ServicePath, entryToRead.ServicePath[:len(entryToRead.ServicePath)-1])
		parentEntry.IsDir = true
		parentEntry.DriverType = DriverTypeLocal
		parentEntry.selectedChildIndex = entryToRead.selectedChildIndex
		parentEntry.isLinkToParentDirectory = true
		result = append(result, parentEntry)
	}

	items := make([]*Entry, 0)
	for _, e := range dirEntries {
		en := NewEntry()
		en.ServicePath = make([]string, len(entryToRead.ServicePath)+1)
		copy(en.ServicePath, entryToRead.ServicePath)
		en.ServicePath[len(en.ServicePath)-1] = e.Name()
		en.IsDir = e.IsDir()
		en.DriverType = DriverTypeLocal

		fileInfo, err := e.Info()
		if err != nil {
			en.Error = err
			continue
		}

		en.Size = fileInfo.Size()
		en.Created = fileInfo.ModTime()
		en.Modified = fileInfo.ModTime()
		items = append(items, en)
	}

	// Add all directories first
	for _, item := range items {
		if item.IsDir {
			result = append(result, item)
		}
	}

	// Add all files
	for _, item := range items {
		if !item.IsDir {
			result = append(result, item)
		}
	}

	return result, nil
}

func readNetworkDirectory(_ *Entry) ([]*Entry, error) {
	return nil, errors.New("not implemented")
}
