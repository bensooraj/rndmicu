package graph

import (
	"fmt"

	"github.com/bensooraj/rndmicu/data/models"
)

func starryNightC(cs *[]models.Creator) []*models.Creator {
	csStar := []*models.Creator{}
	for i := 0; i < len(*cs); i++ {
		csStar = append(csStar, &((*cs)[i]))
	}
	return csStar
}

func starryNightA(cs *[]models.AudioShort) []*models.AudioShort {
	asStar := []*models.AudioShort{}
	for i := 0; i < len(*cs); i++ {
		asStar = append(asStar, &((*cs)[i]))
	}
	return asStar
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
