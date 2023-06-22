package structs

import "math"

type Location struct {
	X     float64
	Y     float64
	Z     float64
	Yaw   float32
	Pitch float32
}

func (l *Location) Add(X float64, Y float64, Z float64) {
	l.X += X
	l.Y += Y
	l.Z += Z
}

func (l *Location) Set(X float64, Y float64, Z float64, Yaw float32, Pitch float32) {
	l.X = X
	l.Y = Y
	l.Z = Z
}

func (l *Location) SetPos(X float64, Y float64, Z float64) {
	l.X = X
	l.Y = Y
	l.Z = Z
}

func (l *Location) SetDirection(Yaw float32, Pitch float32) {
	l.Yaw = Yaw
	l.Pitch = Pitch
}

func (l *Location) GetBlockX() int32 {
	return int32(l.X)
}
func (l *Location) GetBlockY() int32 {
	return int32(l.Y)
}
func (l *Location) GetBlockZ() int32 {
	return int32(l.Z)
}

func (l *Location) GetYawAngle() uint8 {
	return uint8(((math.Mod(float64(l.Yaw), 360)) / 360) * 256)
}

func (l *Location) GetPitchAngle() uint8 {
	return uint8(((math.Mod(float64(l.Pitch), 360)) / 360) * 256)
}

func (l *Location) Clone() *Location {
	return &Location{l.X, l.Y, l.Z, l.Yaw, l.Pitch}
}

func LocationFromPositionFormat(val int64) *Location {
	return &Location{
		X: float64(val >> 38),
		Y: float64(val << 52 >> 52),
		Z: float64(val << 26 >> 38),
	}
}
