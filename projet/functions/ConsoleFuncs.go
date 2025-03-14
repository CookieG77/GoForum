package functions

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var shouldLogDebug = false
var logFolder = "logs"
var defaultLogTime = 3600
var loggerInitialized = false
var currentLogFile *os.File

// ClearCmd clear the console
func ClearCmd() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	}
}

// InitLogger initialize the logger
// If 'latest.log' file does not exist, it will be created
// If it exists, it will be renamed to 'log1.log' and the other log files will be renamed to 'log2.log', 'log3.log' ...
// Also, set up a goroutine to change the log file every 'LOG_FILE_CHANGE_TIME' seconds (default 1 hour)
func InitLogger() {
	if !loggerInitialized {
		InfoPrintln("Initialising the logger")
		// Check if the logs directory exists
		// If it does, check if empty, if not compress the logs in a zip file and create a new 'latest.log' file
		// If not, create it
		if _, err := os.Stat("logs"); os.IsNotExist(err) {
			err := os.Mkdir("logs", os.ModePerm)
			if err != nil {
				DebugPrintf("InitLogger mkdir failed -> %s\n", err)
				return
			}
		} else {
			// Check if the logs directory is empty or not
			// If not, compress the logs in a zip file
			// If it is, create a new 'latest.log' file
			files, err := os.ReadDir("logs")
			if err != nil {
				DebugPrintf("InitLogger readDir failed -> %s\n", err)
				return
			}
			// If there are remaining logs from previous launch, we compress them in a zip file
			if len(files) > 0 {
				CompressCurrentLogs()
			}
		}
		// Create the new 'latest.log' file
		newLogFile, err := os.Create("logs/latest.log")
		if err != nil {
			DebugPrintf("InitLogger create failed -> %s\n", err)
			return
		}
		currentLogFile = newLogFile
		DebugPrintln("New log file created")

		log.SetOutput(io.MultiWriter(os.Stdout, currentLogFile)) // Set the output of the logger to the console and the log file
		loggerInitialized = true
		var loggerResetTime int
		if val := os.Getenv("LOG_FILE_CHANGE_TIME"); val == "" {
			DebugPrintf("LOG_FILE_CHANGE_TIME not set, defaulting to %d seconds\n", defaultLogTime)
			loggerResetTime = defaultLogTime
		} else {
			val2, err := strconv.Atoi(val)
			if err != nil {
				DebugPrintf("LOG_FILE_CHANGE_TIME is not an int, defaulting to %d seconds\n", defaultLogTime)
				loggerResetTime = defaultLogTime
			} else {
				DebugPrintf("LOG_FILE_CHANGE_TIME set to %d\n", val2)
				loggerResetTime = val2
			}
		}
		// Create a goroutine to change the log file every 'LOG_FILE_CHANGE_TIME' seconds
		go func() {
			for {
				time.Sleep(time.Duration(loggerResetTime) * time.Second)
				SwitchLoggerFile()
			}
		}()
	}
}

// SwitchLoggerFile switch the logger file.
// Also close the current log file and rename it to an incremental name (latest.log -> log1.log -> log2.log ...)
func SwitchLoggerFile() {
	if loggerInitialized {
		err := currentLogFile.Close()
		if err != nil {
			DebugPrintf("SwitchLoggerFile close failed -> %s\n", err)
			return
		}
		// Get the list of files in the logs directory
		files, err := os.ReadDir("logs")
		if err != nil {
			DebugPrintf("SwitchLoggerFile readDir failed -> %s\n", err)
			return
		}

		// Sort files to ensure correct order
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() > files[j].Name()
		})

		// Rename the files
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "log") && strings.HasSuffix(file.Name(), ".log") {
				parts := strings.Split(file.Name(), ".")
				if len(parts) == 2 {
					numStr := strings.TrimPrefix(parts[0], "log")
					num, err := strconv.Atoi(numStr)
					if err == nil {
						newName := fmt.Sprintf("logs/log%d.log", num+1)
						err := os.Rename("logs/"+file.Name(), newName)
						if err != nil {
							DebugPrintf("SwitchLoggerFile rename failed -> %s\n", err)
							return
						}
					}
				}
			}
		}

		// Rename latest.log to log1.log
		err = os.Rename("logs/latest.log", "logs/log1.log")
		if err != nil {
			DebugPrintf("SwitchLoggerFile rename latest.log failed -> %s\n", err)
			return
		}

		// Create the new 'latest.log' file
		currentLogFile, err = os.Create("logs/latest.log")
		if err != nil {
			DebugPrintf("SwitchLoggerFile create failed -> %s\n", err)
			return
		}
		log.SetOutput(io.MultiWriter(os.Stdout, currentLogFile)) // Set the output of the logger to the console and the log file
		DebugPrintln("New log file created")
	}
}

