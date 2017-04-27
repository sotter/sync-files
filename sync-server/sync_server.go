package main

import _ "net/http/pprof"
import (
	"github.com/sotter/dovenet/base"
	"github.com/sotter/dovenet/protocol"
	log "github.com/sotter/dovenet/log"
	"time"
	"net"
	"runtime"
	"net/http"
	"sync-files/proto"
	"github.com/BurntSushi/toml"
)

type SyncFileServer struct {
	TcpServer     *base.TCPServer
	ServerNumber  uint16
}

type Session struct {
	TcpConn *base.TcpConnection
	//其他扩展数据
}

/*
	SyncMsgType_SendFileInfoType SyncMsgType = 1
	SyncMsgType_SyncFileReqType  SyncMsgType = 2
	SyncMsgType_SyncFileDataType SyncMsgType = 3
 */
//直接使用Session的回调，可以减少一次从TcpConnection层到Session层的查找
func (this *Session) OnMessageData(conn *base.TcpConnection, msg protocol.Message) error{
	recv_msg := msg.(*protocol.CommMsg)
	if recv_msg.Header.MsgType == uint16(sync_proto.SYNC_Msg_SyncFileInfo) {
		return SendFileInfoHandler(this, recv_msg)
	} else if recv_msg.Header.MsgType == uint16(sync_proto.SYNC_Msg_SyncFileData){
		return SyncFileDataHandler(this, recv_msg)
	}

	return nil
}

func (this *Session)OnConnection(conn *base.TcpConnection) {
	log.Print("New connection from ", conn.Address)
}

func (this *Session)OnDisConnection(conn *base.TcpConnection) {
	log.Print("DisConnection : ", conn.Address)
}

func (this *SyncFileServer)SendData(connId uint64, msg protocol.Message) {
	tcpConn := this.TcpServer.Manager.GetSession(connId)
	tcpConn.Write(msg)
}

func (this *SyncFileServer)GetServerNumber() uint16 {
	//log.Info("Get Server Number : ", this.ServerNumber)
	return this.ServerNumber
}

func (this *SyncFileServer)Loop() {
	//单独绑定一个线程
	defer base.RecoverPrint()
	runtime.LockOSThread()

	for {
		tcpConn, err := this.TcpServer.Accept()
		if err != nil {
			log.Print("Accept ", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}

		//如果4分钟没有任何数据，断开连接
		//tcpConn.SetDeadline(4 * time.Minute)
		session := Session{
		}

		tcpConn.SetReadDeadline(time.Now().Add(time.Minute*4))

		tcpConnection:= base.NewServerConn(
			base.GetNetId(),
			this.TcpServer.Protocol.NewCodec(tcpConn.(* net.TCPConn)),
			&session,
			this.TcpServer.Manager)

		tcpConnection.Address = tcpConn.RemoteAddr().String()
		tcpConnection.NetworkCB = &session
		session.TcpConn = tcpConnection
		tcpConnection.ConnManager  = this.TcpServer.Manager

		tcpConnection.Start()
	}
}

var g_root_path = "./sync-files"

func GetRootPath() string {
	return g_root_path
}

func main() {
	go func() {
		log.Print(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	var Config struct {
		SynServer string
		SaveRootPath string
	}

	_, err := toml.DecodeFile("sync_server.toml", &Config)
	if err != nil {
		log.Println("Decode Config fail : ", err.Error())
		return
	}

	g_root_path = Config.SaveRootPath
	tcp_server, err := base.NewTCPServer(Config.SynServer, &protocol.CommProtocol{})

	if err != nil {
		log.Println("New TcpServer err:", err.Error())
		return
	}

	sync_server := SyncFileServer{
		ServerNumber : 1234,
		TcpServer    : tcp_server,
	}

	sync_server.Loop()
}