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
