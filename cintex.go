package croot

/*
 #include "croot/croot.h"

 #include <stdlib.h>
 #include <string.h>

*/
import "C"

// placeholder for Cintex
type cintex int

const Cintex = cintex(0)

func (c cintex) Enable() {
	C.CRoot_Cintex_Enable()
}

func (c cintex) SetDebug(lvl int) {
	C.CRoot_Cintex_SetDebug(C.int(lvl))
}

func init() {
	Cintex.Enable()
	//Cintex.SetDebug(100000)
}

// eof
