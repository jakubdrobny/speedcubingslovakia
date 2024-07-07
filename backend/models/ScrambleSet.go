package models

type ScrambleSet struct {
	Event CompetitionEvent `json:"event"`
	Scrambles []Scramble `json:"scrambles"`
}

type Scramble struct {
	Scramble string `json:"scramble"`
	Img string `json:"img"`
}

func (s *ScrambleSet) AddScramble(scramble Scramble) {
	s.Scrambles = append(s.Scrambles, scramble)
}