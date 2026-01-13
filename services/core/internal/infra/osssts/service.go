package osssts

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

	"moxuevideo/core/internal/domain"
)

type Service struct {
	client           *sts.Client
	region           string
	bucket           string
	endpoint         string
	roleARN          string
	roleSessionName  string
	durationSeconds  int
	assumeRoleScheme string
}

func New(accessKeyID, accessKeySecret, region, endpoint, bucket, roleARN, roleSessionName string, durationSeconds int) (*Service, error) {
	client, err := sts.NewClientWithAccessKey(region, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	if durationSeconds <= 0 {
		durationSeconds = 900
	}
	if roleSessionName == "" {
		roleSessionName = "moxuevideo-core"
	}
	return &Service{
		client:           client,
		region:           region,
		bucket:           bucket,
		endpoint:         endpoint,
		roleARN:          roleARN,
		roleSessionName:  roleSessionName,
		durationSeconds:  durationSeconds,
		assumeRoleScheme: "https",
	}, nil
}

func (s *Service) GetUploadToken(_ context.Context, purpose string, userID uint64) (domain.UploadToken, error) {
	prefix := buildDir(purpose, userID)
	policy, err := buildPolicy(s.bucket, prefix)
	if err != nil {
		return domain.UploadToken{}, err
	}

	req := sts.CreateAssumeRoleRequest()
	req.Scheme = s.assumeRoleScheme
	req.RoleArn = s.roleARN
	req.RoleSessionName = fmt.Sprintf("%s-%d-%s", s.roleSessionName, userID, randSuffix(6))
	req.DurationSeconds = requests.NewInteger(s.durationSeconds)
	req.Policy = policy

	resp, err := s.client.AssumeRole(req)
	if err != nil {
		return domain.UploadToken{}, err
	}

	host := s.endpoint
	if host != "" && !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}
	if s.bucket != "" && host != "" {
		host = strings.Replace(host, "https://", "https://"+s.bucket+".", 1)
		host = strings.Replace(host, "http://", "http://"+s.bucket+".", 1)
	}

	return domain.UploadToken{
		Region:   s.region,
		Endpoint: s.endpoint,
		Bucket:   s.bucket,
		Host:     host,
		Dir:      prefix,
		Creds: domain.Creds{
			AccessKeyID:     resp.Credentials.AccessKeyId,
			AccessKeySecret: resp.Credentials.AccessKeySecret,
			SecurityToken:   resp.Credentials.SecurityToken,
			Expiration:      resp.Credentials.Expiration,
		},
	}, nil
}

func buildDir(purpose string, userID uint64) string {
	p := strings.ToLower(strings.TrimSpace(purpose))
	switch p {
	case "avatar", "avatars", "user_avatar":
		return fmt.Sprintf("avatars/%d/", userID)
	case "video", "videos":
		return fmt.Sprintf("videos/%d/", userID)
	default:
		return fmt.Sprintf("misc/%d/", userID)
	}
}

func buildPolicy(bucket, prefix string) (string, error) {
	if bucket == "" {
		return "", fmt.Errorf("oss bucket required")
	}
	if prefix == "" {
		prefix = "misc/"
	}
	resource := fmt.Sprintf("acs:oss:*:*:%s/%s*", bucket, prefix)

	p := map[string]any{
		"Version": "1",
		"Statement": []any{
			map[string]any{
				"Effect":   "Allow",
				"Action":   []string{"oss:PutObject", "oss:AbortMultipartUpload", "oss:ListParts"},
				"Resource": []string{resource},
			},
		},
	}

	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func randSuffix(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)[:n]
}
