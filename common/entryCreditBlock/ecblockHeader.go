// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package entryCreditBlock

import (
	"encoding/json"
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

type ECBlockHeader struct {
	BodyHash            interfaces.IHash
	PrevHeaderHash      interfaces.IHash
	PrevFullHash        interfaces.IHash
	DBHeight            uint32
	HeaderExpansionArea []byte
	ObjectCount         uint64
	BodySize            uint64
}

var _ = fmt.Print
var _ interfaces.Printable = (*ECBlockHeader)(nil)

func (c *ECBlockHeader) Init() {
	if c.BodyHash == nil {
		c.BodyHash = primitives.NewZeroHash()
	}
	if c.PrevHeaderHash == nil {
		c.PrevHeaderHash = primitives.NewZeroHash()
	}
	if c.PrevFullHash == nil {
		c.PrevFullHash = primitives.NewZeroHash()
	}
	if c.HeaderExpansionArea == nil {
		c.HeaderExpansionArea = make([]byte, 0)
	}
}

func (e *ECBlockHeader) String() string {
	e.Init()
	var out primitives.Buffer
	out.WriteString(fmt.Sprintf("   %-20s %x\n", "ECChainID", e.GetECChainID().Bytes()[:3]))
	out.WriteString(fmt.Sprintf("   %-20s %x\n", "BodyHash", e.BodyHash.Bytes()[:3]))
	out.WriteString(fmt.Sprintf("   %-20s %x\n", "PrevHeaderHash", e.PrevHeaderHash.Bytes()[:3]))
	out.WriteString(fmt.Sprintf("   %-20s %x\n", "PrevFullHash", e.PrevFullHash.Bytes()[:3]))
	out.WriteString(fmt.Sprintf("   %-20s %d\n", "DBHeight", e.DBHeight))
	out.WriteString(fmt.Sprintf("   %-20s %x\n", "HeaderExpansionArea", e.HeaderExpansionArea))
	out.WriteString(fmt.Sprintf("   %-20s %d\n", "ObjectCount", e.ObjectCount))
	out.WriteString(fmt.Sprintf("   %-20s %d\n", "BodySize", e.BodySize))

	return (string)(out.DeepCopyBytes())
}

func (e *ECBlockHeader) SetBodySize(cnt uint64) {
	e.BodySize = cnt
}

func (e *ECBlockHeader) GetBodySize() uint64 {
	return e.BodySize
}

func (e *ECBlockHeader) SetObjectCount(cnt uint64) {
	e.ObjectCount = cnt
}

func (e *ECBlockHeader) GetObjectCount() uint64 {
	return e.ObjectCount
}

func (e *ECBlockHeader) SetHeaderExpansionArea(area []byte) {
	e.HeaderExpansionArea = area
}

func (e *ECBlockHeader) GetHeaderExpansionArea() (area []byte) {
	return e.HeaderExpansionArea
}

func (e *ECBlockHeader) SetBodyHash(prev interfaces.IHash) {
	e.BodyHash = prev
}

func (e *ECBlockHeader) GetBodyHash() interfaces.IHash {
	return e.BodyHash
}

func (e *ECBlockHeader) GetECChainID() interfaces.IHash {
	h := primitives.NewZeroHash()
	h.SetBytes(constants.EC_CHAINID)
	return h
}

func (e *ECBlockHeader) SetPrevHeaderHash(prev interfaces.IHash) {
	e.PrevHeaderHash = prev
}

func (e *ECBlockHeader) GetPrevHeaderHash() interfaces.IHash {
	return e.PrevHeaderHash
}

func (e *ECBlockHeader) SetPrevFullHash(prev interfaces.IHash) {
	e.PrevFullHash = prev
}

func (e *ECBlockHeader) GetPrevFullHash() interfaces.IHash {
	return e.PrevFullHash
}

func (e *ECBlockHeader) SetDBHeight(height uint32) {
	e.DBHeight = height
}

func (e *ECBlockHeader) GetDBHeight() (height uint32) {
	return e.DBHeight
}

func NewECBlockHeader() *ECBlockHeader {
	h := new(ECBlockHeader)
	h.Init()
	return h
}

func (e *ECBlockHeader) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *ECBlockHeader) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *ECBlockHeader) MarshalBinary() ([]byte, error) {
	e.Init()
	buf := primitives.NewBuffer(nil)

	// 32 byte ECChainID
	err := buf.PushBinaryMarshallable(e.GetECChainID())
	if err != nil {
		return nil, err
	}

	// 32 byte BodyHash
	err = buf.PushBinaryMarshallable(e.GetBodyHash())
	if err != nil {
		return nil, err
	}

	// 32 byte Previous Header Hash
	err = buf.PushBinaryMarshallable(e.GetPrevHeaderHash())
	if err != nil {
		return nil, err
	}

	// 32 byte Previous Full Hash
	err = buf.PushBinaryMarshallable(e.GetPrevFullHash())
	if err != nil {
		return nil, err
	}

	// 4 byte Directory Block Height
	err = buf.PushUInt32(e.GetDBHeight())
	if err != nil {
		return nil, err
	}

	// variable Header Expansion Size
	err = buf.PushVarInt(uint64(len(e.GetHeaderExpansionArea())))
	if err != nil {
		return nil, err
	}

	// varable byte Header Expansion Area
	err = buf.Push(e.GetHeaderExpansionArea())
	if err != nil {
		return nil, err
	}

	// 8 byte Object Count
	err = buf.PushUInt64(e.GetObjectCount())
	if err != nil {
		return nil, err
	}

	// 8 byte size of the Body
	err = buf.PushUInt64(e.GetBodySize())
	if err != nil {
		return nil, err
	}

	return buf.DeepCopyBytes(), nil
}

func (e *ECBlockHeader) UnmarshalBinaryData(data []byte) ([]byte, error) {
	e.Init()
	buf := primitives.NewBuffer(data)

	h := primitives.NewZeroHash()
	err := buf.PopBinaryMarshallable(h)
	if err != nil {
		return nil, err
	}
	if h.String() != "000000000000000000000000000000000000000000000000000000000000000c" {
		return nil, fmt.Errorf("Invalid ChainID - %s", h)
	}

	err = buf.PopBinaryMarshallable(e.BodyHash)
	if err != nil {
		return nil, err
	}
	err = buf.PopBinaryMarshallable(e.PrevHeaderHash)
	if err != nil {
		return nil, err
	}
	err = buf.PopBinaryMarshallable(e.PrevFullHash)
	if err != nil {
		return nil, err
	}

	e.DBHeight, err = buf.PopUInt32()
	if err != nil {
		return nil, err
	}

	// read the Header Expansion Area
	hesize, err := buf.PopVarInt()
	e.HeaderExpansionArea, err = buf.PopLen(int(hesize))
	if err != nil {
		return nil, err
	}

	e.ObjectCount, err = buf.PopUInt64()
	if err != nil {
		return nil, err
	}
	e.BodySize, err = buf.PopUInt64()
	if err != nil {
		return nil, err
	}

	return buf.DeepCopyBytes(), nil
}

func (e *ECBlockHeader) UnmarshalBinary(data []byte) error {
	_, err := e.UnmarshalBinaryData(data)
	return err
}

type ExpandedECBlockHeader ECBlockHeader

func (e ECBlockHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ExpandedECBlockHeader
		ChainID   string
		ECChainID string
	}{
		ExpandedECBlockHeader: ExpandedECBlockHeader(e),
		ChainID:               "000000000000000000000000000000000000000000000000000000000000000c",
		ECChainID:             "000000000000000000000000000000000000000000000000000000000000000c",
	})
}
