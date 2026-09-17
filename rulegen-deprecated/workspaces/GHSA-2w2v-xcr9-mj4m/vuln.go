package main

import (
	"errors"
	"fmt"
	pkg1_raft "github.com/hashicorp/raft"
	codec1978 "github.com/ugorji/go/codec"
	"reflect"
	"runtime"
	time "time"
		panic(err)
	}
	if false { // reference the types, but skip this branch at build/run time
		var v0 pkg1_raft.ServerID
		var v1 time.Duration
		_, _ = v0, v1
	}
}

func (x Bitmap) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	}
}

func (x DiffType) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	z.DecSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x *RaftServer) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
		} else {
			yysep2 := !z.EncBinary()
			yy2arr2 := z.EncBasicHandle().StructToArray
			var yyq2 [1]bool
			_, _, _ = yysep2, yyq2, yy2arr2
			const yyr2 bool = false
			var yynn2 int
			if yyr2 || yy2arr2 {
				r.EncodeArrayStart(1)
			} else {
				yynn2 = 1
				for _, b := range yyq2 {
					if b {
						yynn2++
					r.EncodeBool(bool(x.SystemSchedulerEnabled))
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayEnd100)
			} else {
					*((*bool)(yyv4)) = r.DecodeBool()
				}
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj6 int
	var yyb6 bool
	var yyhl6 bool = l >= 0
	yyj6++
	if yyhl6 {
		yyb6 = yyj6 > l
	} else {
		yyb6 = r.CheckBreak()
	}
	if yyb6 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd100)
		return
	}
	if r.TryDecodeAsNil() {
		x.SystemSchedulerEnabled = false
	} else {
		yyv7 := &x.SystemSchedulerEnabled
		yym8 := z.DecBinary()
		_ = yym8
		if false {
		} else {
			*((*bool)(yyv7)) = r.DecodeBool()
		}
	}
	for {
		yyj6++
		if yyhl6 {
			yyb6 = yyj6 > l
		} else {
			yyb6 = r.CheckBreak()
		}
		if yyb6 {
			break
		}
		z.DecSendContainerState(codecSelfer_containerArrayElem100)
		z.DecStructFieldNotFound(yyj6-1, "")
	}
	z.DecSendContainerState(codecSelfer_containerArrayEnd100)
}
	z.DecSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) encBitmap(v Bitmap, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeStringBytes(codecSelferC_RAW100, []byte(v))
}

func (x codecSelfer100) decBitmap(v *Bitmap, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	*v = r.DecodeBytes(*((*[]byte)(v)), false, false)
}

func (x codecSelfer100) encSlicePtrtoFieldDiff(v []*FieldDiff, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeArrayStart(len(v))
	for _, yyv1 := range v {
		z.EncSendContainerState(codecSelfer_containerArrayElem100)
		if yyv1 == nil {
			r.EncodeNil()
		} else {
			yyv1.CodecEncodeSelf(e)
		}
	}
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decSlicePtrtoFieldDiff(v *[]*FieldDiff, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*FieldDiff{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
			yyc1 = true
		}
	} else if yyl1 > 0 {
		var yyrr1, yyrl1 int
		var yyrt1 bool
		_, _ = yyrl1, yyrt1
		yyrr1 = yyl1 // len(yyv1)
		if yyl1 > cap(yyv1) {

			yyrg1 := len(yyv1) > 0
			yyv21 := yyv1
			yyrl1, yyrt1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 8)
			if yyrt1 {
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*FieldDiff, yyrl1)
				}
			} else {
				yyv1 = make([]*FieldDiff, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			if yyrg1 {
				copy(yyv1, yyv21)
			}
		} else if yyl1 != len(yyv1) {
			yyv1 = yyv1[:yyl1]
			yyc1 = true
		}
		yyj1 := 0
		for ; yyj1 < yyrr1; yyj1++ {
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = FieldDiff{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(FieldDiff)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
			}

		}
		if yyrt1 {
			for ; yyj1 < yyl1; yyj1++ {
				yyv1 = append(yyv1, nil)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = FieldDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(FieldDiff)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
				}

			}
		}

	} else {
		yyj1 := 0
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *FieldDiff
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = FieldDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(FieldDiff)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
				}

			} else {
				z.DecSwallow()
			}

		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*FieldDiff{}
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}

