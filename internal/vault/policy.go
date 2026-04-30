package vault

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// PolicyEntry represents a single secret path and its allowed capabilities.
type PolicyEntry struct {
	Path         string   `json:"path"`
	Capabilities []string `json:"capabilities"`
}

// PolicyReport holds the result of auditing a secret path against a Vault policy.
type PolicyReport struct {
	Path    string        `json:"path"`
	Entries []PolicyEntry `json:"entries"`
}

// PolicyChecker checks whether a given secret path is covered by Vault policies.
type PolicyChecker struct {
	client *Client
}

// NewPolicyChecker creates a new PolicyChecker using the provided Vault client.
func NewPolicyChecker(c *Client) *PolicyChecker {
	return &PolicyChecker{client: c}
}

// CheckPath returns a PolicyReport for the given secret path by querying
// all policies that reference it. mount defaults to "secret" if empty.
func (pc *PolicyChecker) CheckPath(ctx context.Context, secretPath, mount string) (*PolicyReport, error) {
	if mount == "" {
		mount = "secret"
	}

	fullPath := fmt.Sprintf("%s/data/%s", mount, strings.TrimPrefix(secretPath, "/"))

	policies, err := pc.client.Logical().ListWithContext(ctx, "sys/policy")
	if err != nil {
		return nil, fmt.Errorf("listing policies: %w", err)
	}

	report := &PolicyReport{Path: fullPath}

	if policies == nil || policies.Data == nil {
		return report, nil
	}

	keys, _ := policies.Data["keys"].([]interface{})
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
	})

	for _, k := range keys {
		name := fmt.Sprint(k)
		entry, err := pc.resolvePolicy(ctx, name, fullPath)
		if err != nil {
			continue
		}
		if entry != nil {
			report.Entries = append(report.Entries, *entry)
		}
	}

	return report, nil
}

func (pc *PolicyChecker) resolvePolicy(ctx context.Context, name, fullPath string) (*PolicyEntry, error) {
	secret, err := pc.client.Logical().ReadWithContext(ctx, fmt.Sprintf("sys/policy/%s", name))
	if err != nil || secret == nil {
		return nil, err
	}

	rules, _ := secret.Data["rules"].(string)
	if !strings.Contains(rules, fullPath) {
		return nil, nil
	}

	return &PolicyEntry{
		Path:         name,
		Capabilities: extractCapabilities(rules, fullPath),
	}, nil
}

func extractCapabilities(rules, path string) []string {
	lines := strings.Split(rules, "\n")
	for i, line := range lines {
		if strings.Contains(line, path) && i+1 < len(lines) {
			caps := lines[i+1]
			caps = strings.TrimSpace(caps)
			caps = strings.Trim(caps, "capabilities = []")
			var result []string
			for _, c := range strings.Split(caps, ",") {
				c = strings.Trim(strings.TrimSpace(c), `"`)
				if c != "" {
					result = append(result, c)
				}
			}
			return result
		}
	}
	return nil
}
