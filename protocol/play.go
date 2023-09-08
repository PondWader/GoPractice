package protocol

import "github.com/PondWader/GoPractice/protocol/packets"

func (client *ProtocolClient) play() {
	client.state = "play"
}

func (client *ProtocolClient) BeginPacketReader() {
	for {
		packetId, data, err := client.readPacket()
		if err != nil {
			return
		}

		var packetFormat interface{}
		switch packetId {
		case packets.SKeepAliveId:
			packetFormat = &packets.KeepAlivePacket{}
		case packets.SChatId:
			packetFormat = &packets.SChatPacket{}
		case packets.SPlayerPositionId:
			packetFormat = &packets.SPlayerPositionPacket{}
		case packets.SPlayerLookId:
			packetFormat = &packets.SPlayerLookPacket{}
		case packets.SPlayerPositionAndLookId:
			packetFormat = &packets.SPlayerPositionAndLookPacket{}
		case packets.SPlayerBlockPlacementId:
			packetFormat = &packets.SPlayerBlockPlacement{}
		case packets.SPlayerDiggingId:
			packetFormat = &packets.SPlayerDigging{}
		case packets.SCreativeInventoryActionId:
			packetFormat = &packets.SCreaviteInventoryAction{}
		default:
			// utils.Error("Received unrecognized packet of ID", packetId, "from", client.Username)
			continue
		}

		err = client.deserialize(data, packetFormat)
		if err != nil {
			break
		}
		client.HandlePacket(packetFormat)
	}
}
