package cube

import (
	"fmt"
	"strings"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
)

type Cube struct {
	State [][][]int
	Scramble string
	Solution string
}

var VALID_MOVES = []string{"U", "D", "R", "L", "F", "B", "U'", "D'", "R'", "L'", "F'", "B'", "U2", "D2", "R2", "L2", "F2", "B2"}

func InitialState() [][][]int {
	state := make([][][]int, 0)

	SIDES := 6
	for i := 0; i < SIDES; i++ {
		side := make([][]int, 0)

		line := make([]int, 0)
		for j := 0; j < 3; j++ {
			line = append(line, i)
		}

		for j := 0; j < 3; j++ {
			side = append(side, line)
		}

		state = append(state, side)
	}

	return state
}

func IndexFunc[T any](arr []T, test func(T) bool) int {
	for idx, el := range arr {
		if test(el) {
			return idx
		}
	}

	return -1
}

func (c *Cube) ValidMoves() bool {
	moves := strings.Split(c.Solution, " ")
	for _, move := range moves {
		idx := IndexFunc(VALID_MOVES, func (m string) bool { return m == move })
		if idx == -1 { return false }
	}

	return true
}

func ParseFMCSolutionToMilliseconds(scramble string, solution string) int {
	c := Cube{InitialState(), scramble, solution}

	fmt.Println("scramble: ", scramble)
	fmt.Println("solution: ", solution)

	if !c.ValidMoves() { return constants.DNF }
	
	return constants.DNF
}