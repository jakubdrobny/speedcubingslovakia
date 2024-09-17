package cube

import (
	"fmt"
	"strings"

	"github.com/jakubdrobny/speedcubingslovakia/backend/constants"
)

type Cube struct {
	State    [6][3][3]int
	Scramble string
	Solution string
}

var VALID_MOVES = []string{"U", "D", "R", "L", "F", "B", "U'", "D'", "R'", "L'", "F'", "B'", "U2", "D2", "R2", "L2", "F2", "B2", "Uw", "Dw", "Rw", "Lw", "Fw", "Bw", "Uw'", "Dw'", "Rw'", "Lw'", "Fw'", "Bw'", "Uw2", "Dw2", "Rw2", "Lw2", "Fw2", "Bw2", "x", "x'", "x2", "y", "y'", "y2", "z", "z'", "z2"}
var ROTATIONS = []string{"x", "x'", "x2", "y", "y'", "y2", "z", "z'", "z2"}
var COLORS = []string{"W", "O", "G", "R", "B", "Y"}

func InitialState() [6][3][3]int {
	state := [6][3][3]int{}

	SIDES := 6
	for i := 0; i < SIDES; i++ {
		side := [3][3]int{}

		line := [3]int{}
		for j := 0; j < 3; j++ {
			line[j] = i
		}

		for j := 0; j < 3; j++ {
			side[j] = line
		}

		state[i] = side
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
		idx := IndexFunc(VALID_MOVES, func(m string) bool { return m == move })
		if idx == -1 {
			return false
		}
	}

	return true
}

