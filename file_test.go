package go_logger

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewAdapterFile(t *testing.T) {
	NewAdapterFile()
}

func TestAdapterFile_Name(t *testing.T) {
	fileAdapter := NewAdapterFile()

	if fileAdapter.Name() != FILE_ADAPTER_NAME {
		t.Error("file adapter name error")
	}
}

func TestAdapterFile_Write(t *testing.T) {

	fileAdapter := NewAdapterFile()

	fileConfig := &FileConfig{
		Filename:      "./test.log",
		LevelFileName: map[int]string{},
		MaxLine:       2000,
		MaxSize:       10000 * 4,
		JsonFormat:    true,
		DateSlice:     "d",
	}
	err := fileAdapter.Init(fileConfig)
	if err != nil {
		t.Fatal(err.Error())
	}

	loggerMsg := &loggerMessage{
		Timestamp:         time.Now().Unix(),
		TimestampFormat:   time.Now().Format("2006-01-02 15:04:05"),
		Millisecond:       time.Now().UnixNano() / 1e6,
		MillisecondFormat: time.Now().Format("2006-01-02 15:04:05.999"),
		Level:             LOGGER_LEVEL_DEBUG,
		LevelString:       "debug",
		Body:              "logger test file adapter write",
		File:              "file_test.go",
		Line:              50,
		Function:          "TestAdapterFile_Write",
	}
	err = fileAdapter.Write(loggerMsg)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestFileWriterRotatesExistingLogOnNewDay(t *testing.T) {
	directory, err := ioutil.TempDir("", "go-logger-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(directory)
	filename := filepath.Join(directory, "run.log")
	oldDate := time.Now().AddDate(0, 0, -1)
	line := oldDate.Format("01-02 15:04:05.000") + "[Info] old log\n"
	if err := ioutil.WriteFile(filename, []byte(line), 0644); err != nil {
		t.Fatal(err)
	}

	writer := NewFileWrite(filename)
	if err := writer.initFile(); err != nil {
		t.Fatal(err)
	}
	if writer.writer != nil {
		_ = writer.writer.Close()
	}

	archive := filepath.Join(directory, "run_"+oldDate.Format("20060102")+".log")
	if _, err := os.Stat(archive); err != nil {
		t.Fatalf("archive = %s: %v", archive, err)
	}
	if _, err := os.Stat(filename); err != nil {
		t.Fatalf("new active log = %s: %v", filename, err)
	}
}

func TestFileWriterKeepsCurrentDayLogOnRestart(t *testing.T) {
	directory, err := ioutil.TempDir("", "go-logger-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(directory)
	filename := filepath.Join(directory, "run.log")
	line := time.Now().Format("01-02 15:04:05.000") + "[Info] current log\n"
	if err := ioutil.WriteFile(filename, []byte(line), 0644); err != nil {
		t.Fatal(err)
	}

	writer := NewFileWrite(filename)
	if err := writer.initFile(); err != nil {
		t.Fatal(err)
	}
	if writer.writer != nil {
		_ = writer.writer.Close()
	}

	matches, err := filepath.Glob(filepath.Join(directory, "run_*.log"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("unexpected archives = %v", matches)
	}
}

func TestAdapterFile_WriteLevelFile(t *testing.T) {
	fileAdapter := NewAdapterFile()

	fileConfig := &FileConfig{
		Filename: "./test.log",
		LevelFileName: map[int]string{
			LOGGER_LEVEL_DEBUG: "./debug.log",
			LOGGER_LEVEL_INFO:  "./info.log",
			LOGGER_LEVEL_ERROR: "./error.log",
		},
		MaxLine:    2000,
		MaxSize:    10000 * 4,
		JsonFormat: true,
		DateSlice:  "d",
	}
	err := fileAdapter.Init(fileConfig)
	if err != nil {
		t.Fatal(err.Error())
	}

	loggerMsg := &loggerMessage{
		Timestamp:         time.Now().Unix(),
		TimestampFormat:   time.Now().Format("2006-01-02 15:04:05"),
		Millisecond:       time.Now().UnixNano() / 1e6,
		MillisecondFormat: time.Now().Format("2006-01-02 15:04:05.999"),
		Level:             LOGGER_LEVEL_DEBUG,
		LevelString:       "debug",
		Body:              "logger test file adapter write",
		File:              "file_test.go",
		Line:              50,
		Function:          "TestAdapterFile_Write",
	}
	fileAdapter.Write(loggerMsg)
	loggerMsg.Level = LOGGER_LEVEL_INFO
	fileAdapter.Write(loggerMsg)
	loggerMsg.Level = LOGGER_LEVEL_ERROR
	fileAdapter.Write(loggerMsg)
}
