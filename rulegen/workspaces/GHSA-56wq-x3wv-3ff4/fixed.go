package main

		}
	}

	// Validate the copy source as CopyObject does.
	if err := ValidateCopySource(cpSrcPath, srcBucket, srcObject); err != nil {
		glog.V(2).Infof("CopyObjectPartHandler validation error: %v", err)
		s3err.WriteErrorResponse(w, r, MapCopyValidationError(err))
		return
	}

		}
	}

	// `.`/`..` segments are collapsed by the filer's path join; reject them as
	// IsValidObjectKey does for the request URL so the source stays in-bucket.
	if !s3_constants.IsValidBucketName(srcBucket) || !s3_constants.IsValidObjectKey(srcObject) {
		return &CopyValidationError{
			Code:    s3err.ErrInvalidCopySource,
			Message: "Copy source contains invalid path segments",
		}
	}

	return nil
}

