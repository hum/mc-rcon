package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

type packetType int32

const (
	commandPacket packetType = 2
	authPacket    packetType = 3
	authFailed    int32      = -1
)

type Client struct {
	address  string
	port     string
	password string
	conn     net.Conn
}

type header struct {
	Size       int32
	RequestID  int32
	PacketType packetType
}

func CreateClient(address string, port string) (Client, error) {
	if address == "" || port == "" {
		return Client{}, errors.New("Either of the required connection values (address, port) are empty.")
	}

	c := Client{address: address, port: port}
	c.address = address
	c.port = port

	conn, err := net.DialTimeout("tcp", c.address+":"+c.port, 10*time.Second)
	if err != nil {
		return Client{}, err
	}
	c.conn = conn

	return c, err
}

func (c *Client) SendCommand(command string) (string, error) {
	_, payload, err := c.sendPacket(commandPacket, []byte(command))
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

func (c *Client) sendPacket(t packetType, p []byte) (header, []byte, error) {
	packet := encodeRequest(t, p)
	_, err := c.conn.Write(packet)
	if err != nil {
		return header{}, nil, err
	}

	head, payload, err := decodeResponse(c.conn)
	if err != nil {
		return header{}, nil, err
	}

	return head, payload, nil
}

func encodeRequest(t packetType, packet []byte) []byte {
	var buffer bytes.Buffer

	length := int32(len(packet) + 10)
	padding := [2]byte{}

	binary.Write(&buffer, binary.LittleEndian, length)
	binary.Write(&buffer, binary.LittleEndian, int32(0))
	binary.Write(&buffer, binary.LittleEndian, t)
	binary.Write(&buffer, binary.LittleEndian, packet)
	binary.Write(&buffer, binary.LittleEndian, padding)

	return buffer.Bytes()
}

func decodeResponse(reader io.Reader) (header, []byte, error) {
	hdr := header{}
	err := binary.Read(reader, binary.LittleEndian, &hdr)
	if err != nil {
		return header{}, nil, err
	}

	payload := make([]byte, hdr.Size-8)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return header{}, nil, err
	}
	return hdr, payload, nil
}

func (c *Client) Authenticate(password string) error {
	head, _, err := c.sendPacket(authPacket, []byte(password))
	if err != nil {
		return err
	}

	if head.RequestID == authFailed {
		return errors.New("Wrong password. Authentication has failed.")
	}
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}
