package custom

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	mbltypes "github.com/cosmostation/mintscan-backend-library/types"

	//ibc
	interchainaccountstypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
)

const (
	// ibc (1)
	IBCTransferMsgTransfer = "ibctransfer/transfer"

	// IBC 02-client (4)
	IBCClientMsgCreateClient       = "ibcclient/create_client"
	IBCClientMsgUpdateClient       = "ibcclient/update_client"
	IBCClientMsgUpgradeClient      = "ibcclient/upgrade_client"
	IBCClientMsgSubmitMisbehaviour = "ibcclient/submit_misbehaviour"

	// IBC 03 connection (4)
	IBCConnectionMsgConnectionOpenInit    = "ibcconnection/connection_open_init"
	IBCConnectionMsgConnectionOpenTry     = "ibcconnection/connection_open_try"
	IBCConnectionMsgConnectionOpenAck     = "ibcconnection/connection_open_ack"
	IBCConnectionMsgConnectionOpenConfirm = "ibcconnection/connection_open_confirm"

	// IBC 04 channel (10)
	IBCChannelMsgChannelOpenInit     = "ibcchannel/channel_open_init"
	IBCChannelMsgChannelOpenTry      = "ibcchannel/channel_open_try"
	IBCChannelMsgChannelOpenAck      = "ibcchannel/channel_open_ack"
	IBCChannelMsgChannelOpenConfirm  = "ibcchannel/channel_open_confirm"
	IBCChannelMsgChannelCloseInit    = "ibcchannel/channel_close_init"
	IBCChannelMsgChannelCloseConfirm = "ibcchannel/channel_close_confirm"
	IBCChannelMsgRecvPacket          = "ibcchannel/recv_packet"
	IBCChannelMsgTimeout             = "ibcchannel/timeout"
	IBCChannelMsgTimeoutOnClose      = "ibcchannel/timeout_onclose"
	IBCChannelMsgAcknowledgement     = "ibcchannel/acknowledgement"
)

type txParser func(msg *sdktypes.Msg, txHash string) (msgType string, accounts []string)

var CustomTxParsers = make([]txParser, 0)

func init() {
	CustomTxParsers = append(CustomTxParsers, AccountExporterFromIBCMsg)
	CustomTxParsers = append(CustomTxParsers, AccountExporterFromCustomTxMsg)
}

func AccountExporterFromIBCMsg(msg *sdktypes.Msg, txHash string) (msgType string, accounts []string) {
	switch msg := (*msg).(type) {
	//ibc transfer (1)
	case *ibctransfertypes.MsgTransfer:
		msgType = IBCTransferMsgTransfer

	// ibc 02-client (4)
	case *ibcclienttypes.MsgCreateClient:
		msgType = IBCClientMsgCreateClient
	case *ibcclienttypes.MsgUpdateClient:
		msgType = IBCClientMsgUpdateClient
	case *ibcclienttypes.MsgUpgradeClient:
		msgType = IBCClientMsgUpgradeClient
	case *ibcclienttypes.MsgSubmitMisbehaviour:
		msgType = IBCClientMsgSubmitMisbehaviour

	// ibc 03 connection (4)
	case *ibcconnectiontypes.MsgConnectionOpenInit:
		msgType = IBCConnectionMsgConnectionOpenInit
	case *ibcconnectiontypes.MsgConnectionOpenTry:
		msgType = IBCConnectionMsgConnectionOpenTry
	case *ibcconnectiontypes.MsgConnectionOpenAck:
		msgType = IBCConnectionMsgConnectionOpenAck
	case *ibcconnectiontypes.MsgConnectionOpenConfirm:
		msgType = IBCConnectionMsgConnectionOpenConfirm

	// ibc 04 channel (10)
	case *ibcchanneltypes.MsgChannelOpenInit:
		msgType = IBCChannelMsgChannelOpenInit
	case *ibcchanneltypes.MsgChannelOpenTry:
		msgType = IBCChannelMsgChannelOpenTry
	case *ibcchanneltypes.MsgChannelOpenAck:
		msgType = IBCChannelMsgChannelOpenAck
	case *ibcchanneltypes.MsgChannelOpenConfirm:
		msgType = IBCChannelMsgChannelOpenConfirm
	case *ibcchanneltypes.MsgChannelCloseInit:
		msgType = IBCChannelMsgChannelCloseInit
	case *ibcchanneltypes.MsgChannelCloseConfirm:
		msgType = IBCChannelMsgChannelCloseConfirm
	case *ibcchanneltypes.MsgRecvPacket:
		msgType = IBCChannelMsgRecvPacket
		switch msg.Packet.SourcePort {
		case "transfer":
			var pd ibctransfertypes.FungibleTokenPacketData
			AppCodec.UnmarshalJSON(msg.Packet.GetData(), &pd)
			accounts = mbltypes.AddNotNullAccount(pd.Receiver)
		case "icahost":
			var pd interchainaccountstypes.InterchainAccountPacketData
			AppCodec.UnmarshalJSON(msg.Packet.GetData(), &pd)
			icaMsgs, err := interchainaccountstypes.DeserializeCosmosTx(EncodingConfig.Marshaler, pd.GetData())
			if err != nil {
				// TODO :
				// catch error
			}
			for i := range icaMsgs {
				msgType, accounts := mbltypes.AccountExporterFromCosmosTxMsg(&icaMsgs[i])
				for _, customTxParser := range CustomTxParsers {
					if msgType != "" {
						break
					}
					customMsgType, account := customTxParser(&icaMsgs[i], txHash)
					msgType = customMsgType
					accounts = append(accounts, account...)
				}
				if msgType == "" {
					// msgType 이 없을 경우, 해당 건은 수집하지 않는다.
					continue
				}
			}
		}
	case *ibcchanneltypes.MsgTimeout:
		msgType = IBCChannelMsgTimeout
		var pd ibctransfertypes.FungibleTokenPacketData
		AppCodec.UnmarshalJSON(msg.Packet.GetData(), &pd)
		accounts = mbltypes.AddNotNullAccount(pd.Sender)
	case *ibcchanneltypes.MsgTimeoutOnClose:
		msgType = IBCChannelMsgTimeoutOnClose
		var pd ibctransfertypes.FungibleTokenPacketData
		AppCodec.UnmarshalJSON(msg.Packet.GetData(), &pd)
		accounts = mbltypes.AddNotNullAccount(pd.Sender)
	case *ibcchanneltypes.MsgAcknowledgement:
		msgType = IBCChannelMsgAcknowledgement

	default:
		// 전체 case에서 msg를 찾지 못했을 때만, 로깅하도록 하기 위해 주석
		// zap.S().Infof("undefined message type in cosmos : %T, will search msg type in custom module", msg)
	}

	return
}