func (x codecSelfer100) encSlicePtrtoObjectDiff(v []*ObjectDiff, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeArrayStart(len(v))
	for _, yyv1 := range v {
		z.EncSendContainerState(codecSelfer_containerArrayElem100)
		if yyv1 == nil {
			r.EncodeNil()
		} else {
			yyv1.CodecEncodeSelf(e)
		}
	}
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decSlicePtrtoObjectDiff(v *[]*ObjectDiff, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*ObjectDiff{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
			yyc1 = true
		}
	} else if yyl1 > 0 {
		var yyrr1, yyrl1 int
		var yyrt1 bool
		_, _ = yyrl1, yyrt1
		yyrr1 = yyl1 // len(yyv1)
		if yyl1 > cap(yyv1) {

			yyrg1 := len(yyv1) > 0
			yyv21 := yyv1
			yyrl1, yyrt1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 8)
			if yyrt1 {
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*ObjectDiff, yyrl1)
				}
			} else {
				yyv1 = make([]*ObjectDiff, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			if yyrg1 {
				copy(yyv1, yyv21)
			}
		} else if yyl1 != len(yyv1) {
			yyv1 = yyv1[:yyl1]
			yyc1 = true
		}
		yyj1 := 0
		for ; yyj1 < yyrr1; yyj1++ {
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = ObjectDiff{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(ObjectDiff)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
			}

		}
		if yyrt1 {
			for ; yyj1 < yyl1; yyj1++ {
				yyv1 = append(yyv1, nil)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = ObjectDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(ObjectDiff)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
				}

			}
		}

	} else {
		yyj1 := 0
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *ObjectDiff
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = ObjectDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(ObjectDiff)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
				}

			} else {
				z.DecSwallow()
			}

		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*ObjectDiff{}
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}

func (x codecSelfer100) encSlicePtrtoTaskGroupDiff(v []*TaskGroupDiff, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeArrayStart(len(v))
	for _, yyv1 := range v {
		z.EncSendContainerState(codecSelfer_containerArrayElem100)
		if yyv1 == nil {
			r.EncodeNil()
		} else {
			yyv1.CodecEncodeSelf(e)
		}
	}
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decSlicePtrtoTaskGroupDiff(v *[]*TaskGroupDiff, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*TaskGroupDiff{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*TaskGroupDiff, yyrl1)
				}
			} else {
				yyv1 = make([]*TaskGroupDiff, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = TaskGroupDiff{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(TaskGroupDiff)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = TaskGroupDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(TaskGroupDiff)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *TaskGroupDiff
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = TaskGroupDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(TaskGroupDiff)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*TaskGroupDiff{}
			yyc1 = true
		}
	}
	}
}

func (x codecSelfer100) encSlicePtrtoTaskDiff(v []*TaskDiff, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decSlicePtrtoTaskDiff(v *[]*TaskDiff, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*TaskDiff{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*TaskDiff, yyrl1)
				}
			} else {
				yyv1 = make([]*TaskDiff, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = TaskDiff{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(TaskDiff)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = TaskDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(TaskDiff)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *TaskDiff
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = TaskDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(TaskDiff)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
				}

			} else {
				z.DecSwallow()
			}

		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*TaskDiff{}
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}

func (x codecSelfer100) encTaskGroupDiffs(v TaskGroupDiffs, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decTaskGroupDiffs(v *TaskGroupDiffs, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*TaskGroupDiff{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*TaskGroupDiff, yyrl1)
				}
			} else {
				yyv1 = make([]*TaskGroupDiff, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = TaskGroupDiff{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(TaskGroupDiff)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = TaskGroupDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(TaskGroupDiff)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *TaskGroupDiff
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = TaskGroupDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(TaskGroupDiff)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*TaskGroupDiff{}
			yyc1 = true
		}
	}
	}
}

func (x codecSelfer100) encTaskDiffs(v TaskDiffs, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeArrayStart(len(v))
	for _, yyv1 := range v {
		z.EncSendContainerState(codecSelfer_containerArrayElem100)
		if yyv1 == nil {
			r.EncodeNil()
		} else {
			yyv1.CodecEncodeSelf(e)
		}
	}
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decTaskDiffs(v *TaskDiffs, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*TaskDiff{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
			yyc1 = true
		}
	} else if yyl1 > 0 {
		var yyrr1, yyrl1 int
		var yyrt1 bool
		_, _ = yyrl1, yyrt1
		yyrr1 = yyl1 // len(yyv1)
		if yyl1 > cap(yyv1) {

			yyrg1 := len(yyv1) > 0
			yyv21 := yyv1
			yyrl1, yyrt1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 8)
			if yyrt1 {
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*TaskDiff, yyrl1)
				}
			} else {
				yyv1 = make([]*TaskDiff, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			if yyrg1 {
				copy(yyv1, yyv21)
			}
		} else if yyl1 != len(yyv1) {
			yyv1 = yyv1[:yyl1]
			yyc1 = true
		}
		yyj1 := 0
		for ; yyj1 < yyrr1; yyj1++ {
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = TaskDiff{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(TaskDiff)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
			}

		}
		if yyrt1 {
			for ; yyj1 < yyl1; yyj1++ {
				yyv1 = append(yyv1, nil)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = TaskDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(TaskDiff)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
				}

			}
		}

	} else {
		yyj1 := 0
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *TaskDiff
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = TaskDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(TaskDiff)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
				}

			} else {
				z.DecSwallow()
			}

		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*TaskDiff{}
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}

