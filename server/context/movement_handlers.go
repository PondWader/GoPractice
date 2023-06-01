package context

import (
	"github.com/PondWader/GoPractice/protocol"
)

func (p *ContextPlayer) sendToPlayersInView(packetId int, data []byte) {
	for _, entity := range p.EntitiesInView {
		if entity.Type() == "player" {
			player := entity.(*ContextPlayer)
			go func() {
				player.Mu.Lock()
				player.Client.WritePacket(packetId, data)
				player.Mu.Unlock()
			}()
		}
	}
}

func (p *ContextPlayer) handlePlayerPositionUpdate(packet interface{}) {
	playerPositionPacket := packet.(*protocol.SPlayerPositionPacket)
	p.Mu.Lock()

	diffX := playerPositionPacket.X - p.Position.X
	diffY := playerPositionPacket.Y - p.Position.Y
	diffZ := playerPositionPacket.Z - p.Position.Z

	oldX := p.Position.X
	oldY := p.Position.Y
	oldZ := p.Position.Z

	p.handlePositionChange(playerPositionPacket.X, playerPositionPacket.Y, playerPositionPacket.Z)

	GroundStateIsSame := p.IsOnGround == playerPositionPacket.OnGround

	if diffX == 0 && diffY == 0 && diffZ == 0 && GroundStateIsSame {
		p.Mu.Unlock()
		return
	}

	MovementIsWithinLimits := diffX < 4 && diffY < 4 && diffZ < 4 && diffX > -4 && diffY > -4 && diffZ > -4

	if GroundStateIsSame && MovementIsWithinLimits {
		p.sendToPlayersInView(0x15, protocol.Serialize(&protocol.CEntityRelativeMovePacket{
			EntityID: int(p.EntityId),
			DeltaX:   int8(playerPositionPacket.X*32) - int8(oldX*32),
			DeltaY:   int8(playerPositionPacket.Y*32) - int8(oldY*32),
			DeltaZ:   int8(playerPositionPacket.Z*32) - int8(oldZ*32),
			OnGround: playerPositionPacket.OnGround,
		}))
	} else {
		p.IsOnGround = playerPositionPacket.OnGround
		p.sendToPlayersInView(0x18, protocol.Serialize(&protocol.CEntityTeleportPacket{
			EntityID: int(p.EntityId),
			X:        int32(p.Position.X * 32),
			Y:        int32(p.Position.Y * 32),
			Z:        int32(p.Position.Z * 32),
			Yaw:      p.Position.GetYawAngle(),
			Pitch:    p.Position.GetPitchAngle(),
			OnGround: playerPositionPacket.OnGround,
		}))
	}

	p.Mu.Unlock()
}

func (p *ContextPlayer) handlePlayerLookUpdate(packet interface{}) {
	p.Mu.Lock()
	playerLookUpdatePacket := packet.(*protocol.SPlayerLookPacket)

	p.handleDirectionChange(playerLookUpdatePacket.Yaw, playerLookUpdatePacket.Pitch)

	p.sendToPlayersInView(0x16, protocol.Serialize(&protocol.CEntityLookPacket{
		EntityID: int(p.EntityId),
		Yaw:      p.Position.GetYawAngle(),
		Pitch:    p.Position.GetPitchAngle(),
		OnGround: playerLookUpdatePacket.OnGround,
	}))

	p.IsOnGround = playerLookUpdatePacket.OnGround
	p.Mu.Unlock()
}

func (p *ContextPlayer) handlePlayerPositionAndLookUpdate(packet interface{}) {
	playerPositionAndLookUpdatePacket := packet.(*protocol.SPlayerPositionAndLookPacket)
	p.Mu.Lock()

	diffX := playerPositionAndLookUpdatePacket.X - p.Position.X
	diffY := playerPositionAndLookUpdatePacket.Y - p.Position.Y
	diffZ := playerPositionAndLookUpdatePacket.Z - p.Position.Z

	oldX := p.Position.X
	oldY := p.Position.Y
	oldZ := p.Position.Z

	p.handlePositionChange(playerPositionAndLookUpdatePacket.X, playerPositionAndLookUpdatePacket.Y, playerPositionAndLookUpdatePacket.Z)
	p.handleDirectionChange(playerPositionAndLookUpdatePacket.Yaw, playerPositionAndLookUpdatePacket.Pitch)

	GroundStateIsSame := p.IsOnGround == playerPositionAndLookUpdatePacket.OnGround
	MovementIsWithinLimits := diffX < 4 && diffY < 4 && diffZ < 4 && diffX > -4 && diffY > -4 && diffZ > -4

	if GroundStateIsSame && MovementIsWithinLimits {
		p.sendToPlayersInView(0x17, protocol.Serialize(&protocol.CEntityLookAndRelativeMovePacket{
			EntityID: int(p.EntityId),
			DeltaX:   int8(playerPositionAndLookUpdatePacket.X*32) - int8(oldX*32),
			DeltaY:   int8(playerPositionAndLookUpdatePacket.Y*32) - int8(oldY*32),
			DeltaZ:   int8(playerPositionAndLookUpdatePacket.Z*32) - int8(oldZ*32),
			Yaw:      p.Position.GetYawAngle(),
			Pitch:    p.Position.GetPitchAngle(),
			OnGround: playerPositionAndLookUpdatePacket.OnGround,
		}))
	} else {
		p.IsOnGround = playerPositionAndLookUpdatePacket.OnGround
		p.sendToPlayersInView(0x18, protocol.Serialize(&protocol.CEntityTeleportPacket{
			EntityID: int(p.EntityId),
			X:        int32(p.Position.X * 32),
			Y:        int32(p.Position.Y * 32),
			Z:        int32(p.Position.Z * 32),
			Yaw:      p.Position.GetYawAngle(),
			Pitch:    p.Position.GetPitchAngle(),
			OnGround: playerPositionAndLookUpdatePacket.OnGround,
		}))
	}

	p.Mu.Unlock()
}

func (p *ContextPlayer) handlePositionChange(newX float64, newY float64, newZ float64) {
	oldX := p.Position.X
	oldZ := p.Position.Z

	p.Position.SetPos(newX, newY, newZ)

	// Detect if the player has changed chunk if so new chunks need to be loaded and new entities need to be displayed
	newChunkX := p.Position.GetBlockX() >> 4
	newChunkZ := p.Position.GetBlockZ() >> 4
	if int32(oldX)>>4 != newChunkX || int32(oldZ)>>4 != newChunkZ {
		chunk := p.Context.World.GetChunk(newChunkX, newChunkZ)
		p.Context.World.GetChunk(p.currentChunk.X, p.currentChunk.Z).RemoveEntity(p.EntityId)
		chunk.AddEntity(p.EntityId, p)
		p.currentChunk = chunk.GetKey()
		p.Mu.Unlock()

		p.streamChunks()
		p.updateViewedEntities()

		p.Mu.Lock()
	}
}

func (p *ContextPlayer) handleDirectionChange(newYaw float32, newPitch float32) {
	p.Position.SetDirection(newYaw, newPitch)
}
