package world

func CoordsInChunk(x int, y int, z int) (int, int, int) {
	return x & 0xf, y & 0xf, z & 0xf
}