func (x codecSelfer100) encObjectDiffs(v ObjectDiffs, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decObjectDiffs(v *ObjectDiffs, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*ObjectDiff{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*ObjectDiff, yyrl1)
				}
			} else {
				yyv1 = make([]*ObjectDiff, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = ObjectDiff{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(ObjectDiff)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = ObjectDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(ObjectDiff)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *ObjectDiff
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = ObjectDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(ObjectDiff)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*ObjectDiff{}
			yyc1 = true
		}
	}
	}
}

func (x codecSelfer100) encFieldDiffs(v FieldDiffs, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decFieldDiffs(v *FieldDiffs, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*FieldDiff{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*FieldDiff, yyrl1)
				}
			} else {
				yyv1 = make([]*FieldDiff, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = FieldDiff{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(FieldDiff)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = FieldDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(FieldDiff)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *FieldDiff
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = FieldDiff{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(FieldDiff)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*FieldDiff{}
			yyc1 = true
		}
	}
	}
}

func (x codecSelfer100) encSlicePtrtoNetworkResource(v []*NetworkResource, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
		if yyv1 == nil {
			r.EncodeNil()
		} else {
			yym2 := z.EncBinary()
			_ = yym2
			if false {
			} else if z.HasExtensions() && z.EncExt(yyv1) {
			} else {
				z.EncFallback(yyv1)
			}
		}
	}
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decSlicePtrtoNetworkResource(v *[]*NetworkResource, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*NetworkResource{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*NetworkResource, yyrl1)
				}
			} else {
				yyv1 = make([]*NetworkResource, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = NetworkResource{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(NetworkResource)
				}
				yyw2 := yyv1[yyj1]
				yym3 := z.DecBinary()
				_ = yym3
				if false {
				} else if z.HasExtensions() && z.DecExt(yyw2) {
				} else {
					z.DecFallback(yyw2, false)
				}
			}

		}
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = NetworkResource{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(NetworkResource)
					}
					yyw4 := yyv1[yyj1]
					yym5 := z.DecBinary()
					_ = yym5
					if false {
					} else if z.HasExtensions() && z.DecExt(yyw4) {
					} else {
						z.DecFallback(yyw4, false)
					}
				}

			}
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *NetworkResource
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = NetworkResource{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(NetworkResource)
					}
					yyw6 := yyv1[yyj1]
					yym7 := z.DecBinary()
					_ = yym7
					if false {
					} else if z.HasExtensions() && z.DecExt(yyw6) {
					} else {
						z.DecFallback(yyw6, false)
					}
				}

			} else {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*NetworkResource{}
			yyc1 = true
		}
	}
	}
}

func (x codecSelfer100) encMapstringBitmap(v map[string]Bitmap, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	z.EncSendContainerState(codecSelfer_containerMapEnd100)
}

func (x codecSelfer100) decMapstringBitmap(v *map[string]Bitmap, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yyl1 := r.ReadMapStart()
	yybh1 := z.DecBasicHandle()
	if yyv1 == nil {
		yyrl1, _ := z.DecInferLen(yyl1, yybh1.MaxInitLen, 40)
		yyv1 = make(map[string]Bitmap, yyrl1)
		*v = yyv1
	}
	var yymk1 string
	var yymv1 Bitmap
	var yymg1 bool
	if yybh1.MapValueReset {
		yymg1 = true
	}
				}
			}

			if yymg1 {
				yymv1 = yyv1[yymk1]
			} else {
				yymv1 = nil
			}
			z.DecSendContainerState(codecSelfer_containerMapValue100)
			if r.TryDecodeAsNil() {
				yymv1 = nil
			} else {
				yyv4 := &yymv1
				yyv4.CodecDecodeSelf(d)
			}

			if yyv1 != nil {
				yyv1[yymk1] = yymv1
			}
		}
				}
			}

			if yymg1 {
				yymv1 = yyv1[yymk1]
			} else {
				yymv1 = nil
			}
			z.DecSendContainerState(codecSelfer_containerMapValue100)
			if r.TryDecodeAsNil() {
				yymv1 = nil
			} else {
				yyv7 := &yymv1
				yyv7.CodecDecodeSelf(d)
			}

			if yyv1 != nil {
				yyv1[yymk1] = yymv1
			}
		}
	z.DecSendContainerState(codecSelfer_containerMapEnd100)
}