// SetShouldLogDebug set the shouldLogDebug variable to the given value
func SetShouldLogDebug(value bool) {
	shouldLogDebug = value
}

// CompressCurrentLogs compress the current logs in a zip file
// The zip file will be named 'logs-YYYY-MM-DD-HH-MM-SS.zip'
// The logs will be removed after being compressed
func CompressCurrentLogs() {
	DebugPrintln("CompressCurrentLogs started")

	// Create the zip file
	zipAddress := fmt.Sprintf(
		"%s/logs-%s-%s.zip",
		logFolder,
		time.Now().Format("2006-01-02"),
		time.Now().Format("15-04-05"))

	zipFile, err := os.Create(zipAddress)
	if err != nil {
		DebugPrintf("CompressCurrentLogs create failed -> %s\n", err)
		return
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			DebugPrintf("CompressCurrentLogs close failed -> %s\n", err)
		}
	}(zipFile)

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			DebugPrintf("CompressCurrentLogs close failed -> %s\n", err)
		}
	}(zipWriter)

	// Get the list of files in the logs directory
	files, err := os.ReadDir(logFolder)
	if err != nil {
		DebugPrintf("CompressCurrentLogs readDir failed -> %s\n", err)
		return
	}

	// Add the files to the zip file, ignoring .zip files
	for _, file := range files {
		if !file.IsDir() && !strings.HasSuffix(file.Name(), ".zip") {
			err := addFileToZip(zipWriter, logFolder+"/"+file.Name())
			if err != nil {
				DebugPrintf("CompressCurrentLogs addFileToZip failed -> %s\n", err)
				return
			}
		}
	}

	// Remove the files, ignoring .zip files
	for _, file := range files {
		if !file.IsDir() && !strings.HasSuffix(file.Name(), ".zip") {
			err := os.Remove(fmt.Sprintf("%s/%s", logFolder, file.Name()))
			if err != nil {
				DebugPrintf("CompressCurrentLogs remove failed -> %s\n", err)
				return
			}
		}
	}
}

// AddFileToZip add the file at the given path to the given zip.Writer
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	DebugPrintln("addFileToZip started")
	// Check if the file is a directory
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		DebugPrintf("addFileToZip skipping directory %s\n", filePath)
		return nil
	}
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			DebugPrintf("addFileToZip close failed -> %s\n", err)
			return
		}
	}(file)

	// Creating the file header inside the ZIP
	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate // Use the same compression as the ZIP

	// Create the writer for the file inside the ZIP
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copy the file content to the ZIP
	_, err = io.Copy(writer, file)
	return err
}

// InfoPrintf print the given info message formatted with the given arguments
func InfoPrintf(format string, a ...interface{}) {
	log.Printf("\033[1;34m[Info] :\033[0m "+format, a...)
}

// InfoPrintln print the given info message with a new line
func InfoPrintln(s string) {
	log.Println("\033[1;34m[Info] :\033[0m " + s)
}

// ErrorPrintf print the given error message formatted with the given arguments
func ErrorPrintf(format string, a ...interface{}) {
	log.Printf("\033[1;31m[Error] :\033[0m "+format, a...)
}

// ErrorPrintln print the given error message with a new line
func ErrorPrintln(s string) {
	log.Println("\033[1;31m[Error] :\033[0m " + s)
}

// WarningPrintf print the given warning message formatted with the given arguments
func WarningPrintf(format string, a ...interface{}) {
	log.Printf("\033[1;33m[Warning] :\033[0m "+format, a...)
}

// WarningPrintln print the given warning message with a new line
func WarningPrintln(s string) {
	log.Println("\033[1;33m[Warning] :\033[0m " + s)
}

// SuccessPrintf print the given success message formatted with the given arguments
func SuccessPrintf(format string, a ...interface{}) {
	log.Printf("\033[1;32m[Success] :\033[0m "+format, a...)
}

// SuccessPrintln print the given success message with a new line
func SuccessPrintln(s string) {
	log.Println("\033[1;32m[Success] :\033[0m " + s)
}

// FatalPrintf print the given fatal message  formatted with the given arguments and exit the program
func FatalPrintf(format string, a ...interface{}) {
	log.Fatalf("\033[1;31m[Fatal] :\033[0m "+format, a...)
}

// FatalPrintln print the given fatal message with a new line and exit the program
func FatalPrintln(s string) {
	log.Fatalln("\033[1;31m[Fatal] :\033[0m " + s)
}

// DebugPrintf print the given debug message formatted with the given arguments
func DebugPrintf(format string, a ...interface{}) {
	if shouldLogDebug {
		log.Printf("\033[1;35m[Debug] :\033[0m "+format, a...)
	}
}

// DebugPrintln print the given debug message with a new line
func DebugPrintln(s string) {
	if shouldLogDebug {
		log.Println("\033[1;35m[Debug] :\033[0m " + s)
	}
}
