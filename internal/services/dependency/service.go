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
	// Format file content
	decoded, err := base64.StdEncoding.DecodeString(req.FileContent)
	decodedFileContent := string(decoded)
	if err != nil {
		logger.Error("Failed: ", err)
	}
	req.FileContent = decodedFileContent
	currentPackages := extractPackageNames(req.FileContent)

	rawUpgradeVersions, err := s.store.ResolveUpgrades(ctx, *req, currentPackages)
	if err != nil {
		logger.Error("Failed: ", err)
	}

	minimalUpgradableVersions := filterLowestFixVersions(rawUpgradeVersions)

	upgradedRequirementsFile := upgradeRequirements(minimalUpgradableVersions, *req)
	response := types.DependencyCheckResponse{
		FileContent: *upgradedRequirementsFile,
	}
	return response, nil
}

func upgradeRequirements(upgrades []types.UpgradeResponse, payload types.DependencyCheckRequest) *string {
	var lines []string
	// Build a map for quick lookup
	upgradeMap := make(map[string]string)
	for _, u := range upgrades {
		upgradeMap[strings.ToLower(u.PackageName)] = u.FixVersion
	}
	switch payload.Ecosystem {
	case "pypi":
		lines = updatePypiFileContent(upgradeMap, payload)
	default:
		return nil
	}
	// Make it look exactly like the existing file
	encodedFileContent := base64.StdEncoding.EncodeToString([]byte(strings.Join(lines, "\n")))
	return &encodedFileContent
}

func updatePypiFileContent(upgradeMap map[string]string, payload types.DependencyCheckRequest) []string {
	lines := strings.Split(payload.FileContent, "\n")

	// Create map of existingRequirements
	existingRequirements := make(map[string]string)
	for _, er := range lines {
		// Format
		trimmed := strings.TrimSpace(er)
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}

		parts := strings.SplitN(trimmed, "==", 2)
		if len(parts) != 2 {
			continue
		}
		pkgName := strings.ToLower(strings.TrimSpace(parts[0]))
		pkgVersion := strings.ToLower(strings.TrimSpace(parts[1]))
		existingRequirements[pkgName] = pkgVersion
	}

	// Check if upgrade package exist in root
	if _, ok := existingRequirements[payload.Package]; ok {
		existingRequirements[payload.Package] = payload.UpgradeVersion
	} else { // Apply all upgrades
		for pkgName, pkgUpgradeVersion := range upgradeMap {
			if vulnerableVersion, ok := existingRequirements[pkgName]; ok {
				existingVersion, _ := semver.NewVersion(existingRequirements[pkgName])
				upgradeVersion, _ := semver.NewVersion(pkgUpgradeVersion)
				if !(existingVersion.LessThan(upgradeVersion)) {
					continue
				}
				logger.Debug(fmt.Sprintf("Found vulnerable dependency %s for package %s, upgrading to %s", vulnerableVersion, pkgName, pkgUpgradeVersion))
				existingRequirements[pkgName] = pkgUpgradeVersion
			}
		}
	}

	var result []string
	for key, value := range existingRequirements {
		result = append(result, fmt.Sprintf("%s==%s", key, value))
	}
	return result
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

func extractPackageNames(input string) []string {
	lines := strings.Split(input, "\n")
	var packages []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}

		parts := strings.SplitN(trimmed, "==", 2)
		if len(parts) == 2 {
			pkg := strings.TrimSpace(parts[0])
			packages = append(packages, strings.ToLower(pkg))
		}
	}

	return packages
}
