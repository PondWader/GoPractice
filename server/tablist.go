package server

import (
	"github.com/PondWader/GoPractice/protocol"
	"github.com/PondWader/GoPractice/protocol/packets"
)

func (p *Player) loadPlayerList() {
	playerData := []*packets.PlayerListActionAddPlayer{}

	p.server.Mu.RLock()
	players := p.server.Players
	p.server.Mu.RUnlock()

	for _, player := range players {
		playerData = append(playerData, &packets.PlayerListActionAddPlayer{
			UUID:               player.client.Uuid,
			Name:               player.client.Username,
			NumberOfProperties: 1,
			Properties: []*packets.PlayerListPlayerProperties{{
				Name:  "textures",
				Value: player.client.Skin,
			}},
			GameMode: 2,
			Ping:     1,
		})
	}

	p.mu.Lock()
	p.client.WritePacket(packets.CPlayerListItemId, protocol.Serialize(&packets.CPLayerListItemPacket{
		Action:          0,
		NumberOfPlayers: len(playerData),
		Data:            playerData,
	}))

	p.client.WritePacket(packets.CPlayerListHeaderAndFooterId, protocol.Serialize(&packets.CPlayerListHeaderAndFooter{
		Header: packets.ChatComponent{
			Text: "§b§lGoPractice §8v" + p.server.Version + "\n",
		},
		Footer: packets.ChatComponent{
			Text: "\n§3github.com/PondWader/GoPractice",
		},
	}))
	p.mu.Unlock()
}

func (p *Player) addToPlayerlist() {
	p.server.BroadcastPacket(packets.CPlayerListItemId, protocol.Serialize(&packets.CPLayerListItemPacket{
		Action:          0,
		NumberOfPlayers: 1,
		Data: []*packets.PlayerListActionAddPlayer{{
			UUID:               p.client.Uuid,
			Name:               p.client.Username,
			NumberOfProperties: 1,
			Properties: []*packets.PlayerListPlayerProperties{{
				Name:  "textures",
				Value: p.client.Skin,
			}},
			GameMode: 2,
			Ping:     1,
		}},
	}))
}

func (p *Player) removeFromPlayerlist() {
	p.server.BroadcastPacket(packets.CPlayerListItemId, protocol.Serialize(&packets.CPLayerListItemPacket{
		Action:          4,
		NumberOfPlayers: 1,
		Data: []*packets.PlayerListActionRemovePlayer{{
			UUID: p.client.Uuid,
		}},
	}))
}
