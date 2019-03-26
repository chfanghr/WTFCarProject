package blockmap

import "github.com/chfanghr/WTFCarProject/blockmap/bfs"

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
		bs := b.getRelatedBlocks(v)
		for _, v := range bs {
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
	return x > 0 && x < len(b.blocks[0]) && 0 < y && y < len(b.blocks)
}

func (b *BlockMap) Size() (x, y int) {
	return len(b.blocks[0]), len(b.blocks)
}
