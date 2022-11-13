package sigmarsolver

import "fmt"

type Board struct {
	Board                 [91]Tile
	TileTypesRemainingMap map[TileType]int
	AlchemyStage          AlchemyStage
	WhiteUsedWithColored  int
}

type Tile struct {
	Type            TileType
	Lock            bool
	DistanceToAlchs [AlchemyStage_FINAL]int
}

type TileType string
type AlchemyStage int

const (
	TileType_EMPTY   TileType = ""
	TileType_WHITE   TileType = "white"
	TileType_CYAN    TileType = "cyan"
	TileType_ORANGE  TileType = "orange"
	TileType_BLUE    TileType = "blue"
	TileType_GREEN   TileType = "green"
	TileType_LIGHT   TileType = "light"
	TileType_DARK    TileType = "dark"
	TileType_KEY     TileType = "key"
	TileType_L1      TileType = "l1"
	TileType_L2      TileType = "l2"
	TileType_L3      TileType = "l3"
	TileType_L4      TileType = "l4"
	TileType_L5      TileType = "l5"
	TileType_L6      TileType = "l6"
	TileType_L_FINAL TileType = "no-next-alchemy"
)

const (
	AlchemyStage_0 = iota
	AlchemyStage_1
	AlchemyStage_2
	AlchemyStage_3
	AlchemyStage_4
	AlchemyStage_5
	AlchemyStage_FINAL
)

func (this TileType) Valid() {
	switch this {
	case TileType_EMPTY,
		TileType_WHITE,
		TileType_CYAN,
		TileType_ORANGE,
		TileType_BLUE,
		TileType_GREEN,
		TileType_LIGHT,
		TileType_DARK,
		TileType_KEY,
		TileType_L1,
		TileType_L2,
		TileType_L3,
		TileType_L4,
		TileType_L5,
		TileType_L6:
		return
	default:
		panic(fmt.Sprintf("invalid %s", this))
	}
}

func (this TileType) GetAlchemyStage() AlchemyStage {
	switch this {
	case TileType_L1:
		return AlchemyStage_1
	case TileType_L2:
		return AlchemyStage_2
	case TileType_L3:
		return AlchemyStage_3
	case TileType_L4:
		return AlchemyStage_4
	case TileType_L5:
		return AlchemyStage_5
	case TileType_L6, TileType_L_FINAL:
		return AlchemyStage_FINAL
	default:
		return AlchemyStage_0
	}
}

func (this AlchemyStage) GetNextAlchemyType() TileType {
	switch this {
	case AlchemyStage_0:
		return TileType_L1
	case AlchemyStage_1:
		return TileType_L2
	case AlchemyStage_2:
		return TileType_L3
	case AlchemyStage_3:
		return TileType_L4
	case AlchemyStage_4:
		return TileType_L5
	case AlchemyStage_5:
		return TileType_L6
	case AlchemyStage_FINAL:
		return TileType_L_FINAL
	default:
		panic("invalid alchemy stage")
	}
}

const nbLines = 11

var lineSize = [nbLines]int{6, 7, 8, 9, 10, 11, 10, 9, 8, 7, 6}
var startLineValue = [nbLines + 1]int{0, 6, 13, 21, 30, 40, 51, 61, 70, 78, 85, 91}

type Position struct {
	X, Y int
}

func ToXYPos(pos int) (int, int) {
	x := nbLines
	for pos-startLineValue[x] < 0 {
		x--
	}
	return x, pos - startLineValue[x]
}

func FromXYPos(x, y int) int {
	return startLineValue[x] + y
}
