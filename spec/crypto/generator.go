package crypto

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	LowerChars  = "abcdefghijkmnopqrstuvwxyz" // excluded l
	UpperChars  = "ABCDEFGHJKLMNPQRSTUVWXYZ" // excluded I, O
	DigitChars  = "23456789"                  // excluded 0, 1
	SymbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

var defaultWordList = []string{
	"galaxy", "nebula", "cipher", "matrix", "falcon", "beacon", "summit", "timber",
	"canyon", "voyage", "zenith", "crystal", "harbor", "plasma", "shield", "quasar",
	"vector", "shadow", "aurora", "breeze", "anchor", "delta", "frost", "granite",
	"hazard", "ignite", "jungle", "meteor", "orbit", "prism", "quartz", "radar",
	"strata", "titan", "umbra", "vertex", "vortex", "walrus", "yonder", "zephyr",
	"alpine", "badger", "cobalt", "dynamo", "echo", "fusion", "glider", "horizon",
	"island", "jaguar", "kinetic", "legacy", "monolith", "nexus", "oasis", "pulsar",
	"quantum", "relic", "solace", "tracker", "uranium", "viper", "wizard", "xenon",
}

type PasswordOptions struct {
	Length      int  `json:"length"`
	IncludeUpper bool `json:"includeUpper"`
	IncludeLower bool `json:"includeLower"`
	IncludeDigits bool `json:"includeDigits"`
	IncludeSymbols bool `json:"includeSymbols"`
	Passphrase  bool `json:"passphrase"`
	WordCount   int  `json:"wordCount"`
	Separator   string `json:"separator"`
}

func GeneratePassword(opts PasswordOptions) string {
	if opts.Passphrase {
		count := opts.WordCount
		if count < 3 {
			count = 4
		}
		sep := opts.Separator
		if sep == "" {
			sep = "-"
		}
		words := make([]string, count)
		for i := 0; i < count; i++ {
			idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(defaultWordList))))
			if err != nil {
				words[i] = "vault"
			} else {
				words[i] = defaultWordList[idx.Int64()]
			}
		}
		return strings.Join(words, sep)
	}

	length := opts.Length
	if length < 8 {
		length = 16
	}

	var charset string
	var mandatory []byte

	if opts.IncludeLower {
		charset += LowerChars
		if b, err := randomChar(LowerChars); err == nil {
			mandatory = append(mandatory, b)
		}
	}
	if opts.IncludeUpper {
		charset += UpperChars
		if b, err := randomChar(UpperChars); err == nil {
			mandatory = append(mandatory, b)
		}
	}
	if opts.IncludeDigits {
		charset += DigitChars
		if b, err := randomChar(DigitChars); err == nil {
			mandatory = append(mandatory, b)
		}
	}
	if opts.IncludeSymbols {
		charset += SymbolChars
		if b, err := randomChar(SymbolChars); err == nil {
			mandatory = append(mandatory, b)
		}
	}

	if charset == "" {
		charset = LowerChars + UpperChars + DigitChars
	}

	res := make([]byte, length)
	// Place mandatory chars first
	copy(res, mandatory)

	for i := len(mandatory); i < length; i++ {
		ch, err := randomChar(charset)
		if err != nil {
			res[i] = 'x'
		} else {
			res[i] = ch
		}
	}

	// Shuffle
	for i := length - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err == nil {
			j := int(jBig.Int64())
			res[i], res[j] = res[j], res[i]
		}
	}

	return string(res)
}

func randomChar(chars string) (byte, error) {
	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, err
	}
	return chars[idx.Int64()], nil
}

