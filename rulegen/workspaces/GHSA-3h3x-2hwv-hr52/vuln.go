package main


func (h *boringHMAC) Write(p []byte) (int, error) {
	if len(p) > 0 {
		C._goboringcrypto_HMAC_Update(h.ctx, (*C.uint8_t)(unsafe.Pointer(&p[0])), C.size_t(len(p)))
	}
	runtime.KeepAlive(h)
	return len(p), nil
}

func (h *boringHMAC) Sum(in []byte) []byte {
	if h.sum == nil {
		size := h.Size()
		h.sum = make([]byte, size)
	}
	C._goboringcrypto_HMAC_Final(h.ctx, (*C.uint8_t)(unsafe.Pointer(&h.sum[0])), nil)
	return append(in, h.sum...)
}
