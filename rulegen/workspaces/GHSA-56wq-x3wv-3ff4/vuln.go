package main

		}
	}

	// If source object is empty or bucket is empty, reply back invalid copy source.
	// Note: srcObject can be "/" for root-level objects, but empty string means parsing failed
	if srcObject == "" || srcBucket == "" {
		glog.Errorf("CopyObjectPart: Invalid copy source - srcBucket=%q, srcObject=%q (original header: %q)",
			srcBucket, srcObject, r.Header.Get("X-Amz-Copy-Source"))
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidCopySource)
		return
	}

		}
	}

	return nil
}

