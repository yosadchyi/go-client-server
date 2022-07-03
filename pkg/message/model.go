package message

import (
	"encoding/json"
	"fmt"

	"github.com/yosadchyi/go-client-server/pkg/util"
)

type Operation string

const AddOp = Operation("Add")
const RemoveOp = Operation("Remove")
const GetItemOp = Operation("Get")
const GetAllItemsOp = Operation("GetAll")

// Base is a base for message
type Base struct {
	Operation Operation `json:"operation"`
}

// Add is a message representing addItem command
type Add struct {
	Base
	Data string `json:"data"`
}

// Remove is a message representing removeItem command
type Remove struct {
	Base
	ItemID int `json:"itemId"`
}

// Get is a message representing getItem command
type Get struct {
	Base
	ItemID int `json:"itemId"`
}

// GetAll is a message representing getAllItems command
type GetAll struct {
	Base
}

// Any represents any of valid messages, only one message field can be non-nil
type Any struct {
	Base
	Add         *Add
	Remove      *Remove
	GetItem     *Get
	GetAllItems *GetAll
}

func NewAdd(data string) Add {
	return Add{
		Base: Base{
			Operation: AddOp,
		},
		Data: data,
	}
}

func (m Add) ToJSON() *string {
	return util.ToJSON(m)
}

func NewRemove(itemID int) Remove {
	return Remove{
		Base: Base{
			Operation: RemoveOp,
		},
		ItemID: itemID,
	}
}

func (m Remove) ToJSON() *string {
	return util.ToJSON(m)
}

func NewGet(itemID int) Get {
	return Get{
		Base: Base{
			Operation: GetItemOp,
		},
		ItemID: itemID,
	}
}

func (m Get) ToJSON() *string {
	return util.ToJSON(m)
}

func NewGetAll() GetAll {
	return GetAll{
		Base: Base{
			Operation: GetAllItemsOp,
		},
	}
}

func (m GetAll) ToJSON() *string {
	return util.ToJSON(m)
}

func AnyFromJSON(data string) (*Any, error) {
	msg := Any{}
	var err error

	bytes := []byte(data)
	err = json.Unmarshal(bytes, &msg.Base)
	if err != nil {
		return nil, err
	}

	switch msg.Operation {
	case AddOp:
		msg.Add = &Add{}
		err = json.Unmarshal(bytes, msg.Add)
	case RemoveOp:
		msg.Remove = &Remove{}
		err = json.Unmarshal(bytes, msg.Remove)
	case GetItemOp:
		msg.GetItem = &Get{}
		err = json.Unmarshal(bytes, msg.GetItem)
	case GetAllItemsOp:
		msg.GetAllItems = &GetAll{}
		err = json.Unmarshal(bytes, msg.GetAllItems)
	default:
		err = fmt.Errorf("unrecognized operation %q", msg.Operation)
	}

	if err != nil {
		return nil, err
	}

	return &msg, err
}
