package s3fs

import (
	"fmt"
	"io/ioutil"
	"log"
)

type S3Options struct {
	URL             string      `json:"url,omitempty"`
	AccessKeyID     Credentials `json:"access_key_id,omitempty"`
	AccessKeySecret Credentials `json:"access_key_secret,omitempty"`
	Bucket          string      `json:"bucket"`
	PathStyle       bool        `json:"path_style"`
	NoMixUpload     bool        `json:"no_mix_upload"`
	NoMultipart     bool        `json:"no_multipart"`
	NoCache         bool        `json:"no_cache"`
	Debug           bool        `json:"debug"`
}

// Args returns an array of S3FS arguments for the current set of mount options
func (opt *S3Options) Args(path string) []string {
	args := []string{
		"-f",
		opt.Bucket,
		path,
	}

	// debug mode
	if opt.Debug {
		args = append(args, "-o", "dbglevel=info", "-o", "curldbg")
	}

	// cached mode (enabled by default)
	if !opt.NoCache {
		args = append(args, "-o", "use_cache=/tmp")
	}

	// custom url
	if opt.URL != "" {
		args = append(args, "-o", fmt.Sprintf("url=%s", opt.URL))
	}

	// path request style
	if opt.PathStyle {
		args = append(args, "-o", "use_path_request_style")
	}

	// no mix upload
	if opt.NoMixUpload {
		args = append(args, "-o", "nomixupload")
	}

	// no multipart upload
	if opt.NoMultipart {
		args = append(args, "-o", "nomultipart")
	}

	return args
}

func (opt *S3Options) writeCredentials() error {
	// neither are set, do nothing
	if opt.AccessKeyID == "" && opt.AccessKeySecret == "" {
		return nil
	}

	// if one is set, return an error
	if opt.AccessKeyID == "" || opt.AccessKeySecret == "" {
		return ErrCredentialsMissing
	}

	// write credentials to file
	creds := []byte(fmt.Sprintf("%s:%s", opt.AccessKeyID, opt.AccessKeySecret))
	err := ioutil.WriteFile("/etc/passwd-s3fs", creds, 0600)
	if err != nil {
		return fmt.Errorf("Error writing s3fs credentials: %s", err)
	}

	log.Println("s3fs: wrote credentials for", opt.Bucket)
	return nil
}