func (x codecSelfer100) encSlicePtrtoRaftServer(v []*RaftServer, e *codec1978.Encoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	z.EncSendContainerState(codecSelfer_containerArrayEnd100)
}

func (x codecSelfer100) decSlicePtrtoRaftServer(v *[]*RaftServer, d *codec1978.Decoder) {
	var h codecSelfer100
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*RaftServer{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*RaftServer, yyrl1)
				}
			} else {
				yyv1 = make([]*RaftServer, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = RaftServer{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(RaftServer)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = RaftServer{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(RaftServer)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *RaftServer
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = RaftServer{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(RaftServer)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*RaftServer{}
			yyc1 = true
		}
	}
	"io"
	"os"
	"path/filepath"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/nomad/helper"
	"github.com/hashicorp/nomad/helper/tlsutil"
	"github.com/hashicorp/nomad/nomad/structs/config"
	"github.com/mitchellh/mapstructure"
)

// ParseConfigFile parses the given path as a config file.
func ParseConfigFile(path string) (*Config, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
		return nil, err
	}
	defer f.Close()

	config, err := ParseConfig(f)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// ParseConfig parses the config from the given io.Reader.
//
// Due to current internal limitations, the entire contents of the
// io.Reader will be copied into memory first before parsing.
func ParseConfig(r io.Reader) (*Config, error) {
	// Copy the reader into an in-memory buffer first since HCL requires it.
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}

	// Parse the buffer
	root, err := hcl.Parse(buf.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing: %s", err)
	}
	buf.Reset()

	// Top-level item should be a list
	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: root should be an object")
	}

	var config Config
	if err := parseConfig(&config, list); err != nil {
		return nil, fmt.Errorf("error parsing 'config': %v", err)
	}

	return &config, nil
}

