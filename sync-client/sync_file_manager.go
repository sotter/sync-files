package main

import (
)
import (
	"sync"
	"sync-files/base"
	"sync/atomic"
	log "github.com/sotter/dovenet/log"
	sbase "github.com/sotter/dovenet/base"
	"fmt"
)

const sessionMapNum = 8
//怎么把一个超大的任务进行拆分？
type FileInfo struct {
	timestamp  uint64  // begin Sync Time
	size       uint32  // 文件大小
}

//key: requestid, value: 同步等待返回的wait
type FileMap struct {
	rwmutex sync.Mutex
	fileMap map[string](FileInfo)
}

type FileManager struct {
	ScanFinish    bool                   // 待传输的文件是否遍历完成
	fileCnt       int32                  // 当前文件的数目
	fileMapGroups [sessionMapNum]FileMap // 文件管理
	 //disposeFlag bool
	 //disposeOnce sync.Once
	 //disposeWait sync.WaitGroup
	current       uint64
}

func NewFileManager() *FileManager {
	manager := &FileManager{
		ScanFinish : false,
		fileCnt : 0,
	}

	for i := 0; i < sessionMapNum; i++ {
		manager.fileMapGroups[i].fileMap = make(map[string](FileInfo))
	}
	return manager
}

var mutex sync.Mutex

func (this *FileManager) OnScanNewFile(file_name string, info FileInfo) {
	hash := base.RSHash([]byte(file_name))

	smap := &this.fileMapGroups[hash % sessionMapNum]
	smap.rwmutex.Lock()
	_, exist := smap.fileMap[file_name]
	if !exist {
		//log.Println("fileManager insert : ", file_name)
		smap.fileMap[file_name] = info
		atomic.AddInt32(&this.fileCnt, 1)
	}
	smap.rwmutex.Unlock()
}

func (this *FileManager)OnFileComplete(file_name string) (allFinish bool) {
	allFinish = false
	hash := base.RSHash([]byte(file_name))

	smap := &this.fileMapGroups[hash % sessionMapNum]
	smap.rwmutex.Lock()
	info, exist := smap.fileMap[file_name]
	if exist {
		now := base.GetTimeMS()
		delete(smap.fileMap, file_name)
		atomic.AddInt32(&this.fileCnt, -1)
		str := fmt.Sprintf("complete %dms size:%dKB", int(now - info.timestamp),int(info.size/1000))
		log.Println(sbase.FormatOutput(file_name, str, 60))
	}

	if this.ScanFinish && atomic.AddInt32(&this.fileCnt, 0) == 0 {
		allFinish = true
	}
	smap.rwmutex.Unlock()

	return allFinish
}
