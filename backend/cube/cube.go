package cube

import (
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

func (c *Cube) RotateFace(face int) {
	for i, j := 0, len(c.State[face])-1; i<j; i, j = i+1, j-1 {
		c.State[face][i], c.State[face][j] = c.State[face][j], c.State[face][i]
	}

	for i := 0; i < len(c.State[face]); i++ {
		for j := 0; j < i; j++ {
			c.State[face][i][j], c.State[face][j][i] = c.State[face][j][i], c.State[face][i][j]
		}
	}
}

func (c *Cube) ApplyU() {
	c.RotateFace(0)
	c.State[1][0][0], c.State[1][0][1], c.State[1][0][2], c.State[2][0][0], c.State[2][0][1], c.State[2][0][2], c.State[3][0][0], c.State[3][0][1], c.State[3][0][2], c.State[4][0][0], c.State[4][0][1], c.State[4][0][2] = c.State[2][0][0], c.State[2][0][1], c.State[2][0][2], c.State[3][0][0], c.State[3][0][1], c.State[3][0][2], c.State[4][0][0], c.State[4][0][1], c.State[4][0][2], c.State[1][0][0], c.State[1][0][1], c.State[1][0][2]
}

func (c *Cube) ApplyUPrime() {
	c.RotateFace(0)
	c.RotateFace(0)
	c.RotateFace(0)
	c.State[1][0][0], c.State[1][0][1], c.State[1][0][2], c.State[2][0][0], c.State[2][0][1], c.State[2][0][2], c.State[3][0][0], c.State[3][0][1], c.State[3][0][2], c.State[4][0][0], c.State[4][0][1], c.State[4][0][2] = c.State[4][0][0], c.State[4][0][1], c.State[4][0][2], c.State[1][0][0], c.State[1][0][1], c.State[1][0][2], c.State[2][0][0], c.State[2][0][1], c.State[2][0][2], c.State[3][0][0], c.State[3][0][1], c.State[3][0][2]
}

func (c *Cube) ApplyD() {
	c.RotateFace(5)
	c.State[1][2][0], c.State[1][2][1], c.State[1][2][2], c.State[2][2][0], c.State[2][2][1], c.State[2][2][2], c.State[3][2][0], c.State[3][2][1], c.State[3][2][2], c.State[4][2][0], c.State[4][2][1], c.State[4][2][2] = c.State[2][2][0], c.State[2][2][1], c.State[2][2][2], c.State[3][2][0], c.State[3][2][1], c.State[3][2][2], c.State[4][2][0], c.State[4][2][1], c.State[4][2][2], c.State[1][2][0], c.State[1][2][1], c.State[1][2][2]
}

func (c *Cube) ApplyDPrime() {
	c.RotateFace(5)
	c.RotateFace(5)
	c.RotateFace(5)
	c.State[1][2][0], c.State[1][2][1], c.State[1][2][2], c.State[2][2][0], c.State[2][2][1], c.State[2][2][2], c.State[3][2][0], c.State[3][2][1], c.State[3][2][2], c.State[4][2][0], c.State[4][2][1], c.State[4][2][2] = c.State[4][2][0], c.State[4][2][1], c.State[4][2][2], c.State[1][2][0], c.State[1][2][1], c.State[1][2][2], c.State[2][2][0], c.State[2][2][1], c.State[2][2][2], c.State[3][2][0], c.State[3][2][1], c.State[3][2][2]
}

func (c *Cube) ApplyX() {
	c.RotateFace(1)
	c.RotateFace(1)
	c.RotateFace(1)

	c.RotateFace(3)

	c.RotateFace(4)
	c.RotateFace(4)

	c.State[0], c.State[2], c.State[5], c.State[4] = c.State[2], c.State[5], c.State[4], c.State[0]
}

func (c *Cube) ApplyXPrime() {
	c.RotateFace(1)

	c.RotateFace(3)
	c.RotateFace(3)
	c.RotateFace(3)

	c.RotateFace(4)
	c.RotateFace(4)

	c.State[0], c.State[2], c.State[5], c.State[4] = c.State[4], c.State[0], c.State[2], c.State[5]
}

func (c *Cube) ApplyY() {
	c.RotateFace(0)

	c.RotateFace(5)
	c.RotateFace(5)
	c.RotateFace(5)

	c.State[1], c.State[2], c.State[3], c.State[4] = c.State[2], c.State[3], c.State[4], c.State[1]
}

func (c *Cube) ApplyYPrime() {
	c.RotateFace(0)
	c.RotateFace(0)
	c.RotateFace(0)
	
	c.RotateFace(5)

	c.State[1], c.State[2], c.State[3], c.State[4] = c.State[4], c.State[1], c.State[2], c.State[3]
}

func (c *Cube) ApplyZ() {
	c.RotateFace(2)

	c.RotateFace(4)
	c.RotateFace(4)
	c.RotateFace(4)

	for _, f := range []int{0, 1, 5, 3} {
		for range []int{0, 1, 2} {
			c.RotateFace(f)
		}
	}

	c.State[0], c.State[1], c.State[5], c.State[3] = c.State[1], c.State[5], c.State[3], c.State[0]
}

func (c *Cube) ApplyZPrime() {
	c.RotateFace(2)
	c.RotateFace(2)
	c.RotateFace(2)
	
	c.RotateFace(4)

	for _, f := range []int{0, 1, 5, 3} {
		c.RotateFace(f)
	}

	c.State[0], c.State[1], c.State[5], c.State[3] = c.State[3], c.State[0], c.State[1], c.State[5]
}

func (c *Cube) ApplyMove(move string) {
	switch move {
		case "U":
			c.ApplyU()
		case "U'":
			c.ApplyUPrime()
		case "U2":
			c.ApplyU()
			c.ApplyU()
		case "D":
			c.ApplyD()
		case "D'":
			c.ApplyDPrime()
		case "D2":
			c.ApplyD()
			c.ApplyD()
		case "x":
			c.ApplyX()
		case "x'":
			c.ApplyXPrime()
		case "x2":
			c.ApplyX()
			c.ApplyX()
		case "y":
			c.ApplyY()
		case "y'":
			c.ApplyYPrime()
		case "y2":
			c.ApplyY()
			c.ApplyY()
		case "z":
			c.ApplyZ()
		case "z'":
			c.ApplyZPrime()
		case "z2":
			c.ApplyZ()
			c.ApplyZ()
	}
}

func (c *Cube) ApplyAlgorihtm(algorihtm string) {
	moves := strings.Split(algorihtm, " ")
	for _, move := range moves {
		c.ApplyMove(move)
	}
}

func (c *Cube) ApplyScramble() {
	c.ApplyAlgorihtm(c.Scramble)
}

func (c *Cube) ApplySolution() {
	c.ApplyAlgorihtm(c.Solution)
}

func (c *Cube) Solved() bool {
	for face := range c.State {
		all := true
		for i := range c.State[face] {
			for j := range c.State[face][i] {
				if c.State[face][i][j] != c.State[face][0][0] {
					all = false
				}
			}
		}
		if !all {
			return false
		}
	}
	
	return true
}

func (c *Cube) SolutionLength() int {
	return len(strings.Split(c.Solution, " "))
}

func ParseFMCSolutionToMilliseconds(scramble string, solution string) int {
	c := Cube{InitialState(), scramble, solution}

	if !c.ValidMoves() { return constants.DNF }

	c.ApplyScramble()
	c.ApplySolution()
	if !c.Solved() { return constants.DNF }
	
	return c.SolutionLength()
}