func parseConfig(result *Config, list *ast.ObjectList) error {
	// Check for invalid keys
	valid := []string{
		"region",
		"datacenter",
		"name",
		"data_dir",
		"plugin_dir",
		"log_level",
		"log_json",
		"bind_addr",
		"enable_debug",
		"ports",
		"addresses",
		"interfaces",
		"advertise",
		"client",
		"server",
		"telemetry",
		"leave_on_interrupt",
		"leave_on_terminate",
		"enable_syslog",
		"syslog_facility",
		"disable_update_check",
		"disable_anonymous_signature",
		"consul",
		"vault",
		"tls",
		"http_api_response_headers",
		"acl",
		"sentinel",
		"autopilot",
		"plugin",
	}
	if err := helper.CheckHCLKeys(list, valid); err != nil {
		return multierror.Prefix(err, "config:")
	}

	// Decode the full thing into a map[string]interface for ease
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, list); err != nil {
		return err
	}
	delete(m, "ports")
	delete(m, "addresses")
	delete(m, "interfaces")
	delete(m, "advertise")
	delete(m, "client")
	delete(m, "server")
	delete(m, "telemetry")
	delete(m, "consul")
	delete(m, "vault")
	delete(m, "tls")
	delete(m, "http_api_response_headers")
	delete(m, "acl")
	delete(m, "sentinel")
	delete(m, "autopilot")
	delete(m, "plugin")

	// Decode the rest
	if err := mapstructure.WeakDecode(m, result); err != nil {
		return err
	}

	// Parse ports
	if o := list.Filter("ports"); len(o.Items) > 0 {
		if err := parsePorts(&result.Ports, o); err != nil {
			return multierror.Prefix(err, "ports ->")
		}
	}

	// Parse addresses
	if o := list.Filter("addresses"); len(o.Items) > 0 {
		if err := parseAddresses(&result.Addresses, o); err != nil {
			return multierror.Prefix(err, "addresses ->")
		}
	}

	// Parse advertise
	if o := list.Filter("advertise"); len(o.Items) > 0 {
		if err := parseAdvertise(&result.AdvertiseAddrs, o); err != nil {
			return multierror.Prefix(err, "advertise ->")
		}
	}

	// Parse client config
	if o := list.Filter("client"); len(o.Items) > 0 {
		if err := parseClient(&result.Client, o); err != nil {
			return multierror.Prefix(err, "client ->")
		}
	}

	// Parse server config
	if o := list.Filter("server"); len(o.Items) > 0 {
		if err := parseServer(&result.Server, o); err != nil {
			return multierror.Prefix(err, "server ->")
		}
	}

	// Parse ACL config
	if o := list.Filter("acl"); len(o.Items) > 0 {
		if err := parseACL(&result.ACL, o); err != nil {
			return multierror.Prefix(err, "acl ->")
		}
	}

	// Parse telemetry config
	if o := list.Filter("telemetry"); len(o.Items) > 0 {
		if err := parseTelemetry(&result.Telemetry, o); err != nil {
			return multierror.Prefix(err, "telemetry ->")
		}
	}

	// Parse the consul config
	if o := list.Filter("consul"); len(o.Items) > 0 {
		if err := parseConsulConfig(&result.Consul, o); err != nil {
			return multierror.Prefix(err, "consul ->")
		}
	}

	// Parse the vault config
	if o := list.Filter("vault"); len(o.Items) > 0 {
		if err := parseVaultConfig(&result.Vault, o); err != nil {
			return multierror.Prefix(err, "vault ->")
		}
	}

	// Parse the TLS config
	if o := list.Filter("tls"); len(o.Items) > 0 {
		if err := parseTLSConfig(&result.TLSConfig, o); err != nil {
			return multierror.Prefix(err, "tls ->")
		}
	}

	// Parse Sentinel config
	if o := list.Filter("sentinel"); len(o.Items) > 0 {
		if err := parseSentinel(&result.Sentinel, o); err != nil {
			return multierror.Prefix(err, "sentinel->")
		}
	}

	// Parse Autopilot config
	if o := list.Filter("autopilot"); len(o.Items) > 0 {
		if err := parseAutopilot(&result.Autopilot, o); err != nil {
			return multierror.Prefix(err, "autopilot->")
		}
	}

	// Parse Plugin configs
	if o := list.Filter("plugin"); len(o.Items) > 0 {
		if err := parsePlugins(&result.Plugins, o); err != nil {
			return multierror.Prefix(err, "plugin->")
		}
	}

	// Parse out http_api_response_headers fields. These are in HCL as a list so
	// we need to iterate over them and merge them.
	if headersO := list.Filter("http_api_response_headers"); len(headersO.Items) > 0 {
		for _, o := range headersO.Elem().Items {
			var m map[string]interface{}
			if err := hcl.DecodeObject(&m, o.Val); err != nil {
				return err
			}
			if err := mapstructure.WeakDecode(m, &result.HTTPAPIResponseHeaders); err != nil {
				return err
			}
		}
	}

	return nil
}

func parsePorts(result **Ports, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'ports' block allowed")
	}

	// Get our ports object
	listVal := list.Items[0].Val

	// Check for invalid keys
	valid := []string{
		"http",
		"rpc",
		"serf",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	var ports Ports
	if err := mapstructure.WeakDecode(m, &ports); err != nil {
		return err
	}
	*result = &ports
	return nil
}

func parseAddresses(result **Addresses, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'addresses' block allowed")
	}

	// Get our addresses object
	listVal := list.Items[0].Val

	// Check for invalid keys
	valid := []string{
		"http",
		"rpc",
		"serf",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	var addresses Addresses
	if err := mapstructure.WeakDecode(m, &addresses); err != nil {
		return err
	}
	*result = &addresses
	return nil
}

func parseAdvertise(result **AdvertiseAddrs, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'advertise' block allowed")
	}

	// Get our advertise object
	listVal := list.Items[0].Val

	// Check for invalid keys
	valid := []string{
		"http",
		"rpc",
		"serf",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	var advertise AdvertiseAddrs
	if err := mapstructure.WeakDecode(m, &advertise); err != nil {
		return err
	}
	*result = &advertise
	return nil
}

