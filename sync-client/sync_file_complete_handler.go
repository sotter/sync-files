package main

import (
	"github.com/sotter/dovenet/protocol"
	"github.com/golang/protobuf/proto"
	"sync-files/proto"
	log "github.com/sotter/dovenet/log"
	"os"
)

func SyncFileCompleteHandler(service *ServiceClient, msg *protocol.CommMsg) error {
	req := &sync_proto.SyncFileComplete{}
	err := proto.Unmarshal(msg.Body, req)
	if err != nil {
		log.Println("SendFileInfo Unmarshal fail:", err.Error())
		return err
	}

	allFinish := g_file_manager.OnFileComplete(req.GetFileName())

	if allFinish {
		log.Println("All files trans complete !!!")
		os.Exit(0)
	}
	return nil
}
