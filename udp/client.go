package udp

type Client struct {
	DesConnectionID []byte
}

func (c *Client) handshake() (err error) {
	c.DesConnectionID, err = GenerateConnectionID(8)
	if err != nil {
		return
	}
	return err
}
