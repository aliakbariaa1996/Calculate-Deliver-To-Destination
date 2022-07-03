package delivery

import (
	"math"
)

type (
	SourceLocation struct {
		Lat float64
		Lng float64
	}
	DeliverManLocation struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
)

func (s *UseCase) GetDistance(souLoc SourceLocation, deliLoc []DeliverManLocation) interface{} {
	c := make(chan float64)
	for i := 0; i < len(deliLoc); i++ {
		go s.CalculateDist(souLoc.Lat, souLoc.Lng, deliLoc[i].Lat, deliLoc[i].Lng, c)
	}
	var output []float64
	for l := 0; l < len(deliLoc); l++ {
		output = append(output, <-c)
	}
	close(c)
	return output
}

func (s *UseCase) CalculateDist(sourceX float64, sourceY float64, DeliverManX float64, DeliverManY float64, c chan float64) {
	radSourceX := math.Pi * sourceX / 180
	radDeliverManX := math.Pi * DeliverManX / 180
	theta := sourceY - DeliverManY
	radTheTa := math.Pi * theta / 180

	dist := math.Sin(radSourceX)*math.Sin(radDeliverManX) + math.Cos(radSourceX)*math.Cos(radDeliverManX)*math.Cos(radTheTa)
	if dist > 1 {
		dist = 1
	}
	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344

	c <- dist
}
