package main

import (
	"sync-files/base"
	"sync-files/proto"
	"github.com/sotter/dovenet/protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sotter/dovenet/log"
	"os"
)

func SendFileInfo() {
	//只同步小于10M的文件
	files, err := base.LoadLocalFiles(g_root_path, func (file string) bool {
		info,_ := os.Stat(file)
		if info.Size() > 1024 * 1024 * 10 {
			return true
		} else {
			return false
		}
	})

	if err != nil {
		return
	}

	for key, value := range files {
		file_name := g_remote_path + "/" + key
		info := &sync_proto.SendFileInfo{
			FileName : &file_name,
			HashCode : &value,
		}

		body, err := proto.Marshal(info)
		if err != nil {
			log.Println("SendFileInfo Pb Marsal : ", err.Error())
			continue
		}

		req := protocol.NewCommMsg(uint16(sync_proto.SYNC_Msg_SendFileInfo),body)
		SendData("ctl_client", req)
	}
}

