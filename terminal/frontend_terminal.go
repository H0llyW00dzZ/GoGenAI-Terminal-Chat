// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"time"
)

func PrintPrefixWithTimeStamp(prefix string) {
	currentTime := time.Now().Format(TimeFormat)
	fmt.Printf(ObjectHighLevelString, currentTime, prefix)
}
