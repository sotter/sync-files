package base

import (
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
	log "github.com/sotter/dovenet/log"
)

func GetFileHashSum(file_name string) uint64 {
	hash := uint64(0)
	buf,  err := ioutil.ReadFile(file_name)
	if err == nil {
		hash = RSHash(buf)
	}
	return hash
}

type IgnoreFunc func(file string) bool

func LoadLocalFiles(path string,  ignore IgnoreFunc) (map[string]uint64, error) {
	files := map[string]uint64{}
	regulatedPath := filepath.ToSlash(path)
	loadMd5Sums := func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			p := RelativePath(regulatedPath, filepath.ToSlash(filePath))
			if !ignore(filePath) {  //只有非过滤的文件才同步
				files[p] = GetFileHashSum(filePath)
			}
		}
		return nil
	}

	err := filepath.Walk(path, loadMd5Sums)
	if err != nil {
		return files, err
	}

	log.Println("Loaded ", len(files),  " files from ", path)
	return files, nil
}

func SaveDataToLocalFile(data []byte, dir string, name string, append bool, file_mode os.FileMode) error {
	//log.Println("Save Local file size:", len(data), " :", name)
	err := os.MkdirAll(dir, 0775)
	if err == nil || os.IsNotExist(err) {
		//log.Println("Create dir:", dir, " Success")
	} else {
		log.Println("Create dir:", dir, " fail ", err.Error())
		return err
	}

	flag := os.O_CREATE | os.O_RDWR
	if append {
		flag = flag | os.O_APPEND
	} else {
		flag = flag | os.O_TRUNC   //打开并清空文件
	}

	file, err := os.OpenFile(name, flag, file_mode)
	if err != nil {
		log.Println("Create file ", name, " fail:", err.Error())
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	file.Close()
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func RelativePath(path string, filePath string) string {
	if path == "." {
		return strings.TrimPrefix(filePath, "/")
	} else {
		return strings.TrimPrefix(strings.TrimPrefix(filePath, path), "/")
	}
}