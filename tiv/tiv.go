package tiv

/*
#cgo CFLAGS: -I.
#cgo CXXFLAGS: -std=c++11
#include <stdio.h>
#include <stdlib.h>
extern void* encode2data(const char *filename, int *size);
extern void* encode2datafromdata(const unsigned char *data, int *size);
extern void* crop_image(const unsigned char *mem, int *size,
  const int x, const int y, const int w, const int h);
*/
import "C"
import "unsafe"

func Encode2Data(filename string) []byte {
	fn := C.CString(filename)
	defer C.free(unsafe.Pointer(fn))

	var i C.int
	p := C.encode2data(fn, &i)
	defer C.free(p)

	buf := C.GoBytes(p, i)

	return buf
}

func EncodeFromData(data []byte) []byte {
	var i C.int = C.int(len(data))
	p := C.encode2datafromdata((*C.uchar)(unsafe.Pointer(&data[0])), &i)
	defer C.free(p)

	buf := C.GoBytes(p, i)

	return buf

}
