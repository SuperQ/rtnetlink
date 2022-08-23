package rtnetlink

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/sys/unix"
)

func TestLinkMessageMarshalBinary(t *testing.T) {
	skipBigEndian(t)

	tests := []struct {
		name string
		m    Message
		b    []byte
		err  error
	}{
		{
			name: "empty",
			m:    &LinkMessage{},
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "no attributes",
			m: &LinkMessage{
				Family: 0,
				Type:   1,
				Index:  2,
				Flags:  0,
				Change: 0,
			},
			b: []byte{
				0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "no attributes with flags",
			m: &LinkMessage{
				Family: 0,
				Type:   1,
				Index:  2,
				Flags:  unix.IFF_UP,
				Change: 0,
			},
			b: []byte{
				0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "attributes",
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:   []byte{0x40, 0x41, 0x42, 0x43, 0x44, 0x45},
					Broadcast: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
					Name:      "lo",
				},
			},
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x07, 0x00, 0x03, 0x00, 0x6c, 0x6f, 0x00, 0x00,
				0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x05, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x0a, 0x00, 0x01, 0x00, 0x40, 0x41, 0x42, 0x43,
				0x44, 0x45, 0x00, 0x00, 0x0a, 0x00, 0x02, 0x00,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00,
			},
		},
		{
			name: "attributes ipip",
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:   []byte{10, 0, 0, 1},
					Broadcast: []byte{255, 255, 255, 255},
					Name:      "ipip",
				},
			},
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x09, 0x00, 0x03, 0x00, 0x69, 0x70, 0x69, 0x70,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x06, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x01, 0x00,
				0x0a, 0x00, 0x00, 0x01, 0x08, 0x00, 0x02, 0x00,
				0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			name: "info",
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:   []byte{0, 0, 0, 0, 0, 0},
					Broadcast: []byte{0, 0, 0, 0, 0, 0},
					Name:      "lo",
					Info: &LinkInfo{
						Kind:      "data",
						Data:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
						SlaveKind: "foo",
						SlaveData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
					},
				},
			},
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x07, 0x00, 0x03, 0x00, 0x6c, 0x6f, 0x00, 0x00,
				0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x05, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x02, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x38, 0x00, 0x12, 0x00, 0x09, 0x00, 0x01, 0x00,
				0x64, 0x61, 0x74, 0x61, 0x00, 0x00, 0x00, 0x00,
				0x0d, 0x00, 0x02, 0x00, 0x01, 0x02, 0x03, 0x04,
				0x05, 0x06, 0x07, 0x08, 0x09, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x04, 0x00, 0x66, 0x6f, 0x6f, 0x00,
				0x0d, 0x00, 0x05, 0x00, 0x01, 0x02, 0x03, 0x04,
				0x05, 0x06, 0x07, 0x08, 0x09, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "operational state",
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:          []byte{10, 0, 0, 1},
					Broadcast:        []byte{255, 255, 255, 255},
					Name:             "ipip",
					OperationalState: OperStateUp,
				},
			},
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x09, 0x00, 0x03, 0x00, 0x69, 0x70, 0x69, 0x70,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x06, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x01, 0x00,
				0x0a, 0x00, 0x00, 0x01, 0x08, 0x00, 0x02, 0x00,
				0xff, 0xff, 0xff, 0xff, 0x05, 0x00, 0x10, 0x00,
				0x06, 0x00, 0x00, 0x00,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.m.MarshalBinary()

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
			}
			if err != nil {
				return
			}

			if want, got := tt.b, b; !bytes.Equal(want, got) {
				t.Fatalf("unexpected Message bytes:\n- want: [%# x]\n-  got: [%# x]", want, got)
			}
		})
	}
}