func (c *Cube) RotateFace(face int) {
	for i, j := 0, len(c.State[face])-1; i < j; i, j = i+1, j-1 {
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
	c.State[1][2][0], c.State[1][2][1], c.State[1][2][2], c.State[2][2][0], c.State[2][2][1], c.State[2][2][2], c.State[3][2][0], c.State[3][2][1], c.State[3][2][2], c.State[4][2][0], c.State[4][2][1], c.State[4][2][2] = c.State[4][2][0], c.State[4][2][1], c.State[4][2][2], c.State[1][2][0], c.State[1][2][1], c.State[1][2][2], c.State[2][2][0], c.State[2][2][1], c.State[2][2][2], c.State[3][2][0], c.State[3][2][1], c.State[3][2][2]
}

func (c *Cube) ApplyDPrime() {
	c.RotateFace(5)
	c.RotateFace(5)
	c.RotateFace(5)
	c.State[1][2][0], c.State[1][2][1], c.State[1][2][2], c.State[2][2][0], c.State[2][2][1], c.State[2][2][2], c.State[3][2][0], c.State[3][2][1], c.State[3][2][2], c.State[4][2][0], c.State[4][2][1], c.State[4][2][2] = c.State[2][2][0], c.State[2][2][1], c.State[2][2][2], c.State[3][2][0], c.State[3][2][1], c.State[3][2][2], c.State[4][2][0], c.State[4][2][1], c.State[4][2][2], c.State[1][2][0], c.State[1][2][1], c.State[1][2][2]
}

func (c *Cube) ApplyF() {
	c.RotateFace(2)
	c.State[0][2][0], c.State[0][2][1], c.State[0][2][2], c.State[3][0][0], c.State[3][1][0], c.State[3][2][0], c.State[5][0][2], c.State[5][0][1], c.State[5][0][0], c.State[1][2][2], c.State[1][1][2], c.State[1][0][2] = c.State[1][2][2], c.State[1][1][2], c.State[1][0][2], c.State[0][2][0], c.State[0][2][1], c.State[0][2][2], c.State[3][0][0], c.State[3][1][0], c.State[3][2][0], c.State[5][0][2], c.State[5][0][1], c.State[5][0][0]
}

func (c *Cube) ApplyFPrime() {
	c.RotateFace(2)
	c.RotateFace(2)
	c.RotateFace(2)
	c.State[0][2][0], c.State[0][2][1], c.State[0][2][2], c.State[3][0][0], c.State[3][1][0], c.State[3][2][0], c.State[5][0][2], c.State[5][0][1], c.State[5][0][0], c.State[1][2][2], c.State[1][1][2], c.State[1][0][2] = c.State[3][0][0], c.State[3][1][0], c.State[3][2][0], c.State[5][0][2], c.State[5][0][1], c.State[5][0][0], c.State[1][2][2], c.State[1][1][2], c.State[1][0][2], c.State[0][2][0], c.State[0][2][1], c.State[0][2][2]
}

func (c *Cube) ApplyB() {
	c.RotateFace(4)
	c.State[0][0][0], c.State[0][0][1], c.State[0][0][2], c.State[3][0][2], c.State[3][1][2], c.State[3][2][2], c.State[5][2][2], c.State[5][2][1], c.State[5][2][0], c.State[1][2][0], c.State[1][1][0], c.State[1][0][0] = c.State[3][0][2], c.State[3][1][2], c.State[3][2][2], c.State[5][2][2], c.State[5][2][1], c.State[5][2][0], c.State[1][2][0], c.State[1][1][0], c.State[1][0][0], c.State[0][0][0], c.State[0][0][1], c.State[0][0][2]
}

func (c *Cube) ApplyBPrime() {
	c.RotateFace(4)
	c.RotateFace(4)
	c.RotateFace(4)
	c.State[0][0][0], c.State[0][0][1], c.State[0][0][2], c.State[3][0][2], c.State[3][1][2], c.State[3][2][2], c.State[5][2][2], c.State[5][2][1], c.State[5][2][0], c.State[1][2][0], c.State[1][1][0], c.State[1][0][0] = c.State[1][2][0], c.State[1][1][0], c.State[1][0][0], c.State[0][0][0], c.State[0][0][1], c.State[0][0][2], c.State[3][0][2], c.State[3][1][2], c.State[3][2][2], c.State[5][2][2], c.State[5][2][1], c.State[5][2][0]
}

func (c *Cube) ApplyR() {
	c.RotateFace(3)
	c.State[0][0][2], c.State[0][1][2], c.State[0][2][2], c.State[2][0][2], c.State[2][1][2], c.State[2][2][2], c.State[5][0][2], c.State[5][1][2], c.State[5][2][2], c.State[4][2][0], c.State[4][1][0], c.State[4][0][0] = c.State[2][0][2], c.State[2][1][2], c.State[2][2][2], c.State[5][0][2], c.State[5][1][2], c.State[5][2][2], c.State[4][2][0], c.State[4][1][0], c.State[4][0][0], c.State[0][0][2], c.State[0][1][2], c.State[0][2][2]
}

func (c *Cube) ApplyRPrime() {
	c.RotateFace(3)
	c.RotateFace(3)
	c.RotateFace(3)
	c.State[0][0][2], c.State[0][1][2], c.State[0][2][2], c.State[2][0][2], c.State[2][1][2], c.State[2][2][2], c.State[5][0][2], c.State[5][1][2], c.State[5][2][2], c.State[4][2][0], c.State[4][1][0], c.State[4][0][0] = c.State[4][2][0], c.State[4][1][0], c.State[4][0][0], c.State[0][0][2], c.State[0][1][2], c.State[0][2][2], c.State[2][0][2], c.State[2][1][2], c.State[2][2][2], c.State[5][0][2], c.State[5][1][2], c.State[5][2][2]
}

func (c *Cube) ApplyL() {
	c.RotateFace(1)
	c.State[0][0][0], c.State[0][1][0], c.State[0][2][0], c.State[2][0][0], c.State[2][1][0], c.State[2][2][0], c.State[5][0][0], c.State[5][1][0], c.State[5][2][0], c.State[4][2][2], c.State[4][1][2], c.State[4][0][2] = c.State[4][2][2], c.State[4][1][2], c.State[4][0][2], c.State[0][0][0], c.State[0][1][0], c.State[0][2][0], c.State[2][0][0], c.State[2][1][0], c.State[2][2][0], c.State[5][0][0], c.State[5][1][0], c.State[5][2][0]
}

func (c *Cube) ApplyLPrime() {
	c.RotateFace(1)
	c.RotateFace(1)
	c.RotateFace(1)
	c.State[0][0][0], c.State[0][1][0], c.State[0][2][0], c.State[2][0][0], c.State[2][1][0], c.State[2][2][0], c.State[5][0][0], c.State[5][1][0], c.State[5][2][0], c.State[4][2][2], c.State[4][1][2], c.State[4][0][2] = c.State[2][0][0], c.State[2][1][0], c.State[2][2][0], c.State[5][0][0], c.State[5][1][0], c.State[5][2][0], c.State[4][2][2], c.State[4][1][2], c.State[4][0][2], c.State[0][0][0], c.State[0][1][0], c.State[0][2][0]
}

func (c *Cube) ApplyX() {
	c.RotateFace(1)
	c.RotateFace(1)
	c.RotateFace(1)

	c.RotateFace(3)

	c.RotateFace(0)
	c.RotateFace(0)
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
	c.RotateFace(5)
	c.RotateFace(5)

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

	for _, f := range []int{0, 1, 3, 5} {
		c.RotateFace(f)
	}

	c.State[0], c.State[1], c.State[5], c.State[3] = c.State[1], c.State[5], c.State[3], c.State[0]
}

func (c *Cube) ApplyZPrime() {
	c.RotateFace(2)
	c.RotateFace(2)
	c.RotateFace(2)

	c.RotateFace(4)

	for _, f := range []int{0, 1, 3, 5} {
		for range []int{0, 1, 2} {
			c.RotateFace(f)
		}
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

	case "F":
		c.ApplyF()
	case "F'":
		c.ApplyFPrime()
	case "F2":
		c.ApplyF()
		c.ApplyF()

	case "B":
		c.ApplyB()
	case "B'":
		c.ApplyBPrime()
	case "B2":
		c.ApplyB()
		c.ApplyB()

	case "R":
		c.ApplyR()
	case "R'":
		c.ApplyRPrime()
	case "R2":
		c.ApplyR()
		c.ApplyR()

	case "L":
		c.ApplyL()
	case "L'":
		c.ApplyLPrime()
	case "L2":
		c.ApplyL()
		c.ApplyL()

	case "Uw":
		c.ApplyD()
		c.ApplyY()
	case "Uw'":
		c.ApplyDPrime()
		c.ApplyYPrime()
	case "Uw2":
		c.ApplyD()
		c.ApplyY()
		c.ApplyD()
		c.ApplyY()

	case "Dw":
		c.ApplyU()
		c.ApplyYPrime()
	case "Dw'":
		c.ApplyUPrime()
		c.ApplyY()
	case "Dw2":
		c.ApplyU()
		c.ApplyYPrime()
		c.ApplyU()
		c.ApplyYPrime()

	case "Fw":
		c.ApplyB()
		c.ApplyZ()
	case "Fw'":
		c.ApplyBPrime()
		c.ApplyZPrime()
	case "Fw2":
		c.ApplyB()
		c.ApplyZ()
		c.ApplyB()
		c.ApplyZ()

	case "Bw":
		c.ApplyF()
		c.ApplyZPrime()
	case "Bw'":
		c.ApplyFPrime()
		c.ApplyZ()
	case "Bw2":
		c.ApplyF()
		c.ApplyZPrime()
		c.ApplyF()
		c.ApplyZPrime()

	case "Rw":
		c.ApplyL()
		c.ApplyX()
	case "Rw'":
		c.ApplyLPrime()
		c.ApplyXPrime()
	case "Rw2":
		c.ApplyL()
		c.ApplyX()
		c.ApplyL()
		c.ApplyX()

	case "Lw":
		c.ApplyR()
		c.ApplyXPrime()
	case "Lw'":
		c.ApplyRPrime()
		c.ApplyX()
	case "Lw2":
		c.ApplyR()
		c.ApplyXPrime()
		c.ApplyR()
		c.ApplyXPrime()

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

func (c *Cube) OfficialSolutionLength() int {
	moveCount := 0

	for _, move := range strings.Split(c.Solution, " ") {
		idx := IndexFunc(ROTATIONS, func(rot string) bool { return rot == move })
		if idx == -1 {
			moveCount++
		}
	}

	return moveCount * 1000
}

func (c *Cube) TotalSolutionLength() int {
	return len(strings.Split(c.Solution, " "))
}

func (c *Cube) PrintState() {
	for i := range c.State[0] {
		fmt.Print("   ")
		for j := range c.State[0][i] {
			fmt.Print(COLORS[c.State[0][i][j]])
		}
		fmt.Println()
	}

	for i := range []int{0, 1, 2} {
		for _, face := range []int{1, 2, 3, 4} {
			for j := range []int{0, 1, 2} {
				fmt.Print(COLORS[c.State[face][i][j]])
			}
		}
		fmt.Println()
	}

	for i := range c.State[5] {
		fmt.Print("   ")
		for j := range c.State[5][i] {
			fmt.Print(COLORS[c.State[5][i][j]])
		}
		fmt.Println()
	}
}

func ParseFMCSolutionToMilliseconds(scramble string, solution string) int {
	c := Cube{InitialState(), scramble, solution}

	if !c.ValidMoves() {
		return constants.DNF
	}

	c.ApplyScramble()
	c.ApplySolution()

	if !c.Solved() || c.TotalSolutionLength() > 80 {
		return constants.DNF
	}

	return c.OfficialSolutionLength()
}
