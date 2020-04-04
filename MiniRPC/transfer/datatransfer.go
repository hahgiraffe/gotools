/*
 * @Author: haha_giraffe
 * @Date: 2019-12-29 21:19:10
 * @Description: 传输数据格式接口
 */
package transfer

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"net"
)

type Transport struct {
	conn net.Conn
}

type RPCdata struct {
	Name string
	Args []interface{}
	Err  string
}

func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

//向t发送数据（记得加上四字节长度）
func (t *Transport) Send(data []byte) error {
	//需要多四个字节表示消息的长度，剩下的填充消息
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	//发送
	_, err := t.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

//从t接收数据（记得解析四字节长度）
func (t *Transport) Read() ([]byte, error) {
	header := make([]byte, 4)
	//ReadFull是将header填满
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RPCdata ==> []byte
func Encode(data RPCdata) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// []byte ==> RPCdata
func Decode(b []byte) (RPCdata, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data RPCdata
	if err := decoder.Decode(&data); err != nil {
		return RPCdata{}, err
	}
	return data, nil

}
