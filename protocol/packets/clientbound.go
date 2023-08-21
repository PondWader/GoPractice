package packets

import "github.com/google/uuid"

const (
	CKeepAliveId PacketId = iota
	CJoinGameId
	CChatMessageId
	CTimeUpdateId
	CEntityEquipmentId
	CSpawnPositionId
	CUpdateHealthId
	CRespawnId
	CSetPlayerPositionAndLookId
	CHeldItemChangeId
	CUseBedId
	CEntityAnimationId
	CSpawnPlayerId
	CCollectItemId
	CSpawnObjectId
	CSpawnMobId
	CSpawnPaintingId
	CSpawnExperienceOrbId
	CEntityVelocityId
	CDestroyEntitiesId
	CEntityId
	CEntityRelativeMoveId
	CEntityLookId
	CEntityLookAndRelativeMoveId
	CEntityTeleportId
	CEntityHeadRotationId
	CEntityStatusId
	CAttachEntityId
	CEntityMetadataId
	CEntityEffectId
	CRemoveEntityEffectId
	CSetExperienceId
	CEntityPropertiesId
	CChunkDataId
	CMultiBlockChangeId
	CBlockChangeId
	CBlockActionId
	CBlockBreakAnimationId
	CBulkMapChunksId
	CExplosionId
	CParticleOrSoundId
	CSoundId
	CParticleId
	CChangeGameStateId
	CSpawnGlobalEntityId
	COpenWindowId
	CCloseWindowId
	CSetWindowSlotId
	CSetWindowSlotsId
	CWindowPropertyId
	CConfirmTransactionId
	CUpdateSignId
	CMapId
	CUpdateBlockEntityId
	COpenSignEditorId
	CStatisticsId
	CPlayerListItemId
	CPlayAbilitiesId
	CTabCompleteId
	CScoreboardObjectiveId
	CUpdateScoreId
	CDisplayScoreboardId
	CUpdateTeamId
	CPluginMessageId
	CDisconnectId
	CServerDifficultyId
	CCombatEventId
	CCameraId
	CWorldBorderId
	CTitleId
	CSetCompressionId
	CPlayerListHeaderAndFooterId
	CLoadResourcePackId
	CUpdateEntityNbtId
)

type CJoinGamePacket struct {
	EntityID         int32  `type:"Int"`
	GameMode         uint8  `type:"UnsignedByte"`
	Dimension        int8   `type:"Byte"`
	Difficulty       uint8  `type:"UnsignedByte"`
	MaxPlayers       uint8  `type:"UnsignedByte"`
	LevelType        string `type:"String"`
	ReducedDebugInfo bool   `type:"Boolean"`
}

type CChatMessage struct {
	Data     ChatComponent `type:"JSON"`
	Position int8          `type:"Byte"`
}

type CTimeUpdate struct {
	WorldAge  int64 `type:"Long"`
	TimeOfDay int64 `type:"Long"`
}

type CSetPlayerPositionAndLook struct {
	X     float64 `type:"Double"`
	Y     float64 `type:"Double"`
	Z     float64 `type:"Double"`
	Yaw   float32 `type:"Float"`
	Pitch float32 `type:"Float"`
	Flags int8    `type:"Byte"`
}

type CHeldItemChangePacket struct {
	Slot int8 `type:"Byte"`
}

type CPlayerAbilitiesPacket struct {
	Flags        int8    `type:"Byte"`
	FlyingSpeed  float32 `type:"Float"`
	WalkingSpeed float32 `type:"Float"`
}

// Player list item action data

// 0: add player
type PlayerListActionAddPlayer struct {
	UUID               *uuid.UUID                    `type:"UUID"`
	Name               string                        `type:"String"`
	NumberOfProperties int                           `type:"VarInt"`
	Properties         []*PlayerListPlayerProperties `type:"Array"`
	GameMode           int                           `type:"VarInt"`
	Ping               int                           `type:"VarInt"`
	HasDisplayName     bool                          `type:"Boolean"`
	DisplayName        ChatComponent                 `type:"JSON" if:"HasDisplayName"`
}

