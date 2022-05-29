package wire

import "github.com/chuccp/utils/io"

type Handshake struct {
	HandshakeType         HandshakeType
	Length                uint32
	Version               uint16
	Random                []byte
	SessionIdLength       uint8
	SessionId             []byte
	CipherSuiteLength     uint16
	CipherSuites          []byte
	CompressMethodsLength byte
	CompressMethods       []byte
	ExtensionsLength      uint16
	Extensions            map[ExtensionType]*Extension
}

func (h *Handshake) readSession(readStream *io.ReadStream) (err error) {
	h.SessionIdLength, err = readStream.ReadUint8()
	if err != nil {
		return err
	} else {
		if h.SessionIdLength == 0 {
			return nil
		}
		h.SessionId, err = readStream.ReadBytes(int(h.SessionIdLength))
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}
func (h *Handshake) readCipherSuites(readStream *io.ReadStream) (err error) {

	h.CipherSuiteLength, err = readStream.Read2Uint16()
	if err != nil {
		return err
	} else {
		if h.CipherSuiteLength == 0 {
			return nil
		}
		h.CipherSuites, err = readStream.ReadBytes(int(h.CipherSuiteLength))
		if err != nil {
			return err
		}
		return nil
	}
}
func (h *Handshake) readCompression(readStream *io.ReadStream) (err error) {
	h.CompressMethodsLength, err = readStream.ReadUint8()
	if err != nil {
		return err
	}
	h.CompressMethods, err = readStream.ReadBytes(int(h.CompressMethodsLength))
	if err != nil {
		return err
	}
	return nil
}
func (h *Handshake) readExtensions(readStream *io.ReadStream) (err error) {

	h.ExtensionsLength, err = readStream.Read2Uint16()

	bytes, err := readStream.ReadBytes(int(h.ExtensionsLength))
	if err != nil {
		return err
	}
	h.Extensions,err = ParseExtension(bytes)

	return err
}
func ParseHandshake(data []byte) (*Handshake, error) {
	readStream := io.NewReadBytesStream(data)
	b, err := readStream.ReadUint8()
	if err != nil {
		return nil, err
	} else {
		var handshake Handshake
		handshake.HandshakeType = HandshakeType(b)
		handshake.Length, err = readStream.Read3Uint32()
		if err != nil {
			return nil, err
		}
		handshake.Version, err = readStream.Read2Uint16()
		if err != nil {
			return nil, err
		}
		handshake.Random, err = readStream.ReadBytes(32)
		if err != nil {
			return nil, err
		}
		err = handshake.readSession(readStream)
		if err != nil {
			return nil, err
		}

		err = handshake.readCipherSuites(readStream)
		if err != nil {
			return nil, err
		}
		err = handshake.readCompression(readStream)
		if err != nil {
			return nil, err
		}
		err = handshake.readExtensions(readStream)
		if err != nil {
			return nil, err
		}
		return &handshake, err
	}
}
