package common

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

func ExistFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return true
	}
}

func WorkPath(path ...string) (work string) {
	work, _ = os.Getwd()
	joinPath := filepath.Join(path...)
	work = filepath.Join(work, joinPath)

	if !ExistDir(work) {
		err := os.MkdirAll(work, 0755)
		if err != nil {
			panic(
				fmt.Sprintf(
					"Could not create directory %s",
					err,
				),
			)
		}
	}
	return work
}

func GenPath(f ...string) (path string) {
	return filepath.Join(f...)
}

func ReadFilesWithCallback(directory, queryFilename string, callback func(filePath string) error) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	extractWeek := func(str string) string {
		re := regexp.MustCompile(`(第\d+周)`)
		match := re.FindStringSubmatch(str)
		if len(match) > 1 {
			return match[0]
		}
		return ""
	}

	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			filePath := filepath.Join(directory, filename)
			if extractWeek(filename) == queryFilename {
				return callback(filePath)
			}
		}
	}

	return errors.New("未找到")
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return ""
	}
	return value
}

func SetEnv(key, value string) {
	if len(value) == 0 {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}
}

// GetWeekly 获取当前星期
func GetWeekly() (weekly string) {
	// 获取当前周期
	now := time.Now()
	weekday := now.Weekday()

	weekdayMap := map[time.Weekday]string{
		time.Monday:    "星期一",
		time.Tuesday:   "星期二",
		time.Wednesday: "星期三",
		time.Thursday:  "星期四",
		time.Friday:    "星期五",
		time.Saturday:  "星期六",
		time.Sunday:    "星期日",
	}

	return weekdayMap[weekday]
}

// GetWeek 获取当前第几周
func GetWeek(startMon int) string {
	_, weeks := time.Now().ISOWeek()
	return strconv.Itoa(weeks - startMon)
}

func GenMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func FristRun() bool {
	checkFile := GenPath("cache", "data", "第1周.json")
	olderTime := 3 * time.Hour
	return !FileOlder(checkFile, olderTime)
}

func FileOlder(filePath string, t time.Duration) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	fileModTime := fileInfo.ModTime()
	currentTime := time.Now()
	timeDifference := currentTime.Sub(fileModTime)

	if timeDifference > t {
		return false
	} else {
		return true
	}
}
