package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ConfigLogger() error {
	file, err := createLogFile()
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}
func createLogFile() (*os.File, error) {
	dirInfo, err := os.Stat("logfiles")
	if err != nil || (!dirInfo.IsDir()) {
		err = os.Mkdir("logfiles", 0755)
		if err != nil {
			return &os.File{}, err
		}
	}
	formattedDate := time.Now().Format("2006-01-02")
	currentDir, _ := os.Getwd()
	logPath := filepath.Join(currentDir, fmt.Sprintf("logfiles/%s.log", formattedDate))
	return os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
}
