// Apple Lossless (ALAC) decoder
package alac

import (
	"fmt"
	"unsafe"
)

/*
#include <stdint.h>
#include "alac.h"
*/
import "C"

// type Alac *C.struct_alac_file
type Alac struct {
	// alac_file unsafe.Pointer
	alac_file *C.struct_alac_file
}

// New alac decoder. Sample size 16, 2 chan!
func New() (*Alac, error) {
	a := C.create_alac(16, 2)
	if a == nil {
		return nil, fmt.Errorf("can't create alac. No idea why, though")
	}
	// TODO: fmtp stuff
	// fmtp: 96 352 0 16 40 10 14 2 255 0 0 44100
	a.setinfo_max_samples_per_frame = 352 // frame_size;
	a.setinfo_7a = 0                      // fmtp[2];
	a.setinfo_sample_size = 16            // sample_size;
	a.setinfo_rice_historymult = 40       // fmtp[4];
	a.setinfo_rice_initialhistory = 10    // fmtp[5];
	a.setinfo_rice_kmodifier = 14         // fmtp[6];
	a.setinfo_7f = 2                      // fmtp[7];
	a.setinfo_80 = 255                    // fmtp[8];
	a.setinfo_82 = 0                      // fmtp[9];
	a.setinfo_86 = 0                      // fmtp[10];
	a.setinfo_8a_rate = 44100             // fmtp[11];

	C.allocate_buffers(a)
	al := Alac{alac_file: a}
	return &al, nil
}

func (a *Alac) Decode(b []byte) []byte {
	d := make([]byte, 10*1024) // a whole 10K. TODO :)
	var l C.int = 0
	// TODO: how do these arguments work?!?
	C.decode_frame(a.alac_file, (*C.uchar)(&b[0]), unsafe.Pointer(&d[0]), &l)
	return d[:l]
}