func TestLinkMessageUnmarshalBinary(t *testing.T) {
	skipBigEndian(t)

	tests := []struct {
		name string
		b    []byte
		m    Message
		err  error
	}{
		{
			name: "empty",
			err:  errInvalidLinkMessage,
		},
		{
			name: "short",
			b:    make([]byte, 3),
			err:  errInvalidLinkMessage,
		},
		{
			name: "invalid attr",
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x04, 0x00, 0x01, 0x00, 0x04, 0x00, 0x02, 0x00,
				0x05, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x05, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			err: errInvalidLinkMessageAttr,
		},
		{
			name: "zero value",
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x02, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x05, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:   []byte{0, 0, 0, 0, 0, 0},
					Broadcast: []byte{0, 0, 0, 0, 0, 0},
				},
			},
		},
		{
			name: "no data",
			b: []byte{
				0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x02, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x05, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkMessage{
				Family: 0,
				Type:   1,
				Index:  2,
				Flags:  0,
				Change: 0,
				Attributes: &LinkAttributes{
					Address:   []byte{0, 0, 0, 0, 0, 0},
					Broadcast: []byte{0, 0, 0, 0, 0, 0},
				},
			},
		},
		{
			name: "data",
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x02, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x07, 0x00, 0x03, 0x00, 0x6c, 0x6f, 0x00, 0x00,
				0x08, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x05, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:   []byte{0, 0, 0, 0, 0, 0},
					Broadcast: []byte{0, 0, 0, 0, 0, 0},
					Name:      "lo",
				},
			},
		},
		{
			name: "attributes ipip",
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x01, 0x00, 0x0a, 0x00, 0x00, 0x01,
				0x08, 0x00, 0x02, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x09, 0x00, 0x03, 0x00, 0x69, 0x70, 0x69, 0x70,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x04, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x06, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:   []byte{10, 0, 0, 1},
					Broadcast: []byte{255, 255, 255, 255},
					Name:      "ipip",
				},
			},
		},
		{
			name: "info",
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x0a, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x02, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x07, 0x00, 0x03, 0x00, 0x6c, 0x6f, 0x00, 0x00,
				0x08, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x05, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x38, 0x00, 0x12, 0x00, 0x09, 0x00, 0x01, 0x00,
				0x64, 0x61, 0x74, 0x61, 0x00, 0x00, 0x00, 0x00,
				0x0d, 0x00, 0x02, 0x00, 0x01, 0x02, 0x03, 0x04,
				0x05, 0x06, 0x07, 0x08, 0x09, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x04, 0x00, 0x66, 0x6f, 0x6f, 0x00,
				0x0d, 0x00, 0x05, 0x00, 0x01, 0x02, 0x03, 0x04,
				0x05, 0x06, 0x07, 0x08, 0x09, 0x00, 0x00, 0x00,
			},
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:   []byte{0, 0, 0, 0, 0, 0},
					Broadcast: []byte{0, 0, 0, 0, 0, 0},
					Name:      "lo",
					Info: &LinkInfo{
						Kind:      "data",
						Data:      []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
						SlaveKind: "foo",
						SlaveData: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
					},
				},
			},
		},
		{
			name: "operational state",
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x01, 0x00, 0x0a, 0x00, 0x00, 0x01,
				0x08, 0x00, 0x02, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x09, 0x00, 0x03, 0x00, 0x69, 0x70, 0x69, 0x70,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x04, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x08, 0x00, 0x05, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x06, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x10, 0x00,
				0x06, 0x00, 0x00, 0x00,
			},
			m: &LinkMessage{
				Attributes: &LinkAttributes{
					Address:          []byte{10, 0, 0, 1},
					Broadcast:        []byte{255, 255, 255, 255},
					Name:             "ipip",
					OperationalState: OperStateUp,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LinkMessage{}
			err := (m).UnmarshalBinary(tt.b)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
			}
			if err != nil {
				return
			}

			if want, got := tt.m, m; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}