func parseClient(result **ClientConfig, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'client' block allowed")
	}

	// Get our client object
	obj := list.Items[0]

	// Value should be an object
	var listVal *ast.ObjectList
	if ot, ok := obj.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return fmt.Errorf("client value: should be an object")
	}

	// Check for invalid keys
	valid := []string{
		"enabled",
		"state_dir",
		"alloc_dir",
		"servers",
		"node_class",
		"options",
		"meta",
		"chroot_env",
		"network_interface",
		"network_speed",
		"memory_total_mb",
		"cpu_total_compute",
		"max_kill_timeout",
		"client_max_port",
		"client_min_port",
		"reserved",
		"stats",
		"gc_interval",
		"gc_disk_usage_threshold",
		"gc_inode_usage_threshold",
		"gc_parallel_destroys",
		"gc_max_allocs",
		"no_host_uuid",
		"server_join",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	delete(m, "options")
	delete(m, "meta")
	delete(m, "chroot_env")
	delete(m, "reserved")
	delete(m, "stats")
	delete(m, "server_join")

	var config ClientConfig
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &config,
	})
	if err != nil {
		return err
	}
	if err := dec.Decode(m); err != nil {
		return err
	}

	// Parse out options fields. These are in HCL as a list so we need to
	// iterate over them and merge them.
	if optionsO := listVal.Filter("options"); len(optionsO.Items) > 0 {
		for _, o := range optionsO.Elem().Items {
			var m map[string]interface{}
			if err := hcl.DecodeObject(&m, o.Val); err != nil {
				return err
			}
			if err := mapstructure.WeakDecode(m, &config.Options); err != nil {
				return err
			}
		}
	}

	// Parse out options meta. These are in HCL as a list so we need to
	// iterate over them and merge them.
	if metaO := listVal.Filter("meta"); len(metaO.Items) > 0 {
		for _, o := range metaO.Elem().Items {
			var m map[string]interface{}
			if err := hcl.DecodeObject(&m, o.Val); err != nil {
				return err
			}
			if err := mapstructure.WeakDecode(m, &config.Meta); err != nil {
				return err
			}
		}
	}

	// Parse out chroot_env fields. These are in HCL as a list so we need to
	// iterate over them and merge them.
	if chrootEnvO := listVal.Filter("chroot_env"); len(chrootEnvO.Items) > 0 {
		for _, o := range chrootEnvO.Elem().Items {
			var m map[string]interface{}
			if err := hcl.DecodeObject(&m, o.Val); err != nil {
				return err
			}
			if err := mapstructure.WeakDecode(m, &config.ChrootEnv); err != nil {
				return err
			}
		}
	}

	// Parse reserved config
	if o := listVal.Filter("reserved"); len(o.Items) > 0 {
		if err := parseReserved(&config.Reserved, o); err != nil {
			return multierror.Prefix(err, "reserved ->")
		}
	}

	// Parse ServerJoin config
	if o := listVal.Filter("server_join"); len(o.Items) > 0 {
		if err := parseServerJoin(&config.ServerJoin, o); err != nil {
			return multierror.Prefix(err, "server_join->")
		}
	}

	*result = &config
	return nil
}

func parseReserved(result **Resources, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'reserved' block allowed")
	}

	// Get our reserved object
	obj := list.Items[0]

	// Value should be an object
	var listVal *ast.ObjectList
	if ot, ok := obj.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return fmt.Errorf("client value: should be an object")
	}

	// Check for invalid keys
	valid := []string{
		"cpu",
		"memory",
		"disk",
		"reserved_ports",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	var reserved Resources
	if err := mapstructure.WeakDecode(m, &reserved); err != nil {
		return err
	}
	if err := reserved.CanParseReserved(); err != nil {
		return err
	}

	*result = &reserved
	return nil
}

func parseServer(result **ServerConfig, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'server' block allowed")
	}

	// Get our server object
	obj := list.Items[0]

	// Value should be an object
	var listVal *ast.ObjectList
	if ot, ok := obj.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return fmt.Errorf("client value: should be an object")
	}

	// Check for invalid keys
	valid := []string{
		"enabled",
		"bootstrap_expect",
		"data_dir",
		"protocol_version",
		"raft_protocol",
		"num_schedulers",
		"enabled_schedulers",
		"node_gc_threshold",
		"eval_gc_threshold",
		"job_gc_threshold",
		"deployment_gc_threshold",
		"heartbeat_grace",
		"min_heartbeat_ttl",
		"max_heartbeats_per_second",
		"rejoin_after_leave",
		"encrypt",
		"authoritative_region",
		"non_voting_server",
		"redundancy_zone",
		"upgrade_version",

		"server_join",

		// For backwards compatibility
		"start_join",
		"retry_join",
		"retry_max",
		"retry_interval",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	delete(m, "server_join")

	var config ServerConfig
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &config,
	})
	if err != nil {
		return err
	}
	if err := dec.Decode(m); err != nil {
		return err
	}

	if config.UpgradeVersion != "" {
		if _, err := version.NewVersion(config.UpgradeVersion); err != nil {
			return fmt.Errorf("error parsing upgrade_version: %v", err)
		}
	}

	// Parse ServerJoin config
	if o := listVal.Filter("server_join"); len(o.Items) > 0 {
		if err := parseServerJoin(&config.ServerJoin, o); err != nil {
			return multierror.Prefix(err, "server_join->")
		}
	}

	*result = &config
	return nil
}

