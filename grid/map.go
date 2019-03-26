package grid

import "github.com/chfanghr/WTFCarProject/grid/bfs"

type block struct {
	x, y      int
	idx       int
	isBarrier bool
}

type BlockMap struct {
	blocks    [][]block
	idxBlocks map[int]*block
}

func NewBlockMap(sizeX, sizeY int) *BlockMap {
	if sizeX <= 0 || sizeY <= 0 {
		return nil
	}
	blm := &BlockMap{}
	blm.idxBlocks = make(map[int]*block)

	tmpBlks := make([]block, sizeX*sizeY)
	for i := 0; i < len(tmpBlks); i++ {
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

func (b *BlockMap) toGraph() *bfs.Graph {
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

func (b *BlockMap) getRelatedBlocks(blk *block) (res []*block) {
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

func (b *BlockMap) getBlock(x, y int) *block {
	if !b.isBlockExist(x, y) {
		return nil
	}
	return &b.blocks[y][x]
}

func (b *BlockMap) isBlockExist(x, y int) bool {
	return 0 < x && x < len(b.blocks[0]) && 0 < y && y < len(b.blocks)
}

func (b *BlockMap) Size() (x, y int) {
	return len(b.blocks[0]), len(b.blocks)
}

func (b *BlockMap) SetBarrier(x, y int, s bool) {
	if b.isBlockExist(x, y) {
		b.getBlock(x, y).isBarrier = s
	}
}

func (b *BlockMap) GetShortestPath(fromX, fromY, toX, toY int) (res []struct{ x, y int }) {
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
		res = append(res, struct{ x, y int }{x: b.idxBlocks[v].x, y: b.idxBlocks[v].y})
	}
	return
}
