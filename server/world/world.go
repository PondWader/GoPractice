package world

// The world package handles the loading & saving of world data

type World struct {
	chunks map[string]*Chunk
}

func New() *World {
	return &World{
		chunks: make(map[string]*Chunk),
	}
}
