package inform

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Codec struct {
	// KeyBag contains a mapping of colon-separated MAC addresses to their AES
	// keys
	KeyBag map[string]string
}

func (c *Codec) Unmarshal(fp io.Reader) (*InformWrapper, error) {
	w := NewInformWrapper()

	var magic int32
	binary.Read(fp, binary.BigEndian, &magic)
	if magic != PROTOCOL_MAGIC {
		return nil, errors.New("Invalid magic number")
	}

	binary.Read(fp, binary.BigEndian, &w.Version)
	io.ReadFull(fp, w.MacAddr)
	binary.Read(fp, binary.BigEndian, &w.Flags)

	iv := make([]byte, 16)
	io.ReadFull(fp, iv)

	binary.Read(fp, binary.BigEndian, &w.DataVersion)

	var dataLen int32
	binary.Read(fp, binary.BigEndian, &dataLen)

	p := make([]byte, dataLen)
	io.ReadFull(fp, p)

	key, ok := c.KeyBag[w.FormattedMac()]
	if !ok {
		return nil, errors.New("No key found")
	}

	u, err := Decrypt(p, iv, key)
	if err != nil {
		return nil, err
	}

	w.Payload = u

	return w, nil
}

func (c *Codec) Marshal(msg *InformWrapper) ([]byte, error) {
	b := &bytes.Buffer{}
	payload := msg.Payload
	var iv []byte

	if msg.IsEncrypted() {
		key, ok := c.KeyBag[msg.FormattedMac()]
		if !ok {
			return nil, errors.New("No key found")
		}

		var err error
		payload, iv, err = Encrypt(payload, key)
		if err != nil {
			return nil, err
		}
	}

	binary.Write(b, binary.BigEndian, PROTOCOL_MAGIC)
	binary.Write(b, binary.BigEndian, msg.Version)
	b.Write(msg.MacAddr)
	binary.Write(b, binary.BigEndian, msg.Flags)
	b.Write(iv)
	binary.Write(b, binary.BigEndian, msg.DataVersion)
	binary.Write(b, binary.BigEndian, int32(len(payload)))
	b.Write(payload)

	return b.Bytes(), nil
}
