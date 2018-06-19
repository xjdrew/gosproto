package sproto_test

import (
	"bytes"
	"testing"

	"github.com/xjdrew/gosproto"
)

type PackTestCase struct {
	Name     string
	Unpacked []byte
	Origin   []byte
	Packed   []byte
}

var packTestCases []*PackTestCase = []*PackTestCase{
	&PackTestCase{
		Name:     "SimplePack",
		Origin:   []byte{0x08, 0x00, 0x00, 0x00, 0x03, 0x00, 0x02, 0x00, 0x19, 0x00, 0x00, 0x00, 0xaa, 0x01, 0x00, 0x00},
		Packed:   []byte{0x51, 0x08, 0x03, 0x02, 0x31, 0x19, 0xaa, 0x01},
		Unpacked: []byte{0x08, 0x00, 0x00, 0x00, 0x03, 0x00, 0x02, 0x00, 0x19, 0x00, 0x00, 0x00, 0xaa, 0x01, 0x00, 0x00},
	},
	&PackTestCase{
		Name: "FFPack",
		Origin: bytes.Join([][]byte{
			bytes.Repeat([]byte{0x8a}, 30),
			[]byte{0x00, 0x00},
		}, nil),
		Packed: bytes.Join([][]byte{
			[]byte{0xff, 0x03},
			bytes.Repeat([]byte{0x8a}, 30),
			[]byte{0x00, 0x00},
		}, nil),
		Unpacked: bytes.Join([][]byte{
			bytes.Repeat([]byte{0x8a}, 30),
			[]byte{0x00, 0x00},
		}, nil),
	},
	&PackTestCase{
		Name:     "FFPack2",
		Origin:   []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E},
		Packed:   []byte{0xFF, 0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x00, 0x00},
		Unpacked: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x00, 0x00},
	},
}

func TestPack(t *testing.T) {
	for _, tc := range packTestCases {
		packed := sproto.Pack(tc.Origin)
		if !bytes.Equal(packed, tc.Packed) {
			t.Log("packed:", packed)
			t.Log("expected:", tc.Packed)
			t.Fatalf("test case *%s* failed", tc.Name)
		}
	}
}

func TestUnpack(t *testing.T) {
	var allUnpacked, allPacked []byte
	for _, tc := range packTestCases {
		unpacked, err := sproto.Unpack(tc.Packed)
		if err != nil {
			t.Fatalf("test case *%s* failed with error:%s", tc.Name, err)
		}
		if !bytes.Equal(unpacked, tc.Unpacked) {
			t.Log("unpacked:", unpacked)
			t.Log("expected:", tc.Unpacked)
			t.Fatalf("test case *%s* failed", tc.Name)
		}
		allUnpacked = sproto.Append(allUnpacked, tc.Unpacked)
		allPacked = sproto.Append(allPacked, tc.Packed)
	}
	unpacked, err := sproto.Unpack(allPacked)
	if err != nil {
		t.Fatalf("test case *total* failed with error:%s", err)
	}
	if !bytes.Equal(unpacked, allUnpacked) {
		t.Log("unpacked:", unpacked)
		t.Log("expected:", allUnpacked)
		t.Fatal("test case *total* failed")
	}
}
