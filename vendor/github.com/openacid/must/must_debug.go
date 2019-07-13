// +build debug

package must

import (
	"github.com/openacid/must/enabled"
)

var (
	// Be is the container of all checking APIs, such as "must.Be.Equal(a, b)".
	//
	// Since 0.1.0
	Be = enabled.Be
)
