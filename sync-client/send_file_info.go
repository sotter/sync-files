package main

import (
	"sync-files/base"
	"sync-files/proto"
	"github.com/sotter/dovenet/protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sotter/dovenet/log"
	"os"
	"path/filepath"
)

func Ignore(info os.FileInfo) bool {
	if info.Size() > 1024 * 1024 * 10 {
		return true
	} else {
		return false
	}
}

func LoadLocalFiles(path string) (error) {
	regulatedPath := filepath.ToSlash(path)
	loadMd5Sums := func(filePath string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if !info.IsDir() {
			p := base.RelativePath(regulatedPath, filepath.ToSlash(filePath))
			if info.Size() == 0 {
				log.Println(p, "size is 0, ignore!!!")
				return nil
			}

			if Ignore(info) {
				log.Println(p, "is ignored!!!")
				return nil
			}

			//只有非过滤的文件才同步
			hash_code := base.GetFileHashSum(filePath)
			if err := SendFileInfo(p, hash_code); err == nil {
				g_file_manager.OnScanNewFile(g_remote_path + "/" + p, FileInfo{
					timestamp : base.GetTimeMS(),
					size : uint32(info.Size()),
				})
			}
		}
		return nil
	}
	return filepath.Walk(path, loadMd5Sums)
}

func SendFileInfo(file_name string, hash_code uint64) error {
	file_name = g_remote_path + "/" + file_name
	info := &sync_proto.SyncFileInfo{
		FileName : &file_name,
		HashCode : &hash_code,
	}

	body, err := proto.Marshal(info)
	if err != nil {
		log.Println("SendFileInfo Pb Marsal : ", err.Error())
		return err
	}

	req := protocol.NewCommMsg(uint16(sync_proto.SYNC_Msg_SyncFileInfo), body)
	SendData("ctl_client", req)

	return nil
}

