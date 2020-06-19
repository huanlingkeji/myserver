package util

import (
	"encoding/binary"
	"haidao/backend/core/rpc"
)

const (
	srcIdSize   = 1
	modeSize    = 1
	seqIdSize   = 4
	protoIdSize = 2
)

//协议包格式 最终发送数据会在包头再加上协议包的长度 4字节
//+---------------------+-----------------+--------+--------+----------+
//|模块来源:1(游戏/战斗)|模式:1(应答/推送)|序列号:4|协议id:2|数据内容:n|
//+---------------------+-----------------+--------+--------+----------+

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

func DecodePacket(payload []byte) (*rpc.Msg_Front, error) {
	reader := packet.Reader(payload)
	//srcId
	srcId, err := reader.ReadS8()
	if err != nil {
		log.Error("read client srcId failed:", err)
		return nil, err
	}
	//mode
	mode, err := reader.ReadS8()
	if err != nil {
		log.Error("read client mode failed:", err)
		return nil, err
	}
	//read reserve
	reserve, err := reader.ReadU16()
	if err != nil {
		log.Error("read client reserve failed:", err)
		return nil, err
	}
	// 读客户端数据包序列号(1,2,3...)
	// 客户端发送的数据包必须包含一个自增的序号，必须严格递增
	// 加密后，可避免重放攻击-REPLAY-ATTACK
	seqId, err := reader.ReadU32()
	if err != nil {
		log.Error("read client sequence failed:", err)
		return nil, err
	}
	// 读协议号 2B
	protoId, err := reader.ReadS16()
	if err != nil {
		log.Error("read protocol number failed.")
		return nil, err
	}

	data, err := reader.ReadBytes(reader.Length() - reqDataStartIndex)
	if err != nil {
		log.Error("read data failed.")
		return nil, err
	}
	return &rpc.Msg_Front{
		SrcId:   int32(srcId),
		Mode:    int32(mode),
		Reserve: uint32(reserve),
		SeqId:   uint32(seqId),
		ProtoId: uint32(protoId),
		Data:    data,
	}, nil
}
