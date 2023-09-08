package packets

const (
	SKeepAliveId PacketId = iota
	SChatId
	SUseEntityId
	SGroundStatusId
	SPlayerPositionId
	SPlayerLookId
	SPlayerPositionAndLookId
	SPlayerDiggingId
	SPlayerBlockPlacementId
	SHeldItemChangeId
	SArmSwingId
	SPlayerActionId
	SSteerVehicleId
	SCloseWindowId
	SClickWindowId
	SConfirmTransaction
	SCreativeInventoryActionId
	SEnchantItemId
	SUpdateSignId
	SPlayerAbilitiesId
	STabCompleteId
	SClientSettingsId
	SClientStatusId
	SPluginMessageId
	SSpectateId
	SResourcePackStatusId
)

type SChatPacket struct {
	Message string `type:"String"`
}

type SPlayerPositionPacket struct {
	X        float64 `type:"Double"`
	Y        float64 `type:"Double"`
	Z        float64 `type:"Double"`
	OnGround bool    `type:"Boolean"`
}

type SPlayerLookPacket struct {
	Yaw      float32 `type:"Float"`
	Pitch    float32 `type:"Float"`
	OnGround bool    `type:"Boolean"`
}

type SPlayerPositionAndLookPacket struct {
	X        float64 `type:"Double"`
	Y        float64 `type:"Double"`
	Z        float64 `type:"Double"`
	Yaw      float32 `type:"Float"`
	Pitch    float32 `type:"Float"`
	OnGround bool    `type:"Boolean"`
}

type SPlayerBlockPlacement struct {
	Location *Position `type:"Position"`
	Face     int8      `type:"Byte"`
}

type SPlayerDigging struct {
	Status   int8      `type:"Byte"`
	Location *Position `type:"Position"`
	Face     int8      `type:"Byte"`
}

type SCreaviteInventoryAction struct {
	Slot        int16 `type:"Short"`
	ClickedItem *Slot `type:"Struct"`
}
