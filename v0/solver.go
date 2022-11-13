package v0

import (
	"fmt"
	"sort"
	"time"
)

type Action struct {
	X1, Y1, X2, Y2 int
	Type1, Type2   TileType
	Unlocked       []Position
}

type Possibilities map[TileType][]Position

func (this Possibilities) Insert(tileType TileType, pos Position) {
	this[tileType] = append(this[tileType], pos)
}

func (this Possibilities) Remove(pos Position) {
	for k, v := range this {
		for i, p := range v {
			if p == pos {
				this[k][i] = this[k][len(v)-1]
				this[k] = this[k][:len(v)-1]
				return
			}
		}
	}
	panic(fmt.Sprintf("pos %v not found in possibilities", pos))
}

func (this Possibilities) Sort() {
	for k, v := range this {
		sort.Slice(this[k], func(i, j int) bool { return v[i].X < v[j].X || (v[i].X == v[j].X && v[i].Y < v[j].Y) })
	}
}

func (this Possibilities) IsEmpty() bool {
	for _, arr := range this {
		if len(arr) > 0 {
			return false
		}
	}
	return true
}

func (this *Board) Solve() ([]Action, int64) {
	start := time.Now()

	actions := make([]Action, 0, 91)
	possibilities := Possibilities{}

	doAction := func(x1, y1, x2, y2 int) {
		action := Action{
			X1: x1, Y1: y1, Type1: this.Board[FromXYPos(x1, y1)].Type,
			X2: x2, Y2: y2, Type2: this.Board[FromXYPos(x2, y2)].Type,
		}

		this.TileTypesRemainingMap[action.Type1]--
		this.TileTypesRemainingMap[action.Type2]--

		if action.Type1 == TileType_WHITE && action.Type2 != TileType_WHITE {
			if this.TileTypesRemainingMap[action.Type2]%2 == 1 {
				this.WhiteUsedWithColored++
			} else {
				this.WhiteUsedWithColored--
			}
		}

		this.Board[FromXYPos(x1, y1)].Type = TileType_EMPTY
		possibilities.Remove(Position{x1, y1})

		if x1 != x2 || y1 != y2 {
			this.Board[FromXYPos(x2, y2)].Type = TileType_EMPTY
			possibilities.Remove(Position{x2, y2})
		}

		if action.Type1.GetAlchemyStage() != AlchemyStage_0 || action.Type2.GetAlchemyStage() != AlchemyStage_0 {
			this.AlchemyStage++
		}

		action.Unlocked = append(this.CheckLockFromTileRemoving(action.Type1, x1, y1), this.CheckLockFromTileRemoving(action.Type2, x2, y2)...)

		for _, pos := range action.Unlocked {
			possibilities.Insert(this.Board[FromXYPos(pos.X, pos.Y)].Type, pos)
		}

		actions = append(actions, action)
		possibilities.Sort()
	}

	undoLastAction := func() {
		action := actions[len(actions)-1]

		this.TileTypesRemainingMap[action.Type1]++
		this.TileTypesRemainingMap[action.Type2]++

		if action.Type1 == TileType_WHITE && action.Type2 != TileType_WHITE {
			if this.TileTypesRemainingMap[action.Type2]%2 == 1 {
				this.WhiteUsedWithColored++
			} else {
				this.WhiteUsedWithColored--
			}
		}

		this.Board[FromXYPos(action.X1, action.Y1)].Type = action.Type1
		possibilities.Insert(action.Type1, Position{action.X1, action.Y1})

		if action.X1 != action.X2 || action.Y1 != action.Y2 {
			this.Board[FromXYPos(action.X2, action.Y2)].Type = action.Type2
			possibilities.Insert(action.Type2, Position{action.X2, action.Y2})
		}

		if action.Type1.GetAlchemyStage() != AlchemyStage_0 || action.Type2.GetAlchemyStage() != AlchemyStage_0 {
			this.AlchemyStage--
		}

		for _, pos := range action.Unlocked {
			this.Board[FromXYPos(pos.X, pos.Y)].Lock = true

			possibilities.Remove(pos)
		}

		actions = actions[:len(actions)-1]
		possibilities.Sort()
	}

	// Fill the possibilities with all unlocked tiles
	for i, tile := range this.Board {
		if tile.Type != TileType_EMPTY && !tile.Lock {
			x, y := ToXYPos(i)
			possibilities.Insert(tile.Type, Position{x, y})
		}
	}
	possibilities.Sort()

	n := int64(0)
	k := 0
	iterators := []int{0}
	for len(iterators) > 0 && !possibilities.IsEmpty() {
		if k == 100000 {
			fmt.Println(n, time.Since(start), iterators)
			k = 0
		}
		n++
		k++

		lastItPos := len(iterators) - 1
		pos1, pos2, nextIt := possibilities.GetNextPossibility(this, iterators[lastItPos])
		if nextIt != -1 {
			iterators[lastItPos] = nextIt
			doAction(pos1.X, pos1.Y, pos2.X, pos2.Y)
			iterators = append(iterators, 0)
		} else {
			undoLastAction()
			iterators = iterators[:lastItPos]
		}
	}
	return actions, n
}