type PlayerListPlayerProperties struct {
	Name     string `type:"String"`
	Value    string `type:"String"`
	IsSigned bool   `type:"Boolean"`
}

// 1: update gamemode
type PlayerListActionUpdateGamemode struct {
	UUID     *uuid.UUID `type:"UUID"`
	GameMode int        `type:"VarInt"`
}

// 2: update latency
type PlayerListActionUpdateLatency struct {
	UUID *uuid.UUID `type:"UUID"`
	Ping int        `type:"VarInt"`
}

// 3: update display name
type PlayerListActionUpdateDisplayName struct {
	UUID           *uuid.UUID    `type:"UUID"`
	HasDisplayName bool          `type:"Boolean"`
	DisplayName    ChatComponent `type:"JSON" if:"HasDisplayName"`
}

// 4: remove player
type PlayerListActionRemovePlayer struct {
	UUID *uuid.UUID `type:"UUID"`
}

type CPLayerListItemPacket struct {
	Action          int         `type:"VarInt"`
	NumberOfPlayers int         `type:"VarInt"`
	Data            interface{} `type:"Array"`
}

type CPlayerListHeaderAndFooter struct {
	Header ChatComponent `type:"JSON"`
	Footer ChatComponent `type:"JSON"`
}

type CChunkData struct {
	ChunkX             int32  `type:"Int"`
	ChunkZ             int32  `type:"Int"`
	GroundUpContinuous bool   `type:"Boolean"`
	PrimaryBitMask     uint16 `type:"UnsignedShort"`
	Size               int    `type:"VarInt"`
	Data               []byte `type:"ByteArray"`
}

type CDestroyEntitiesPacket struct {
	Count     int         `type:"VarInt"`
	EntityIDs []*EntityID `type:"Array"`
}

type CSpawnPlayerPacket struct {
	EntityID    int        `type:"VarInt"`
	UUID        *uuid.UUID `type:"UUID"`
	X           float64    `type:"FixedPoint"`
	Y           float64    `type:"FixedPoint"`
	Z           float64    `type:"FixedPoint"`
	Yaw         uint8      `type:"UnsignedByte"`
	Pitch       uint8      `type:"UnsignedByte"`
	CurrentItem int16      `type:"Short"`
	Metadata    []byte     `type:"ByteArray"`
}

type CEntityRelativeMovePacket struct {
	EntityID int  `type:"VarInt"`
	DeltaX   int8 `type:"Byte"`
	DeltaY   int8 `type:"Byte"`
	DeltaZ   int8 `type:"Byte"`
	OnGround bool `type:"Boolean"`
}

type CEntityLookPacket struct {
	EntityID int   `type:"VarInt"`
	Yaw      uint8 `type:"UnsignedByte"`
	Pitch    uint8 `type:"UnsignedByte"`
	OnGround bool  `type:"Boolean"`
}

type CEntityLookAndRelativeMovePacket struct {
	EntityID int   `type:"VarInt"`
	DeltaX   int8  `type:"Byte"`
	DeltaY   int8  `type:"Byte"`
	DeltaZ   int8  `type:"Byte"`
	Yaw      uint8 `type:"UnsignedByte"`
	Pitch    uint8 `type:"UnsignedByte"`
	OnGround bool  `type:"Boolean"`
}

type CEntityTeleportPacket struct {
	EntityID int   `type:"VarInt"`
	X        int32 `type:"Int"`
	Y        int32 `type:"Int"`
	Z        int32 `type:"Int"`
	Yaw      uint8 `type:"UnsignedByte"`
	Pitch    uint8 `type:"UnsignedByte"`
	OnGround bool  `type:"Boolean"`
}

type CEntityHeadRotationPacket struct {
	EntityID int   `type:"VarInt"`
	Yaw      uint8 `type:"UnsignedByte"`
}