func parseServerJoin(result **ServerJoin, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'server_join' block allowed")
	}

	// Get our object
	listVal := list.Items[0].Val

	// Check for invalid keys
	valid := []string{
		"start_join",
		"retry_join",
		"retry_max",
		"retry_interval",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	var serverJoinInfo ServerJoin
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &serverJoinInfo,
	})
	if err != nil {
		return err
	}
	if err := dec.Decode(m); err != nil {
		return err
	}

	*result = &serverJoinInfo
	return nil
}

func parseACL(result **ACLConfig, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'acl' block allowed")
	}

	// Get our server object
	obj := list.Items[0]

	// Value should be an object
	var listVal *ast.ObjectList
	if ot, ok := obj.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return fmt.Errorf("acl value: should be an object")
	}

	// Check for invalid keys
	valid := []string{
		"enabled",
		"token_ttl",
		"policy_ttl",
		"replication_token",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	var config ACLConfig
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &config,
	})
	if err != nil {
		return err
	}
	if err := dec.Decode(m); err != nil {
		return err
	}

	*result = &config
	return nil
}

func parseTelemetry(result **Telemetry, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'telemetry' block allowed")
	}

	// Get our telemetry object
	listVal := list.Items[0].Val

	// Check for invalid keys
	valid := []string{
		"statsite_address",
		"statsd_address",
		"disable_hostname",
		"use_node_name",
		"collection_interval",
		"publish_allocation_metrics",
		"publish_node_metrics",
		"datadog_address",
		"datadog_tags",
		"prometheus_metrics",
		"circonus_api_token",
		"circonus_api_app",
		"circonus_api_url",
		"circonus_submission_interval",
		"circonus_submission_url",
		"circonus_check_id",
		"circonus_check_force_metric_activation",
		"circonus_check_instance_id",
		"circonus_check_search_tag",
		"circonus_check_display_name",
		"circonus_check_tags",
		"circonus_broker_id",
		"circonus_broker_select_tag",
		"disable_tagged_metrics",
		"backwards_compatible_metrics",
		"prefix_filter",
		"filter_default",
		"disable_dispatched_job_summary_metrics",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	var telemetry Telemetry
	if err := mapstructure.WeakDecode(m, &telemetry); err != nil {
		return err
	}
	if telemetry.CollectionInterval != "" {
		if dur, err := time.ParseDuration(telemetry.CollectionInterval); err != nil {
			return fmt.Errorf("error parsing value of %q: %v", "collection_interval", err)
		} else {
			telemetry.collectionInterval = dur
		}
	}
	*result = &telemetry
	return nil
}

func parseConsulConfig(result **config.ConsulConfig, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'consul' block allowed")
	}

	// Get our Consul object
	listVal := list.Items[0].Val

	// Check for invalid keys
	valid := []string{
		"address",
		"auth",
		"auto_advertise",
		"ca_file",
		"cert_file",
		"checks_use_advertise",
		"client_auto_join",
		"client_service_name",
		"client_http_check_name",
		"key_file",
		"server_auto_join",
		"server_service_name",
		"server_http_check_name",
		"server_serf_check_name",
		"server_rpc_check_name",
		"ssl",
		"timeout",
		"token",
		"verify_ssl",
	}

	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	consulConfig := config.DefaultConsulConfig()
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &consulConfig,
	})
	if err != nil {
		return err
	}
	if err := dec.Decode(m); err != nil {
		return err
	}

	*result = consulConfig
	return nil
}

func parseTLSConfig(result **config.TLSConfig, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'tls' block allowed")
	}

	// Get the TLS object
	listVal := list.Items[0].Val

	valid := []string{
		"http",
		"rpc",
		"verify_server_hostname",
		"rpc_upgrade_mode",
		"ca_file",
		"cert_file",
		"key_file",
		"verify_https_client",
		"tls_cipher_suites",
		"tls_min_version",
		"tls_prefer_server_cipher_suites",
	}

	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	var tlsConfig config.TLSConfig
	if err := mapstructure.WeakDecode(m, &tlsConfig); err != nil {
		return err
	}

	if _, err := tlsutil.ParseCiphers(&tlsConfig); err != nil {
		return err
	}

	if _, err := tlsutil.ParseMinVersion(tlsConfig.TLSMinVersion); err != nil {
		return err
	}

	*result = &tlsConfig
	return nil
}

