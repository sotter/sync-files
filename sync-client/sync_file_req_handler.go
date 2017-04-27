package main

import (
	"github.com/sotter/dovenet/protocol"
	log "github.com/sotter/dovenet/log"
	"sync-files/proto"
	"github.com/golang/protobuf/proto"
	sbase "github.com/sotter/dovenet/base"
	"io/ioutil"
	"os"
	"strings"
	"sync-files/base"
)

const  PACK_SIZE = (1024 * 10)

//收到发送文件的请求
func SyncFileReqHandler(service *ServiceClient, msg *protocol.CommMsg) error {
	req := &sync_proto.SyncFileInfoResp{}
	err := proto.Unmarshal(msg.Body, req)
	if err != nil {
		log.Println("SendFileInfo Unmarshal fail:", err.Error())
		return err
	}

	relative_name := req.GetFileName()
	if req.GetNeedSync() == false {
		log.Print(sbase.FormatOutput(relative_name, "Same", 60))
		return nil
	}

	//生成本地目录时，要把远程目录去掉
	local_file := GetRootPath() + "/" + strings.TrimPrefix(relative_name, g_remote_path + "/")
	hashCode := base.RSHash([]byte(local_file))

	//TODO: 增加判定，如果超过10M的文件，那么就不再传输， 过滤条件可以后面统一加进来；
	//获取文件权限
	info, err := os.Stat(local_file)
	if err != nil {
		log.Println("File Stat ", local_file, " fail:", err.Error())
		return nil
	}

	mode_perm := uint32(info.Mode().Perm())

	buf,  err := ioutil.ReadFile(local_file)
	buf_size := len(buf)

	if err != nil || buf_size <= 0 {
		log.Println("Read file ", local_file, " err:", err.Error())
		return nil
	}

	total_cnt := uint32(buf_size  - 1 + PACK_SIZE) / uint32(PACK_SIZE)
	current_cnt := uint32(0)
	for i := 0 ; i < int(total_cnt) ; i ++ {
		data := &sync_proto.SyncFileData{
			FileName : &relative_name,
			TotalPacks: &total_cnt,
			CurrentPacks: &current_cnt,
			FileMode: &mode_perm,
		}

		if (current_cnt < total_cnt - 1) {
			data.Data = buf[current_cnt * PACK_SIZE : (current_cnt+1) * PACK_SIZE]
		} else {
			data.Data = buf[current_cnt * PACK_SIZE:]
		}

		body , err := proto.Marshal(data)
		if err != nil {
			log.Println("SyncFileData Marshal:", err.Error())
			return nil
		}

		resp := protocol.NewCommMsg(uint16(sync_proto.SYNC_Msg_SyncFileData), body)
		service.Transport.SendDataWouldBlock("data_client", resp, hashCode)

		current_cnt = current_cnt + 1
	}
	return nil
}