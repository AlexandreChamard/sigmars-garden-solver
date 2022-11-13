package sigmarsolver

import (
	"fmt"
)

func NewBoard(tiles [][]TileType) Board {
	board := Board{
		TileTypesRemainingMap: make(map[TileType]int),
	}

	if len(tiles) != nbLines {
		panic(fmt.Sprintf("tiles: wrong number of lines (expected %d got %d)", nbLines, len(tiles)))
	}
	for x, line := range tiles {
		if len(line) != lineSize[x] {
			panic(fmt.Sprintf("line %d: wrong line lentgh (expected %d got %d)", x, lineSize[x], len(line)))
		}
		for y, tile := range line {
			board.Board[FromXYPos(x, y)] = Tile{
				Type: tile,
			}
			if tile != TileType_EMPTY {
				board.TileTypesRemainingMap[tile]++
			}
		}
	}
	for x, line := range tiles {
		for y := range line {
			board.CheckLockState(x, y)
			board.setAlchemyDistance(x, y)
		}
	}
	return board
}

func (this *Board) CheckLockState(x, y int) bool {
	// Check empty
	if this.Board[FromXYPos(x, y)].Type == TileType_EMPTY {
		this.Board[FromXYPos(x, y)].Lock = false
		return false
	}

	// Check is alchemy locked
	isAlchemyLocked := func() bool {
		tileType := this.Board[FromXYPos(x, y)].Type
		tileAlchemyStage := tileType.GetAlchemyStage()
		return tileAlchemyStage != AlchemyStage_0 && tileAlchemyStage > this.AlchemyStage+1
	}
	if isAlchemyLocked() {
		this.Board[FromXYPos(x, y)].Lock = true
		return true
	}

	// Check is locked by
	joinedTilesPos := getAllPossibleJoinedTiles(x, y)

	var bits uint8
	count := 0

	for i, pos := range joinedTilesPos {
		if IsPossitionValid(pos.X, pos.Y) && this.Board[FromXYPos(pos.X, pos.Y)].Type != TileType_EMPTY {
			bits |= 1 << i
			count++
		}
	}
	if count <= 1 || count >= 4 {
		isLocked := count >= 4
		this.Board[FromXYPos(x, y)].Lock = isLocked
		return isLocked
	}

	// cpy the firse two bits to cycle the bits
	// because we check 3 bits at a time, we have to cycle at least 2 bits
	bits |= (bits & 0b11) << 6

	// the tile will be unlock if there is at least 3 empty tiles in a row
	isLocked := true
	for i := 0; i < 6; i++ {
		isLocked = isLocked && (bits>>i)&0b111 != 0
	}

	if !isLocked {
		Logf("%v (%d,%d):\t%d %d %d %d %d %d\n", this.Board[FromXYPos(x, y)], x, y, (bits>>0)&1, (bits>>1)&1, (bits>>2)&1, (bits>>3)&1, (bits>>4)&1, (bits>>5)&1)
	}

	this.Board[FromXYPos(x, y)].Lock = isLocked
	return isLocked
}

func (this *Board) setAlchemyDistance(x, y int) {
	alchemyStage := this.Board[FromXYPos(x, y)].Type.GetAlchemyStage()
	if alchemyStage != AlchemyStage_0 {
		possibleJoinedTiles := getAllPossibleJoinedTiles(x, y)
		this.setAlchemyDistanceTopRec(possibleJoinedTiles[0].X, possibleJoinedTiles[0].Y, int(alchemyStage)-1, 1)
		this.setAlchemyDistanceTopRec(possibleJoinedTiles[1].X, possibleJoinedTiles[1].Y, int(alchemyStage)-1, 1)
		this.setAlchemyDistanceBottomRec(possibleJoinedTiles[3].X, possibleJoinedTiles[3].Y, int(alchemyStage)-1, 1)
		this.setAlchemyDistanceBottomRec(possibleJoinedTiles[4].X, possibleJoinedTiles[4].Y, int(alchemyStage)-1, 1)
		this.setAlchemyDistanceLeftRightRec(possibleJoinedTiles[2].X, possibleJoinedTiles[2].Y, int(alchemyStage)-1, 1, +1)
		this.setAlchemyDistanceLeftRightRec(possibleJoinedTiles[5].X, possibleJoinedTiles[5].Y, int(alchemyStage)-1, 1, -1)
	}
}

