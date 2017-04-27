package main

import (
	"github.com/sotter/dovenet/protocol"
	sbase "github.com/sotter/dovenet/base"
	"github.com/golang/protobuf/proto"
	"sync-files/proto"
	log "github.com/sotter/dovenet/log"
	"sync-files/base"
)

func SendFileInfoHandler(session *Session, msg *protocol.CommMsg) error{
	// 比对下本地文件的看是否需要比对
	req := &sync_proto.SyncFileInfo{}
	err := proto.Unmarshal(msg.Body, req)
	if err != nil {
		log.Println("SendFileInfo Unmarshal fail:", err.Error())
		return err
	}

	file_name := req.GetFileName()
	local_hash := base.GetFileHashSum(GetRootPath() + "/" + file_name)
	remote_hash := req.GetHashCode()

	need_sync := true
	if local_hash== remote_hash {
		log.Print(sbase.FormatOutput(file_name, "Same", 60))
		need_sync = false
	}

	log.Print(sbase.FormatOutput(file_name, "Diff", 60))

	//文件不相同，告知客户端发送数据
	sync_file_req := &sync_proto.SyncFileInfoResp{
		FileName : &file_name,
		NeedSync : &need_sync,
	}

	body, err := proto.Marshal(sync_file_req)
	if err != nil {
		log.Println("sync_file_req Marshal fail:", err.Error())
		return err
	}

	resp := protocol.NewCommMsg(uint16(sync_proto.SYNC_Msg_SyncFileInfoResp), body)
	session.TcpConn.WriteWouldBlock(resp)

	return nil
}
