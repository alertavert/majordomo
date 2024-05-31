/*
 * Copyright (c) 2024 AlertAvert.com. All rights reserved.
 */

package main

import (
	"fmt"
	"sample/pkg"
	"time"
)

func main() {
	fmt.Println("This is a wonderful world!")
	fmt.Println("Today's date is:", time.Now().Format("2006-01-02"))
	pkg.Simple("Marco")
}
