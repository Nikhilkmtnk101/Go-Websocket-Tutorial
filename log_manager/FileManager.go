package log_manager

import (
	"fmt"
	"sync"
	"time"
)

var (
	fileManager *FileManager
	once        = sync.Once{}
)

type FileManager struct {
	File      []string
	cursor    int
	isUpdated bool
	mu        *sync.Mutex
}

func NewFileManager() {

	once.Do(func() {
		fileManager = &FileManager{
			File:      []string{},
			cursor:    -1,
			mu:        &sync.Mutex{},
			isUpdated: false,
		}
	})
}

func GetFileManager() *FileManager {
	if fileManager == nil {
		NewFileManager()
	}
	return fileManager
}

func (f *FileManager) WriteFile(text string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.File = append(f.File, text)
	f.cursor++
	f.isUpdated = true
}

func (f *FileManager) WriteRandomLogs() {
	count := 0

	for {
		logMessage := fmt.Sprintf("Log Message # %d \n", count)
		f.WriteFile(logMessage)
		count++
		time.Sleep(1 * time.Second)
	}
}

func (f *FileManager) ReadLastLine() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	text := ""

	if f.isUpdated {
		text = f.File[f.cursor]
		f.isUpdated = false
	}

	return text
}

func (f *FileManager) ReadNLastLines(n int) []string {
	f.mu.Lock()
	defer f.mu.Unlock()

	var result []string
	tempCursor := f.cursor

	for tempCursor >= 0 && n > 0 {
		result = append(result, f.File[tempCursor])
		tempCursor--
		n--
	}

	return result
}
