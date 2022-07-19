package udp

type RawConn struct {

}

func (rc *RawConn) Read(data []byte)(n int, err error)  {
	return 0,nil
}
func (rc *RawConn) Write(p []byte) (n int, err error) {
	return 0,nil
}
func (rc *RawConn) push(p []byte) (n int, err error) {
	return 0,nil
}
func newRawConn() *RawConn {
	return &RawConn{}
}