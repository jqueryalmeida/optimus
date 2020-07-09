package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails"
	"image/jpeg"
	"image/png"
	"sync"
)

type FileManager struct {
	Files   []*File
	Runtime *wails.Runtime
	Logger  *wails.CustomLogger
}

// NewFileManager creates a new FileManager.
func NewFileManager() *FileManager {
	return &FileManager{}
}

func (fm *FileManager) WailsInit(runtime *wails.Runtime) error {
	fm.Runtime = runtime
	fm.Logger = fm.Runtime.Log.New("FileManager")
	fm.Logger.Info("FileManager initialized...")
	return nil
}

// HandleFile processes a file from the client.
func (fm *FileManager) HandleFile(fileJson string) (err error) {
	file := &File{}
	if err := json.Unmarshal([]byte(fileJson), &file); err != nil {
		return err
	}

	if file.MimeType == "image/jpg" || file.MimeType == "image/jpeg" {
		file.Image, err = jpeg.Decode(bytes.NewReader(file.Data))
		fm.Files = append(fm.Files, file)
		fm.Logger.Info(fmt.Sprintf("added file to file manager: %s", file.Name))
	} else if file.MimeType == "image/png" {
		file.Image, err = png.Decode(bytes.NewReader(file.Data))
		fm.Files = append(fm.Files, file)
		fm.Logger.Info(fmt.Sprintf("added file to file manager: %s", file.Name))
	}
	if err != nil {
		return err
	}

	return nil
}

// Convert runs the conversion on all files in the FileManager.
func (fm *FileManager) Convert() (errs []error) {
	var wg sync.WaitGroup
	wg.Add(len(fm.Files))

	for _, file := range fm.Files {
		file := file
		go func(wg *sync.WaitGroup) {
			if err := file.Write(); err != nil {
				fm.Logger.Errorf("failed to convert file: %s", file.Name)
				errs = append(errs, fmt.Errorf("failed to convert file: %s", file.Name))
			}
			fm.Logger.Info(fmt.Sprintf("converted file: %s", file.Name))
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	return errs
}