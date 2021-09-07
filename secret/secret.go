package secret

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var region = os.Getenv("AWS_REGION")

type Svid struct {
	// The SPIFFE ID of that identify this SVID
	SpiffeID string `json:"spiffeId,omitempty"`
	// PEM encoded certificate chain. MAY invlude intermediates,
	// the leaf certificate (or SVID itself) MUST come first
	X509Svid string `json:"x509Svid,omitempty"`
	// PEM encoded PKCS#8 private key.
	X509SvidKey string `json:"x509SvidKey,omitempty"`
	// PEM encoded X.509 bundle for the trust domain
	Bundle string `json:"bundle,omitempty"`
	// CA certificate bundles belonging to foreign trust domains that the workload should trust,
	// keyed by trust domain. Bundles are in encoded in PEM format.
	FederatedBundles map[string]string `json:"federated_bundles,omitempty"`
}

type SecretManager interface {
	GetSecret(ctx context.Context, secretID string) (*Svid, error)
}

type secret struct {
	// TODO: add an interface
	client *secretsmanager.Client
}

func New(ctx context.Context) (SecretManager, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &secret{
		client: secretsmanager.NewFromConfig(cfg),
	}, nil
}

func (s *secret) GetSecret(ctx context.Context, secretID string) (*Svid, error) {
	startAt := time.Now()
	defer func() {
		elapse := time.Since(startAt)
		log.Printf("[EXTENSION] GetSecretValue takes: %s", strconv.FormatInt(elapse.Milliseconds(), 10))
	}()

	// Get the specified secret ID (AWSCURRENT version)
	resp, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretID,
	})
	if err != nil {
		return nil, err
	}

	secretBinary := new(Svid)
	if err := json.Unmarshal(resp.SecretBinary, secretBinary); err != nil {
		log.Printf("[EXTENSION] failed to unmarshal: %v", err)
		return nil, err
	}

	return secretBinary, nil
}
