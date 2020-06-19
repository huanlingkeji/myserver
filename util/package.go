package util

import (
	"encoding/binary"
)

const (
	srcIdSize   = 1
	modeSize    = 1
	seqIdSize   = 4
	protoIdSize = 2
)

//协议包格式 最终发送数据会在包头再加上协议包的长度 4字节
//+---------------------+-----------------+------------------------------+--------+--------+----------+
//|模块来源:1(游戏/战斗)|模式:1(应答/推送)|编码类型:1(string/json/proto3)|序列号:4|协议id:2|数据内容:n|
//+---------------------+-----------------+------------------------------+--------+--------+----------+

//解包
func Unpack(msg []byte) (srcId int8, mode int8, seqId uint32, protoId int16, reqData []byte) {
	pos := 0
	srcId = int8(msg[pos])
	pos = pos + srcIdSize
	mode = int8(msg[pos])
	pos = pos + modeSize
	seqId = binary.BigEndian.Uint32(msg[pos : pos+seqIdSize])
	pos = pos + seqIdSize
	protoId = int16(binary.BigEndian.Uint16(msg[pos : pos+protoIdSize]))
	pos = pos + protoIdSize
	reqData = msg[pos:]
	return srcId, mode, seqId, protoId, reqData
}

func PacketPayload(srcId int8, mode int8, seqId int32, protoId int16, data []byte) []byte {
	writer := Writer()
	writer.WriteS8(srcId)
	writer.WriteS8(mode)
	writer.WriteU32(uint32(seqId))
	writer.WriteS16(protoId)
	if data == nil {
		return writer.Data()
	}
	writer.WriteRawBytes(data)
	return writer.Data()
}

func FinalPkg(req []byte) []byte {
	w := Writer()
	w.WriteU32(uint32(len(req))) //TODO 数据包超大怎么办？
	w.WriteRawBytes(req)
	return w.Data()
}
