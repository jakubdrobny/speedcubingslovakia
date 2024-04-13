package models

type ScrambleSet struct {
	Event CompetitionEvent `json:"event"`
	Scrambles []string `json:"scrambles"`
}

func (s *ScrambleSet) AddScramble(scramble string) {
	s.Scrambles = append(s.Scrambles, scramble)
}