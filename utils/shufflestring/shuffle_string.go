package shufflestring

import (
	"math/rand"
	"time"
)

func Shuffle(in string) string {
	rand.Seed(time.Now().Unix())
	inRune := []rune(in)
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}
