/*
 * Copyright (c) 2024 AlertAvert.com. All rights reserved.
 */

// Author: M. Massenzio (marco@alertavert.com), 5/30/24

package pkg

import (
	"fmt"
)

func Simple(name string) error {
	fmt.Println(fmt.Sprintf("Name is %d characters long", len(name)))
	return nil
}