func TestLinkStatsUnmarshalBinary(t *testing.T) {
	skipBigEndian(t)

	tests := []struct {
		name string
		b    []byte
		m    *LinkStats
		err  error
	}{
		{
			name: "empty",
			err:  fmt.Errorf("incorrect size, want: 92 or 96"),
		},
		{
			name: "invalid < 92",
			b:    make([]byte, 91),
			err:  fmt.Errorf("incorrect size, want: 92 or 96"),
		},
		{
			name: "invalid > 96",
			b:    make([]byte, 97),
			err:  fmt.Errorf("incorrect size, want: 92 or 96"),
		},
		{
			name: "invalid > 92 < 96",
			b:    make([]byte, 93),
			err:  fmt.Errorf("incorrect size, want: 92 or 96"),
		},
		{
			name: "kernel <4.6",
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkStats{
				RXPackets:         0,
				TXPackets:         0,
				RXBytes:           0,
				TXBytes:           0,
				RXErrors:          0,
				TXErrors:          0,
				RXDropped:         0,
				TXDropped:         0,
				Multicast:         0,
				Collisions:        0,
				RXLengthErrors:    0,
				RXOverErrors:      0,
				RXCRCErrors:       0,
				RXFrameErrors:     0,
				RXFIFOErrors:      0,
				RXMissedErrors:    0,
				TXAbortedErrors:   0,
				TXCarrierErrors:   0,
				TXFIFOErrors:      0,
				TXHeartbeatErrors: 0,
				TXWindowErrors:    0,
				RXCompressed:      0,
				TXCompressed:      0,
			},
		},
		{
			name: "kernel 4.6+",
			b: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkStats{
				RXPackets:         0,
				TXPackets:         0,
				RXBytes:           0,
				TXBytes:           0,
				RXErrors:          0,
				TXErrors:          0,
				RXDropped:         0,
				TXDropped:         0,
				Multicast:         0,
				Collisions:        0,
				RXLengthErrors:    0,
				RXOverErrors:      0,
				RXCRCErrors:       0,
				RXFrameErrors:     0,
				RXFIFOErrors:      0,
				RXMissedErrors:    0,
				TXAbortedErrors:   0,
				TXCarrierErrors:   0,
				TXFIFOErrors:      0,
				TXHeartbeatErrors: 0,
				TXWindowErrors:    0,
				RXCompressed:      0,
				TXCompressed:      0,
				RXNoHandler:       0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LinkStats{}
			err := (m).unmarshalBinary(tt.b)
			if err != nil {
				if want, got := fmt.Sprintf("%s", tt.err), fmt.Sprintf("%s", err); want != got {
					t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
				}
				return
			}

			if want, got := tt.m, m; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}

func TestLinkStats64UnmarshalBinary(t *testing.T) {
	skipBigEndian(t)

	tests := []struct {
		name string
		b    []byte
		m    *LinkStats64
		err  error
	}{
		{
			name: "empty",
			err:  fmt.Errorf("incorrect size, want: 184 or 192 or 200"),
		},
		{
			name: "invalid < 184",
			b:    make([]byte, 183),
			err:  fmt.Errorf("incorrect size, want: 184 or 192 or 200"),
		},
		{
			name: "invalid > 184 < 192",
			b:    make([]byte, 188),
			err:  fmt.Errorf("incorrect size, want: 184 or 192 or 200"),
		},
		{
			name: "invalid > 192 < 200",
			b:    make([]byte, 193),
			err:  fmt.Errorf("incorrect size, want: 184 or 192 or 200"),
		},
		{
			name: "invalid > 200",
			b:    make([]byte, 201),
			err:  fmt.Errorf("incorrect size, want: 184 or 192 or 200"),
		},
		{
			name: "kernel <4.6",
			b: []byte{
				0x50, 0xb6, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0xc9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0xa9, 0x41, 0xcd, 0x09, 0x00, 0x00, 0x00, 0x00,
				0x96, 0x96, 0x2a, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkStats64{
				RXPackets:         0x1b650,
				TXPackets:         0xc906,
				RXBytes:           0x9cd41a9,
				TXBytes:           0x2a9696,
				RXErrors:          0x0,
				TXErrors:          0x0,
				RXDropped:         0x0,
				TXDropped:         0x0,
				Multicast:         0x0,
				Collisions:        0x0,
				RXLengthErrors:    0x0,
				RXOverErrors:      0x0,
				RXCRCErrors:       0x0,
				RXFrameErrors:     0x0,
				RXFIFOErrors:      0x0,
				RXMissedErrors:    0x0,
				TXAbortedErrors:   0x0,
				TXCarrierErrors:   0x0,
				TXFIFOErrors:      0x0,
				TXHeartbeatErrors: 0x0,
				TXWindowErrors:    0x0,
				RXCompressed:      0x0,
				TXCompressed:      0x0,
			},
		},
		{
			name: "kernel 4.6+",
			b: []byte{
				0x50, 0xb6, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0xc9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0xa9, 0x41, 0xcd, 0x09, 0x00, 0x00, 0x00, 0x00,
				0x96, 0x96, 0x2a, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkStats64{
				RXPackets:         0x1b650,
				TXPackets:         0xc906,
				RXBytes:           0x9cd41a9,
				TXBytes:           0x2a9696,
				RXErrors:          0x0,
				TXErrors:          0x0,
				RXDropped:         0x0,
				TXDropped:         0x0,
				Multicast:         0x0,
				Collisions:        0x0,
				RXLengthErrors:    0x0,
				RXOverErrors:      0x0,
				RXCRCErrors:       0x0,
				RXFrameErrors:     0x0,
				RXFIFOErrors:      0x0,
				RXMissedErrors:    0x0,
				TXAbortedErrors:   0x0,
				TXCarrierErrors:   0x0,
				TXFIFOErrors:      0x0,
				TXHeartbeatErrors: 0x0,
				TXWindowErrors:    0x0,
				RXCompressed:      0x0,
				TXCompressed:      0x0,
				RXNoHandler:       0x1,
			},
		},
		{
			name: "kernel 5.19+",
			b: []byte{
				0x50, 0xb6, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x06, 0xc9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0xa9, 0x41, 0xcd, 0x09, 0x00, 0x00, 0x00, 0x00,
				0x96, 0x96, 0x2a, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			m: &LinkStats64{
				RXPackets:          0x1b650,
				TXPackets:          0xc906,
				RXBytes:            0x9cd41a9,
				TXBytes:            0x2a9696,
				RXErrors:           0x0,
				TXErrors:           0x0,
				RXDropped:          0x0,
				TXDropped:          0x0,
				Multicast:          0x0,
				Collisions:         0x0,
				RXLengthErrors:     0x0,
				RXOverErrors:       0x0,
				RXCRCErrors:        0x0,
				RXFrameErrors:      0x0,
				RXFIFOErrors:       0x0,
				RXMissedErrors:     0x0,
				TXAbortedErrors:    0x0,
				TXCarrierErrors:    0x0,
				TXFIFOErrors:       0x0,
				TXHeartbeatErrors:  0x0,
				TXWindowErrors:     0x0,
				RXCompressed:       0x0,
				TXCompressed:       0x0,
				RXNoHandler:        0x1,
				RXOtherhostDropped: 0x2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LinkStats64{}
			err := (m).unmarshalBinary(tt.b)
			if err != nil {
				if want, got := fmt.Sprintf("%s", tt.err), fmt.Sprintf("%s", err); want != got {
					t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
				}
				return
			}

			if want, got := tt.m, m; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Message:\n- want: %#v\n-  got: %#v", want, got)
			}
		})
	}
}
