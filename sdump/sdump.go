package sdump

/*
#cgo CFLAGS: -I.
#include <stdio.h>
#include <stdlib.h>
extern void *image_from_data(unsigned char *data, int size);
extern void image_crop(void *handle, int x, int y, int w, int h);
extern void *sixel_output2(void *handle, int *output_size);
*/
import "C"
import "unsafe"

func EncodeCropImage(data []byte, x, y, w, h int) []byte {
	var i C.int = C.int(len(data))
	handle := C.image_from_data((*C.uchar)(unsafe.Pointer(&data[0])), i)
	C.image_crop(handle, C.int(x), C.int(y), C.int(w), C.int(h))

	p := C.sixel_output2(handle, &i)
	defer C.free(p)

	buf := C.GoBytes(p, i)

	return buf
}