func (this *Board) setAlchemyDistanceTopRec(x, y, alchNum, count int) {
	if !IsPossitionValid(x, y) {
		return
	}
	this.Board[FromXYPos(x, y)].DistanceToAlchs[alchNum] = count
	possibleJoinedTiles := getAllPossibleJoinedTiles(x, y)
	this.setAlchemyDistanceTopRec(possibleJoinedTiles[0].X, possibleJoinedTiles[0].Y, alchNum, count+1)
	this.setAlchemyDistanceTopRec(possibleJoinedTiles[1].X, possibleJoinedTiles[1].Y, alchNum, count+1)
	this.setAlchemyDistanceLeftRightRec(possibleJoinedTiles[2].X, possibleJoinedTiles[2].Y, alchNum, count+1, +1)
	this.setAlchemyDistanceLeftRightRec(possibleJoinedTiles[5].X, possibleJoinedTiles[5].Y, alchNum, count+1, -1)
}

func (this *Board) setAlchemyDistanceBottomRec(x, y, alchNum, count int) {
	if !IsPossitionValid(x, y) {
		return
	}
	this.Board[FromXYPos(x, y)].DistanceToAlchs[alchNum] = count
	possibleJoinedTiles := getAllPossibleJoinedTiles(x, y)
	this.setAlchemyDistanceBottomRec(possibleJoinedTiles[3].X, possibleJoinedTiles[3].Y, alchNum, count+1)
	this.setAlchemyDistanceBottomRec(possibleJoinedTiles[4].X, possibleJoinedTiles[4].Y, alchNum, count+1)
	this.setAlchemyDistanceLeftRightRec(possibleJoinedTiles[2].X, possibleJoinedTiles[2].Y, alchNum, count+1, +1)
	this.setAlchemyDistanceLeftRightRec(possibleJoinedTiles[5].X, possibleJoinedTiles[5].Y, alchNum, count+1, -1)
}

func (this *Board) setAlchemyDistanceLeftRightRec(x, y, alchNum, count, direction int) {
	if !IsPossitionValid(x, y) {
		return
	}
	this.Board[FromXYPos(x, y)].DistanceToAlchs[alchNum] = count
	this.setAlchemyDistanceLeftRightRec(x, y+direction, alchNum, count+1, direction)
}

func IsPossitionValid(x, y int) bool {
	return x >= 0 && x < nbLines && y >= 0 && y < lineSize[x]
}

func getAllPossibleJoinedTiles(x, y int) [6]Position {
	isUpper := 0
	if x >= nbLines/2 {
		isUpper = 1
	}

	return [6]Position{
		{x - 1, y - 1 + isUpper},
		{x - 1, y + isUpper},
		{x, y + 1},
		{x + 1, y + 1 - isUpper},
		{x + 1, y - isUpper},
		{x, y - 1},
	}
}

func getAllJoinedTiles(x, y int) []Position {
	possible := getAllPossibleJoinedTiles(x, y)
	joined := make([]Position, 0, 6)
	for _, pos := range possible {
		if IsPossitionValid(pos.X, pos.Y) {
			joined = append(joined, pos)
		}
	}
	return joined
}

func (this Board) getAllLockedJoinedTiles(x, y int) []Position {
	joined := getAllJoinedTiles(x, y)
	locked := make([]Position, 0, 6)
	for _, pos := range joined {
		if this.Board[FromXYPos(pos.X, pos.Y)].Lock {
			locked = append(locked, pos)
		}
	}
	return locked
}

func (this *Board) CheckLockFromTileRemoving(tileType TileType, x, y int) []Position {
	if this.Board[FromXYPos(x, y)].Lock {
		panic("trying to remove locked tile")
	}

	alchemyStage := tileType.GetAlchemyStage()
	if alchemyStage != AlchemyStage_0 && alchemyStage != this.AlchemyStage {
		panic("trying to remove bad alchemy tile")
	}

	var unlockedPosition []Position

	if alchemyStage != AlchemyStage_0 {
		nextAlchemy := this.AlchemyStage.GetNextAlchemyType()
		for i, tile := range this.Board {
			if tile.Type == nextAlchemy {
				if tile.Lock {
					x, y := ToXYPos(i)
					if isLocked := this.CheckLockState(x, y); !isLocked {
						unlockedPosition = append(unlockedPosition, Position{x, y})
					}
				}
				break
			}
		}
	}

	for _, pos := range this.getAllLockedJoinedTiles(x, y) {
		if isLocked := this.CheckLockState(pos.X, pos.Y); !isLocked {
			unlockedPosition = append(unlockedPosition, Position{pos.X, pos.Y})
		}
	}

	return unlockedPosition
}
