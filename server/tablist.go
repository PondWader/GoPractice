package server

import (
	"github.com/PondWader/GoPractice/protocol"
)

func (p *Player) loadPlayerList() {
	playerData := []*protocol.PlayerListActionAddPlayer{}

	p.server.Mu.RLock()
	players := p.server.Players
	p.server.Mu.RUnlock()

	for _, player := range players {
		playerData = append(playerData, &protocol.PlayerListActionAddPlayer{
			UUID:               player.client.Uuid,
			Name:               player.client.Username,
			NumberOfProperties: 1,
			Properties: []*protocol.PlayerListPlayerProperties{{
				Name:  "textures",
				Value: player.client.Skin,
			}},
			GameMode: 2,
			Ping:     1,
		})
	}

	p.mu.Lock()
	p.client.WritePacket(0x38, protocol.Serialize(&protocol.CPLayerListItemPacket{
		Action:          0,
		NumberOfPlayers: len(playerData),
		Data:            playerData,
	}))

	p.client.WritePacket(0x47, protocol.Serialize(&protocol.CPlayerListHeaderAndFooter{
		Header: protocol.ChatComponent{
			Text: "§b§lGoPractice §8v" + p.server.Version + "\n",
		},
		Footer: protocol.ChatComponent{
			Text: "\n§3github.com/PondWader/GoPractice",
		},
	}))
	p.mu.Unlock()
}

func (p *Player) addToPlayerlist() {
	p.server.BroadcastPacket(0x38, protocol.Serialize(&protocol.CPLayerListItemPacket{
		Action:          0,
		NumberOfPlayers: 1,
		Data: []*protocol.PlayerListActionAddPlayer{{
			UUID:               p.client.Uuid,
			Name:               p.client.Username,
			NumberOfProperties: 1,
			Properties: []*protocol.PlayerListPlayerProperties{{
				Name:  "textures",
				Value: p.client.Skin,
			}},
			GameMode: 2,
			Ping:     1,
		}},
	}))
}

func (p *Player) removeFromPlayerlist() {
	p.server.BroadcastPacket(0x38, protocol.Serialize(&protocol.CPLayerListItemPacket{
		Action:          4,
		NumberOfPlayers: 1,
		Data: []*protocol.PlayerListActionRemovePlayer{{
			UUID: p.client.Uuid,
		}},
	}))
}
