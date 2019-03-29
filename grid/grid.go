package grid

import "github.com/chfanghr/WTFCarProject/grid/bfs"

type block struct {
	x, y      int
	idx       int
	isBarrier bool
}

type Grid struct {
	blocks    [][]block
	idxBlocks map[int]*block
}

func NewGrid(sizeX, sizeY int) *Grid {
	if sizeX <= 0 || sizeY <= 0 {
		return nil
	}
	blm := &Grid{}
	blm.idxBlocks = make(map[int]*block)

	tmpBlks := make([]block, sizeX*sizeY)
	for i := range tmpBlks {
		blm.idxBlocks[i] = &tmpBlks[i]
		tmpBlks[i].idx = i
		tmpBlks[i].isBarrier = false
	}

	blm.blocks = make([][]block, sizeY)
	for i := 0; i < sizeY; i++ {
		blm.blocks[i] = tmpBlks[:sizeX]
		tmpBlks = tmpBlks[sizeX:]
	}

	return blm
}

func (b *Grid) toGraph() *bfs.Graph {
	g := bfs.NewGraph()
	for k, v := range b.idxBlocks {
		if v.isBarrier {
			continue
		}
		bs := b.getRelatedBlocks(v)
		for _, v := range bs {
			if v.isBarrier {
				continue
			}
			g.AddEdge(k, v.idx)
		}
	}
	return g
}

func (b *Grid) getRelatedBlocks(blk *block) (res []*block) {
	left := b.getBlock(blk.x-1, blk.y)
	right := b.getBlock(blk.x+1, blk.y)
	below := b.getBlock(blk.x, blk.y-1)
	above := b.getBlock(blk.x, blk.y+1)

	//* obliquely *
	loa := b.getBlock(blk.x-1, blk.y+1)
	lob := b.getBlock(blk.x-1, blk.y-1)
	roa := b.getBlock(blk.x+1, blk.y+1)
	rob := b.getBlock(blk.x+1, blk.y-1)

	tmp := []*block{left, right, below, above, loa, lob, roa, rob}
	for _, v := range tmp {
		if v != nil {
			res = append(res, v)
		}
	}
	return
}

func (b *Grid) getBlock(x, y int) *block {
	if !b.isBlockExist(x, y) {
		return nil
	}
	return &b.blocks[y][x]
}

func (b *Grid) isBlockExist(x, y int) bool {
	return 0 < x && x < len(b.blocks[0]) && 0 < y && y < len(b.blocks)
}

func (b *Grid) Size() (x, y int) {
	return len(b.blocks[0]), len(b.blocks)
}

func (b *Grid) SetBarrier(x, y int, s bool) {
	if b.isBlockExist(x, y) {
		b.getBlock(x, y).isBarrier = s
	}
}

func (b *Grid) GetShortestPath(fromX, fromY, toX, toY int) (res []struct{ X, Y int }) {
	if !(b.isBlockExist(fromX, fromY) && b.isBlockExist(toX, toY)) {
		return nil
	}
	g := b.toGraph()
	resChan := bfs.NewBFSPath(g, b.getBlock(fromX, fromY).idx).PathTo(b.getBlock(toX, toY).idx)
	var resIdxes []int
	for v := range resChan {
		resIdxes = append(resIdxes, v.(int))
	}
	for _, v := range resIdxes {
		res = append(res, struct{ X, Y int }{X: b.idxBlocks[v].x, Y: b.idxBlocks[v].y})
	}
	return
}
