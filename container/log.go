package container

import (
	"fmt"
	"io"
	"os"
	"path"
)

func LogContainer(containerId string) {
	logFile := fmt.Sprintf("%s-json.log", containerId)
	logFileLocation := path.Join(fmt.Sprintf(InfoLocalFormat, containerId), logFile)

	file, err := os.Open(logFileLocation)
	if err != nil {
		fmt.Printf("error while open container log file: %v\n", err)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("error while read container log file: %v\n", err)
		return
	}

	_, err = fmt.Fprint(os.Stdout, string(content))
	if err != nil {
		fmt.Printf("error while print container log file: %v\n", err)
		return
	}
}