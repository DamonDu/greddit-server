package dataloader

import "fmt"

// Int64Key implements the Key interface for a string
type Int64Key int64

// String is an identity method. Used to implement String interface
func (k Int64Key) String() string { return fmt.Sprintf("%d", k) }

func (k Int64Key) Raw() interface{} { return k }
