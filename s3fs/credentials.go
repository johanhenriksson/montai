package s3fs

// Credentials store sensitive strings that should not be marshalled to JSON
type Credentials string

// MarshalJSON implementation that ensures credentials are never returned from the API
func (c Credentials) MarshalJSON() ([]byte, error) {
	if len(c) > 0 {
		return []byte(`"[redacted]"`), nil
	}
	return []byte(`""`), nil
}
