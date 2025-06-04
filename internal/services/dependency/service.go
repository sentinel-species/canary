package dependency

import (
	"canary/internal/logger"
	"canary/internal/store"
	"canary/internal/types"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"log"
	"strings"
)

type Service interface {
	Resolve(ctx context.Context, req *types.DependencyCheckRequest) (types.DependencyCheckResponse, error)
}

type service struct {
	store store.DependencyStore
}

func NewService(store store.DependencyStore) Service {
	return &service{
		store: store,
	}
}

func (s *service) Resolve(ctx context.Context, req *types.DependencyCheckRequest) (types.DependencyCheckResponse, error) {
	payload := types.DependencyCheckRequest{
		Ecosystem:      req.Ecosystem,
		Package:        req.Package,
		CurrentVersion: req.CurrentVersion,
		UpgradeVersion: req.UpgradeVersion,
		FileContent:    req.FileContent,
	}
	rawUpgradeVersions, err := s.store.ResolveUpgrades(ctx, payload)
	if err != nil {
		logger.Error("Failed: ", err)
	}

	// TODO This chose something off
	upgradeVersions := filterLowestFixVersions(rawUpgradeVersions)
	decodedFileContent, err := base64.StdEncoding.DecodeString(req.FileContent)
	if err != nil {
		logger.Error("Failed: ", err)
	}

	upgradedRequirementsFile := upgradeRequirements("pypi", string(decodedFileContent), upgradeVersions)
	response := types.DependencyCheckResponse{
		FileContent: *upgradedRequirementsFile,
	}
	// TODO Response failed for reason
	return response, nil
}

func upgradeRequirements(environment string, requirementsFile string, upgrades []types.UpgradeResponse) *string {
	var lines []string
	// Build a map for quick lookup
	upgradeMap := make(map[string]string)
	for _, u := range upgrades {
		upgradeMap[strings.ToLower(u.PackageName)] = u.FixVersion
	}
	switch environment {
	case "pypi":
		lines = updateRequirementsTxt(upgradeMap, requirementsFile)
	default:
		return nil
	}
	upgradedRequirementsFile := strings.Join(lines, "\n")
	return &upgradedRequirementsFile
}

func updateRequirementsTxt(upgradeMap map[string]string, requirementsFile string) []string {
	lines := strings.Split(requirementsFile, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Ignore nonsense in requirements.txt
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}

		parts := strings.SplitN(trimmed, "==", 2)
		if len(parts) != 2 {
			continue
		}

		pkgName := strings.ToLower(strings.TrimSpace(parts[0]))
		if newVer, ok := upgradeMap[pkgName]; ok {
			lines[i] = fmt.Sprintf("%s==%s", parts[0], newVer)
		}
	}
	return lines
}

func filterLowestFixVersions(upgrades []types.UpgradeResponse) []types.UpgradeResponse {
	filtered := make(map[string]types.UpgradeResponse)

	for _, u := range upgrades {
		key := u.PackageName

		v, err := semver.NewVersion(u.FixVersion)
		if err != nil {
			log.Printf("Skipping invalid version %q for package %q: %v", u.FixVersion, u.PackageName, err)
			continue
		}

		existing, exists := filtered[key]
		if !exists {
			filtered[key] = u
			continue
		}

		ev, err := semver.NewVersion(existing.FixVersion)
		if err != nil || v.LessThan(ev) {
			filtered[key] = u
		}
	}

	result := make([]types.UpgradeResponse, 0, len(filtered))
	for _, u := range filtered {
		result = append(result, u)
	}

	return result
}
