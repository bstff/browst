/*
 * Copyright (c) 2014 Hayaki Saito
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

#include <stddef.h>  /* for size_t */

#ifndef LIBSIXEL_SIXEL_H
#define LIBSIXEL_SIXEL_H

#define LIBSIXEL_VRTSION 1.0.3
#define LIBSIXEL_ABI_VRTSION 1:0:0

#define SIXEL_OUTPUT_PACKET_SIZE 1024
#define SIXEL_PALETTE_MAX 256

/* output character size */
enum characterSize {
    CSIZE_7BIT = 0,  /* 7bit character */
    CSIZE_8BIT = 1,  /* 8bit character */
};

/* method for finding the largest dimention for splitting,
 * and sorting by that component */
enum methodForLargest {
    LARGE_AUTO = 0,  /* choose automatically the method for finding the largest
                        dimention */
    LARGE_NORM = 1,  /* simply comparing the range in RGB space */
    LARGE_LUM  = 2   /* transforming into luminosities before the comparison */
};

/* method for choosing a color from the box */
enum methodForRep {
    REP_AUTO           = 0, /* choose automatically the method for selecting
                               representative color from each box */
    REP_CENTER_BOX     = 1, /* choose the center of the box */
    REP_AVERAGE_COLORS = 2, /* choose the average all the color
                               in the box (specified in Heckbert's paper) */
    REP_AVERAGE_PIXELS = 3  /* choose the averate all the pixels in the box */
};

/* method for diffusing */
enum methodForDiffuse {
    DIFFUSE_AUTO     = 0, /* choose diffusion type automatically */
    DIFFUSE_NONE     = 1, /* don't diffuse */
    DIFFUSE_ATKINSON = 2, /* diffuse with Bill Atkinson's method */
    DIFFUSE_FS       = 3, /* diffuse with Floyd-Steinberg method */
    DIFFUSE_JAJUNI   = 4, /* diffuse with Jarvis, Judice & Ninke method */
    DIFFUSE_STUCKI   = 5, /* diffuse with Stucki's method */
    DIFFUSE_BURKES   = 6  /* diffuse with Burkes' method */
};

/* quality modes */
enum qualityMode {
    QUALITY_AUTO = 0, /* choose quality mode automatically */
    QUALITY_HIGH = 1, /* high quality */
    QUALITY_LOW  = 2, /* low quality */
};

/* built-in dither */
enum builtinDither {
    BUILTIN_MONO_DARK  = 0, /* monochrome terminal with dark background */
    BUILTIN_MONO_LIGHT = 1, /* monochrome terminal with dark background */
    BUILTIN_XTERM16    = 2, /* xterm 16color */
    BUILTIN_XTERM256   = 3, /* xterm 256color */
};

#ifndef LIBSIXEL_DITHER_H

/* dither context object */
typedef struct sixel_dither {} sixel_dither_t;

#endif /* SIXEL_DITHER_H */

#ifndef LIBSIXEL_OUTPUT_H

typedef int (* sixel_write_function)(char *data, int size, void *priv);

typedef struct sixel_output {} sixel_output_t;

#endif  /* LIBSIXEL_OUTPUT_H */


/* output context manipulation API */

#ifdef __cplusplus
extern "C" {
#endif

/* create output context object */
sixel_output_t *const
sixel_output_create(sixel_write_function /* in */ fn_write, /* callback function for output sixel */
                    void /* in */ *priv);                   /* private data given as 
                                                               3rd argument of fn_write */
/* destroy output context object */
void
sixel_output_destroy(sixel_output_t /* in */ *output); /* output context */

/* increment reference count of output context object (thread-unsafe) */
void
sixel_output_ref(sixel_output_t /* in */ *output);     /* output context */

/* decrement reference count of output context object (thread-unsafe) */
void
sixel_output_unref(sixel_output_t /* in */ *output);   /* output context */

int
sixel_output_get_8bit_availability(sixel_output_t *output);

void
sixel_output_set_8bit_availability(sixel_output_t *output, int availability);

#ifdef __cplusplus
}
#endif


/* color quantization API */

#ifdef __cplusplus
extern "C" {
#endif

/* create dither context object */
sixel_dither_t *
sixel_dither_create(int /* in */ ncolors); /* number of colors */

/* get built-in dither context object */
sixel_dither_t *
sixel_dither_get(int builtin_dither); /* ID of built-in dither object */

/* destroy dither context object */
void
sixel_dither_destroy(sixel_dither_t *dither); /* dither context object */

/* increment reference count of dither context object (thread-unsafe) */
void
sixel_dither_ref(sixel_dither_t *dither); /* dither context object */

/* decrement reference count of dither context object (thread-unsafe) */
void
sixel_dither_unref(sixel_dither_t *dither); /* dither context object */

/* initialize internal palette from specified pixel buffer */
int
sixel_dither_initialize(sixel_dither_t *dither,          /* dither context object */
                        unsigned char /* in */ *data,    /* sample image */
                        int /* in */ width,              /* image width */
                        int /* in */ height,             /* image height */
                        int /* in */ depth,              /* pixel depth, now only '3' is supported */
                        int /* in */ method_for_largest, /* method for finding the largest dimention */
                        int /* in */ method_for_rep,     /* method for choosing a color from the box */
                        int /* in */ quality_mode);      /* quality of histgram processing */

/* set diffusion type, choose from enum methodForDiffuse */
void
sixel_dither_set_diffusion_type(sixel_dither_t /* in */ *dither,  /* dither context object */
                                int /* in */ method_for_diffuse); /* one of enum methodForDiffuse */

/* get number of palette colors */
int
sixel_dither_get_num_of_palette_colors(sixel_dither_t /* in */ *dither);   /* dither context object */

/* get number of histgram colors */
int
sixel_dither_get_num_of_histgram_colors(sixel_dither_t /* in */ *dither);  /* dither context object */

/* get palette */
unsigned char *
sixel_dither_get_palette(sixel_dither_t /* in */ *dither);  /* dither context object */

#ifdef __cplusplus
}
#endif

/* converter API */

typedef void * (* sixel_allocator_function)(size_t size);

#ifdef __cplusplus
extern "C" {
#endif

/* convert pixels into sixel format and write it to output context */
int
sixel_encode(unsigned char  /* in */ *pixels,     /* pixel bytes */
             int            /* in */  width,      /* image width */
             int            /* in */  height,     /* image height */
             int            /* in */  depth,      /* pixel depth */
             sixel_dither_t /* in */ *dither,     /* dither context */
             sixel_output_t /* in */ *context);   /* output context */

/* convert sixel data into indexed pixel bytes and palette data */
int
sixel_decode(unsigned char              /* in */  *sixels,    /* sixel bytes */
             int                        /* in */  size,       /* size of sixel bytes */
             unsigned char              /* out */ **pixels,   /* decoded pixels */
             int                        /* out */ *pwidth,    /* image width */
             int                        /* out */ *pheight,   /* image height */
             unsigned char              /* out */ **palette,  /* RGBA palette */
             int                        /* out */ *ncolors,   /* palette size (<= 256) */
             sixel_allocator_function   /* out */ allocator); /* malloc function */

#ifdef __cplusplus
}
#endif

#endif  /* LIBSIXEL_SIXEL_H */

/* emacs, -*- Mode: C; tab-width: 4; indent-tabs-mode: nil -*- */
/* vim: set expandtab ts=4 : */
/* EOF */
