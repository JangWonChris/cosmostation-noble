package handler

import (
	"github.com/cosmostation/cosmostation-cosmos/mintscan/client"
	"github.com/cosmostation/cosmostation-cosmos/mintscan/db"
)

// Sessions is shorten for s will be used throughout this handler pakcage.
var (
	s              *Session
	ChainNumMap    = map[int]string{}
	ChainIDMap     = map[string]int{}
	ChainID        string
	MessageIDMap   = map[int]string{}
	MessageTypeMap = map[string]int{}
)

// Session is struct for wrapping both client and db structs.
type Session struct {
	Client *client.Client
	DB     *db.Database
}

// SetSession set Session object.
func SetSession(client *client.Client, db *db.Database) *Session {
	s = &Session{client, db}
	return s
}

func SetChainID() {
	ChainID, err := s.Client.RPC.GetNetworkChainID()
	if err != nil {
		panic(err)
	}

	chainInfo, err := s.DB.QueryChainInfo()
	if err != nil {
		panic(err)
	}
	for _, c := range chainInfo {
		ChainNumMap[int(c.ID)] = c.ChainID
		ChainIDMap[c.ChainID] = int(c.ID)
	}
	_, ok := ChainIDMap[ChainID]
	if !ok {
		panic("chain id does not exist")
	}
}

func SetMessageInfo() {
	messageInfo, err := s.DB.QueryMessageInfo()
	if err != nil {
		panic(err)
	}
	for _, m := range messageInfo {
		MessageIDMap[int(m.ID)] = m.Type
		MessageTypeMap[m.Type] = int(m.ID)
	}
}
