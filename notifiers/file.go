package notifiers

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type FileNotifier struct {
	Path string
}

func (file *FileNotifier) Open(configPath string) bool {
	pathPath := fmt.Sprintf("%s.path", configPath)
	if !viper.IsSet(pathPath) {
		log.Println("Need 'path' for FileNotifier")
		return false
	}
	file.Path = viper.GetString(pathPath)

	f, err := os.OpenFile(file.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Could not open:", file.Path)
		return false
	}
	f.Close()
	log.Println("File notifier path: ", file.Path)
	return true
}

func (file *FileNotifier) Message(message string) bool {
	f, err := os.OpenFile(file.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Could not open:", file.Path)
		return false
	}
	defer f.Close()
	_, err = f.WriteString(message + "\n")
	if err != nil {
		log.Println("Could not write notification to:", file.Path)
		return false
	}
	return true
}

func (file *FileNotifier) Close() bool {
	return true
}
