package structs

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
