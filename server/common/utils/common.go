package utils

import (
	"context"
	"strings"
)

func GetPayload(ctx context.Context) (username, role string) {
	payload := ctx.Value("payloads").(map[string]interface{})
	username = payload["username"].(string)
	role = payload["role"].(string)
	return
}

func GetS3BucketAndPath(inputpath string) (bucket, path string) {
	// path必须是<bucket>/开头，如果第一位是/就去掉它
	inputpath = strings.TrimPrefix(inputpath, "/")
	pathParts := strings.Split(inputpath, "/")
	bucket = pathParts[0]
	path = strings.TrimPrefix(inputpath, bucket)
	return
}