func parseVaultConfig(result **config.VaultConfig, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'vault' block allowed")
	}

	// Get our Vault object
	listVal := list.Items[0].Val

	// Check for invalid keys
	valid := []string{
		"address",
		"allow_unauthenticated",
		"enabled",
		"task_token_ttl",
		"ca_file",
		"ca_path",
		"cert_file",
		"create_from_role",
		"key_file",
		"tls_server_name",
		"tls_skip_verify",
		"token",
	}

	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	vaultConfig := config.DefaultVaultConfig()
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &vaultConfig,
	})
	if err != nil {
		return err
	}
	if err := dec.Decode(m); err != nil {
		return err
	}

	*result = vaultConfig
	return nil
}

func parseSentinel(result **config.SentinelConfig, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'sentinel' block allowed")
	}

	// Get our sentinel object
	obj := list.Items[0]

	// Value should be an object
	var listVal *ast.ObjectList
	if ot, ok := obj.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	} else {
		return fmt.Errorf("sentinel value: should be an object")
	}

	// Check for invalid keys
	valid := []string{
		"import",
	}
	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var config config.SentinelConfig
	if err := hcl.DecodeObject(&config, listVal); err != nil {
		return err
	}

	*result = &config
	return nil
}

func parseAutopilot(result **config.AutopilotConfig, list *ast.ObjectList) error {
	list = list.Elem()
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'autopilot' block allowed")
	}

	// Get our Autopilot object
	listVal := list.Items[0].Val

	// Check for invalid keys
	valid := []string{
		"cleanup_dead_servers",
		"server_stabilization_time",
		"last_contact_threshold",
		"max_trailing_logs",
		"enable_redundancy_zones",
		"disable_upgrade_migration",
		"enable_custom_upgrades",
	}

	if err := helper.CheckHCLKeys(listVal, valid); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, listVal); err != nil {
		return err
	}

	autopilotConfig := config.DefaultAutopilotConfig()
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           &autopilotConfig,
	})
	if err != nil {
		return err
	}
	if err := dec.Decode(m); err != nil {
		return err
	}

	*result = autopilotConfig
	return nil
}

func parsePlugins(result *[]*config.PluginConfig, list *ast.ObjectList) error {
	listLen := len(list.Items)
	plugins := make([]*config.PluginConfig, listLen)

	// Check for invalid keys
	valid := []string{
		"args",
		"config",
	}

	for i := 0; i < listLen; i++ {
		// Get the current plugin object
		listVal := list.Items[i]

		// Deal with json->hcl AST parsing incorrectness when directly nested
		// items show up as additional keys.
		// TODO(preetha): Add additional tests and fix other places that have the same issue
		unwrapLegacyHCLObjectKeysFromJSON(listVal, 1)
		if err := helper.CheckHCLKeys(listVal.Val, valid); err != nil {
			return fmt.Errorf("invalid keys in plugin config %d: %v", i+1, err)
		}

		// Ensure there is a key
		if len(listVal.Keys) != 1 {
			return fmt.Errorf("plugin config %d doesn't incude a name key", i+1)
		}

		var plugin config.PluginConfig
		if err := hcl.DecodeObject(&plugin, listVal); err != nil {
			return fmt.Errorf("error decoding plugin config %d: %v", i+1, err)
		}

		plugins[i] = &plugin
	}

	*result = plugins
	return nil
}

// unwrapLegacyHCLObjectKeysFromJSON cleans up an edge case that can occur when
// parsing JSON as input: if we're parsing JSON then directly nested
// items will show up as additional "keys".
//
// For objects that expect a fixed number of keys, this breaks the
// decoding process. This function unwraps the object into what it would've
// looked like if it came directly from HCL by specifying the number of keys
// you expect.
//
// Example:
//
// { "foo": { "baz": {} } }
//
// Will show up with Keys being: []string{"foo", "baz"}
// when we really just want the first two. This function will fix this.
func unwrapLegacyHCLObjectKeysFromJSON(item *ast.ObjectItem, depth int) {
	if len(item.Keys) > depth && item.Keys[0].Token.JSON {
		for len(item.Keys) > depth {
			// Pop off the last key
			n := len(item.Keys)
			key := item.Keys[n-1]
			item.Keys[n-1] = nil
			item.Keys = item.Keys[:n-1]

			// Wrap our value in a list
			item.Val = &ast.ObjectType{
				List: &ast.ObjectList{
					Items: []*ast.ObjectItem{
						{
							Keys: []*ast.ObjectKey{key},
							Val:  item.Val,
						},
					},
				},
			}
		}
	}
}
	z.DecSendContainerState(codecSelfer_containerArrayEnd102)
}

func (x *AllocStatsRequest) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer102
	z, r := codec1978.GenHelperEncoder(e)
