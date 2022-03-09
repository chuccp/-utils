package math

import "testing"

func TestName(t *testing.T) {

	t.Log(U32BE(BEU32(12)))

}
func Test2Name(t *testing.T) {

	t.Log(RandUInt32())

}