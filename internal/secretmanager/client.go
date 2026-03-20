package secretmanager

import (
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Client wraps the GCP Secret Manager API client (uses Application Default Credentials).
type Client struct {
	inner *secretmanager.Client
}

// NewClient creates a client. Close it when the manager stops if you need clean shutdown.
func NewClient(ctx context.Context) (*Client, error) {
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("secretmanager.NewClient: %w", err)
	}
	return &Client{inner: c}, nil
}

// Close releases client resources.
func (c *Client) Close() error {
	if c == nil || c.inner == nil {
		return nil
	}
	return c.inner.Close()
}

// GetSecretData fetches a secret version and maps it to string key/value data for Kubernetes Secret.Data.
func (c *Client) GetSecretData(ctx context.Context, projectID, secretID, version, decodeFormat, secretKey string) (map[string]string, error) {
	if version == "" {
		version = "latest"
	}
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, secretID, version)
	req := &secretmanagerpb.AccessSecretVersionRequest{Name: name}
	resp, err := c.inner.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AccessSecretVersion %s: %w", name, err)
	}
	raw := string(resp.Payload.GetData())

	switch decodeFormat {
	case "json":
		return decodeJSONPayload(raw)
	default:
		key := secretKey
		if key == "" {
			key = "value"
		}
		return map[string]string{key: raw}, nil
	}
}

func decodeJSONPayload(raw string) (map[string]string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return nil, fmt.Errorf("decode json payload: %w", err)
	}
	out := make(map[string]string, len(obj))
	for k, v := range obj {
		switch t := v.(type) {
		case string:
			out[k] = t
		case float64:
			out[k] = fmt.Sprintf("%g", t)
		case bool:
			out[k] = fmt.Sprintf("%v", t)
		case nil:
			out[k] = ""
		default:
			b, err := json.Marshal(t)
			if err != nil {
				return nil, fmt.Errorf("key %q: %w", k, err)
			}
			out[k] = string(b)
		}
	}
	return out, nil
}
