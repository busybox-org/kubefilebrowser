package utils

import "encoding/base64"

func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base64Decode(s string) ([]byte, error) {
	ds, err := base64.StdEncoding.DecodeString(s)
	return ds, err
}

func Base64UrlEncode(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

func Base64UrlDecode(s string) ([]byte, error) {
	ds, err := base64.URLEncoding.DecodeString(s)
	return ds, err
}
