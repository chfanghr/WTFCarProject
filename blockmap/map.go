package blockmap

type ID struct{ x, y int }

type block struct{ isBarrier bool }

type blockmap struct{ blocks [][]block }

func (i *ID) Init(x, y int) {
	if x > 0 && y > 0 {
		i.x, i.y = x, y
		return
	}
	i.x, i.y = 0, 0
}
func (i *ID) IncreaseX() { i.x++ }
func (i *ID) IncreaseY() { i.y++ }
func (i *ID) DecreaseX() {
	if i.x--; i.x < 0 {
		i.x = 0
	}
}
func (i *ID) DecreaseY() {
	if i.y--; i.y < 0 {
		i.y = 0
	}
}
func (i *ID) GetX() int             { return i.x }
func (i *ID) GetY() int             { return i.y }
func (i *ID) MoveUpDiagonally()     { i.IncreaseX(); i.IncreaseY() }
func (i *ID) MoveDownDiagonally()   { i.DecreaseX(); i.DecreaseY() }
func (i *ID) IsEqualTo(id *ID) bool { return id.x == i.x && id.y == i.y }

func (b *blockmap) IsValidID(id *ID) bool {
	if id.y > len(b.blocks) || id.x > len(b.blocks[0]) {
		return false
	}
	return true
}
func (b *blockmap) GetBlock(id *ID) *block {
	if !b.IsValidID(id) {
		return nil
	}
	return &b.blocks[id.y][id.x]
}
