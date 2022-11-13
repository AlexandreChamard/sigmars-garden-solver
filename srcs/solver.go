package sigmarsolver

import (
	"fmt"
	"sort"
	"strings"
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

func (this *Board) Solve() ([]Action, int64, time.Duration) {
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

	var n int64 // number of iteration
	var k int   // used to show advancement
	start := time.Now()
	iterators := []*iterator{newIterator(this, possibilities)}

	for len(iterators) > 0 && !possibilities.IsEmpty() {
		if k == 100000 {
			fmt.Println(n, time.Since(start), iteratorsToString(iterators))
			k = 0
		}
		n++
		k++

		pos1, pos2, found := iterators[len(iterators)-1].Next()
		if found {
			doAction(pos1.X, pos1.Y, pos2.X, pos2.Y)
			iterators = append(iterators, newIterator(this, possibilities))
		} else {
			undoLastAction()
			iterators = iterators[:len(iterators)-1]
		}
	}

	return actions, n, time.Since(start)
}

func SolutionToString(actions []Action) string {
	var s string
	for _, action := range actions {
		s += fmt.Sprintf("%6s %6s {x:%2d y:%2d} {x:%2d y:%2d}\n", action.Type1, action.Type2, action.X1, action.Y1, action.X2, action.Y2)
	}
	return s
}

type iterator struct {
	n                  int
	foundPossibilities []possibleSolution
}

type possibleSolution struct {
	p1, p2 Position
}

func newIterator(board *Board, possibilities Possibilities) *iterator {
	var foundPossibilities []possibleSolution

	// KEY + (L1 | L2 | L3 | L4 | L5)
	for _, pos1 := range possibilities[TileType_KEY] {
		for _, pos2 := range possibilities[TileType_L1] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
		for _, pos2 := range possibilities[TileType_L2] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
		for _, pos2 := range possibilities[TileType_L3] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
		for _, pos2 := range possibilities[TileType_L4] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
		for _, pos2 := range possibilities[TileType_L5] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
	}

	// L6
	if len(possibilities[TileType_L6]) == 1 {
		foundPossibilities = append(foundPossibilities, possibleSolution{possibilities[TileType_L6][0], possibilities[TileType_L6][0]})
	}

	// LIGHT + DARK
	for _, pos1 := range possibilities[TileType_LIGHT] {
		for _, pos2 := range possibilities[TileType_DARK] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
	}

	// CYAN + CYAN
	for i, pos1 := range possibilities[TileType_CYAN] {
		for _, pos2 := range possibilities[TileType_CYAN][i+1:] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
	}

	// ORANGE + ORANGE
	for i, pos1 := range possibilities[TileType_ORANGE] {
		for _, pos2 := range possibilities[TileType_ORANGE][i+1:] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
	}

	// BLUE + BLUE
	for i, pos1 := range possibilities[TileType_BLUE] {
		for _, pos2 := range possibilities[TileType_BLUE][i+1:] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
	}

	// GREEN + GREEN
	for i, pos1 := range possibilities[TileType_GREEN] {
		for _, pos2 := range possibilities[TileType_GREEN][i+1:] {
			foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
		}
	}

	// WHITE + ( CYAN | ORANGE | BLUE | GREEN | WHITE )
	for i, pos1 := range possibilities[TileType_WHITE] {
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] || board.TileTypesRemainingMap[TileType_CYAN]%2 == 1 {
			for _, pos2 := range possibilities[TileType_CYAN] {
				foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
			}
		}
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] || board.TileTypesRemainingMap[TileType_ORANGE]%2 == 1 {
			for _, pos2 := range possibilities[TileType_ORANGE] {
				foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
			}
		}
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] || board.TileTypesRemainingMap[TileType_BLUE]%2 == 1 {
			for _, pos2 := range possibilities[TileType_BLUE] {
				foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
			}
		}
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] || board.TileTypesRemainingMap[TileType_GREEN]%2 == 1 {
			for _, pos2 := range possibilities[TileType_GREEN] {
				foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
			}
		}
		if board.WhiteUsedWithColored < board.TileTypesRemainingMap[TileType_WHITE] {
			for _, pos2 := range possibilities[TileType_WHITE][i+1:] {
				foundPossibilities = append(foundPossibilities, possibleSolution{pos1, pos2})
			}
		}
	}

	currentStage := int(board.AlchemyStage)

	// sort possibilities
	if currentStage != AlchemyStage_FINAL {
		sort.Slice(foundPossibilities, func(i, j int) bool {
			p11Dist := board.Board[FromXYPos(foundPossibilities[i].p1.X, foundPossibilities[i].p1.Y)].DistanceToAlchs[currentStage]
			p21Dist := board.Board[FromXYPos(foundPossibilities[i].p2.X, foundPossibilities[i].p2.Y)].DistanceToAlchs[currentStage]
			p12Dist := board.Board[FromXYPos(foundPossibilities[j].p1.X, foundPossibilities[j].p1.Y)].DistanceToAlchs[currentStage]
			p22Dist := board.Board[FromXYPos(foundPossibilities[j].p2.X, foundPossibilities[j].p2.Y)].DistanceToAlchs[currentStage]

			if p11Dist < p12Dist && p11Dist < p22Dist || (p21Dist < p12Dist && p21Dist < p22Dist) {
				return true // i < j
			} else if p11Dist > p12Dist && p11Dist > p22Dist || (p21Dist > p12Dist && p21Dist > p22Dist) {
				return false // i > j
			} else {
				return p11Dist+p21Dist < p12Dist+p22Dist
			}
		})
	}

	return &iterator{
		n:                  0,
		foundPossibilities: foundPossibilities,
	}
}

func (this *iterator) Next() (Position, Position, bool) {
	if this == nil || this.n >= len(this.foundPossibilities) {
		return Position{-1, -1}, Position{-1, -1}, false
	}
	p1, p2 := this.foundPossibilities[this.n].p1, this.foundPossibilities[this.n].p2
	this.n++
	return p1, p2, true
}

func iteratorsToString(its []*iterator) string {
	var a []string
	for _, it := range its {
		a = append(a, fmt.Sprintf("%d/%d", it.n, len(it.foundPossibilities)))
	}
	return fmt.Sprintf("[ %s ]", strings.Join(a, ", "))
}
