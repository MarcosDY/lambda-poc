package secret

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/golang/protobuf/proto"
	"github.com/spiffe/go-spiffe/v2/proto/spiffe/workload"
)

var region = os.Getenv("AWS_REGION")

type SecretManager interface {
	GetSecret(ctx context.Context, secretID string) (*workload.X509SVIDResponse, error)
}

type secret struct {
	// TODO: add an interface
	client *secretsmanager.SecretsManager
}

func New() SecretManager {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: &region,
	}))

	return &secret{
		client: secretsmanager.New(sess),
	}
}

func (s *secret) GetSecret(ctx context.Context, secretID string) (*workload.X509SVIDResponse, error) {
	startAt := time.Now()
	defer func() {
		elapse := time.Since(startAt)
		log.Printf("[EXTENSION] GetSecretValue takes: %s", strconv.FormatInt(elapse.Milliseconds(), 10))
	}()

	// Get AWSCURRENT secrert
	resp, err := s.client.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &secretID,
	})
	if err != nil {
		return nil, err
	}

	x509SVID := new(workload.X509SVIDResponse)
	if err := proto.Unmarshal(resp.SecretBinary, x509SVID); err != nil {
		log.Printf("[EXTENSION] failed to unmarshal: %v", err)
		return nil, err
	}

	return x509SVID, nil
}
