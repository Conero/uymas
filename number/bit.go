package number

import (
	"fmt"
	"math"
)

// BitSize is an example like `time.Duration`
// The Bytes unit for resource file.
type BitSize int64

// 1 Bytes = 8Bit					byte
//
// 1 B 	  = 8Bit					byte			字节
//
// 1 KB    = 1024 Bytes			    Kilobyte		千字节
//
// 1 MB	  = 1024 KB				    Megabyte		百万字节		   兆
//
// 1 GB    = 1024 MB				Gigabyte		千兆			   吉
//
// 1 TB	  = 1024 GB				    Terabyte		万亿字节		   太
//
// 1 PB	  = 1024 TB				    Petabyte		千万亿字节	   拍
//
// 1 EB	  = 1024 PB				    Exabyte			百亿亿字节	   艾
//
// 1 ZB	  = 1024 EB				    Zettabyte		十万亿亿字节	   泽
//
// 1 YB	  = 1024 ZB				    Yottabyte		一亿亿亿字节	   尧
//
// 1 BB	  = 1024 YB				    Brontobyte
//
// 1 NB	  = 1024 BB				    NonaByte
//
// 1 DB	  = 1024 NB				    DoggaByte
const (
	Bit  BitSize = 1
	Byte         = 8 * Bit
	KB           = 1000 * Byte //kilobyte
	MB           = 1000 * KB   //megabyte
	GB           = 1000 * MB   //gigabyte
	TB           = 1000 * GB   //terabyte
	PB           = 1000 * TB   //petabyte
	//EB           = 1000 * PB   //exabyte
	//ZB           = 1000 * EB   //zettabyte
	//YB           = 1000 * ZB   //yottabyte

	KiB = 1024 * Byte // kibibyte
	MiB = 1024 * KiB  //	mebibyte
	GiB = 1024 * MiB  //	gibibyte
	TiB = 1024 * GiB  //	tebibyte
	PiB = 1024 * TiB  //	pebibyte
	//EiB = 1024 * PiB  //	exbibyte
	//ZiB = 1024 * EiB  //	zebibyte
	//YiB = 1024 * ZiB  //yobibyte
)

// Format get the format of byte size
func (b BitSize) Format() (float64, string) {
	if b == 0 {
		return 0, "bit"
	}
	// Byte
	if b < Byte {
		return float64(b), "bit"
	}
	// Byte
	if b < KB {
		return float64(b) / float64(Byte), "Byte"
	}
	var sizes = []string{"", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	var i = math.Floor(math.Log10(float64(b)) / math.Log10(1000))
	//the max data unit is to `YB`
	var sizesLen = float64(len(sizes))
	if i > sizesLen {
		i = sizesLen - 1
	}
	return float64(b) / math.Pow(1000, i), sizes[int(i)]
}

// Format2 get the format of byte size
func (b BitSize) Format2() (float64, string) {
	if b == 0 {
		return 0, "bit"
	}
	// Byte
	if b < Byte {
		return float64(b), "bit"
	}
	// Byte
	if b < KiB {
		return float64(b) / float64(Byte), "Byte"
	}
	var sizes = []string{"", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	var i = math.Floor(math.Log(float64(b)/float64(Byte)) / math.Log(1024))
	//the max data unit is to `YB`
	var sizesLen = float64(len(sizes))
	if i > sizesLen {
		i = sizesLen - 1
	}
	return float64(b/Byte) / math.Pow(1024, i), sizes[int(i)]
}

func (b BitSize) Bit() float64 {
	return float64(b)
}

func (b BitSize) Byte() float64 {
	return float64(b / Byte)
}

func (b BitSize) KB() float64 {
	return float64(b / KB)
}

func (b BitSize) MB() float64 {
	return float64(b / MB)
}

func (b BitSize) GB() float64 {
	return float64(b / GB)
}

func (b BitSize) TB() float64 {
	return float64(b / TB)
}

func (b BitSize) PB() float64 {
	return float64(b / PB)
}

func (b BitSize) KiB() float64 {
	return float64(b / KiB)
}

func (b BitSize) MiB() float64 {
	return float64(b / MiB)
}

func (b BitSize) GiB() float64 {
	return float64(b / GiB)
}

func (b BitSize) TiB() float64 {
	return float64(b / TiB)
}

func (b BitSize) PiB() float64 {
	return float64(b / PiB)
}

// The file size default use 1024.
func (b BitSize) String() string {
	v, unit := b.Format2()
	if v == 0 {
		return fmt.Sprintf("%v %v", v, unit)
	}
	return fmt.Sprintf("%.4f %v", v, unit)
}

// Bytes get the bit size by bytes
func Bytes(bytes int64) BitSize {
	return BitSize(bytes) * Byte
}
