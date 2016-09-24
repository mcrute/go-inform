package inform

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	PROTOCOL_MAGIC int32 = 1414414933 // UBNT
	INFORM_VERSION int32 = 0
	DATA_VERSION   int32 = 1

	ENCRYPTED_FLAG  = 1
	COMPRESSED_FLAG = 2
)

// Wrapper around an inform message, serializes directly into the wire
// protocol
type InformWrapper struct {
	Version     int32
	MacAddr     []byte
	Flags       int16
	DataVersion int32
	Payload     []byte
}

// Create InformWrapper with sane defaults
func NewInformWrapper() *InformWrapper {
	return &InformWrapper{
		Version:     INFORM_VERSION,
		MacAddr:     make([]byte, 6),
		Flags:       0,
		DataVersion: DATA_VERSION,
	}
}

// Update the payload data with JSON value
func (i *InformWrapper) UpdatePayload(v interface{}) error {
	if d, err := json.Marshal(v); err != nil {
		return err
	} else {
		i.Payload = d
		return nil
	}
}

// Format Mac address bytes as lowercase string with colons
func (i *InformWrapper) FormattedMac() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		i.MacAddr[0], i.MacAddr[1], i.MacAddr[2],
		i.MacAddr[3], i.MacAddr[4], i.MacAddr[5])
}

func (i *InformWrapper) String() string {
	b := &bytes.Buffer{}

	fmt.Fprintf(b, "Version:      \t%d\n", i.Version)
	fmt.Fprintf(b, "Mac Addr:     \t%s\n", i.FormattedMac())
	fmt.Fprintf(b, "Flags:        \t%d\n", i.Flags)
	fmt.Fprintf(b, " Encrypted:   \t%t\n", i.IsEncrypted())
	fmt.Fprintf(b, " Compressed:  \t%t\n", i.IsCompressed())
	fmt.Fprintf(b, "Data Version: \t%d\n", i.DataVersion)
	fmt.Fprintf(b, "Payload:      \t%q\n", i.Payload)

	return b.String()
}

func (i *InformWrapper) IsEncrypted() bool {
	return i.Flags&ENCRYPTED_FLAG != 0
}

func (i *InformWrapper) SetEncrypted(e bool) {
	if e {
		i.Flags |= ENCRYPTED_FLAG
	} else {
		i.Flags &= ENCRYPTED_FLAG
	}
}

func (i *InformWrapper) IsCompressed() bool {
	return i.Flags&COMPRESSED_FLAG != 0
}

func (i *InformWrapper) SetCompressed(c bool) {
	if c {
		i.Flags |= COMPRESSED_FLAG
	} else {
		i.Flags &= COMPRESSED_FLAG
	}
}
