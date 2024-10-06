package randoms

import (
	crand "crypto/rand"
	"log"
	"math/big"
	mrand "math/rand"

	"github.com/pkg/errors"
)

const DEBUG_MODE = false

func UniformIntN(n int) int {
	if DEBUG_MODE {
		return mrand.Intn(n)
	} else {
		maxInt := big.NewInt(int64(n))
		randomBigInt, err := crand.Int(crand.Reader, maxInt)
		if err != nil {
			log.Panicln(errors.Wrap(err, "Could not call RNG for URandIntN"))
		}
		return int(randomBigInt.Int64())
	}
}
