package main


func (h *boringHMAC) Write(p []byte) (int, error) {
	if len(p) > 0 {
		if C._goboringcrypto_HMAC_Update(h.ctx, (*C.uint8_t)(unsafe.Pointer(&p[0])), C.size_t(len(p))) == 0 {
			panic("boringcrypto: HMAC_Update failed")
		}
	}
	runtime.KeepAlive(h)
	return len(p), nil
}

func (h *boringHMAC) Sum(in []byte) []byte {
	size := h.Size()
	if h.sum == nil {
		h.sum = make([]byte, size)
	}
	if C._goboringcrypto_HMAC_Final(h.ctx, (*C.uint8_t)(unsafe.Pointer(&h.sum[0])), C.uint(size)) == 0 {
		panic("boringcrypto: HMAC_Final failed")
	}
	return append(in, h.sum...)
}