func SolutionToString(actions []Action) string {
	s := ""
	for _, action := range actions {
		s += fmt.Sprintf("%6s %6s {x:%2d y:%2d} {x:%2d y:%2d}\n", action.Type1, action.Type2, action.X1, action.Y1, action.X2, action.Y2)
	}
	return s
}

func (this Possibilities) GetNextPossibility(board *Board, iterator int) (Position, Position, int) {
	if iterator < 0 {
		panic("invalid iterator")
	}

	count := 0

	// KEY + (L1 | L2 | L3 | L4 | L5)
	for _, pos1 := range this[TileType_KEY] {
		for _, pos2 := range this[TileType_L1] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
		for _, pos2 := range this[TileType_L2] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
		for _, pos2 := range this[TileType_L3] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
		for _, pos2 := range this[TileType_L4] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
		for _, pos2 := range this[TileType_L5] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
	}

	// L6
	if len(this[TileType_L6]) == 1 {
		if count == iterator {
			return this[TileType_L6][0], this[TileType_L6][0], count + 1
		}
		count++
	}

	// LIGHT + DARK
	for _, pos1 := range this[TileType_LIGHT] {
		for _, pos2 := range this[TileType_DARK] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
	}

	// CYAN + CYAN
	for i, pos1 := range this[TileType_CYAN] {
		for _, pos2 := range this[TileType_CYAN][i+1:] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
	}

	// ORANGE + ORANGE
	for i, pos1 := range this[TileType_ORANGE] {
		for _, pos2 := range this[TileType_ORANGE][i+1:] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
	}

	// BLUE + BLUE
	for i, pos1 := range this[TileType_BLUE] {
		for _, pos2 := range this[TileType_BLUE][i+1:] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
	}

	// GREEN + GREEN
	for i, pos1 := range this[TileType_GREEN] {
		for _, pos2 := range this[TileType_GREEN][i+1:] {
			if count == iterator {
				return pos1, pos2, count + 1
			}
			count++
		}
	}

	// WHITE + ( CYAN | ORANGE | BLUE | GREEN | WHITE )
	for i, pos1 := range this[TileType_WHITE] {
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] || board.TileTypesRemainingMap[TileType_CYAN]%2 == 1 {
			for _, pos2 := range this[TileType_CYAN] {
				if count == iterator {
					return pos1, pos2, count + 1
				}
				count++
			}
		}
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] || board.TileTypesRemainingMap[TileType_ORANGE]%2 == 1 {
			for _, pos2 := range this[TileType_ORANGE] {
				if count == iterator {
					return pos1, pos2, count + 1
				}
				count++
			}
		}
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] || board.TileTypesRemainingMap[TileType_BLUE]%2 == 1 {
			for _, pos2 := range this[TileType_BLUE] {
				if count == iterator {
					return pos1, pos2, count + 1
				}
				count++
			}
		}
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] || board.TileTypesRemainingMap[TileType_GREEN]%2 == 1 {
			for _, pos2 := range this[TileType_GREEN] {
				if count == iterator {
					return pos1, pos2, count + 1
				}
				count++
			}
		}
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] {
			for _, pos2 := range this[TileType_WHITE][i+1:] {
				if count == iterator {
					return pos1, pos2, count + 1
				}
				count++
			}
		}
	}

	return Position{}, Position{}, -1
}
