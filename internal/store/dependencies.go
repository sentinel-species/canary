package store

import (
	db "canary/internal/database"
	"canary/internal/types"
	"context"
)

type DependencyStore struct {
	store db.Pool
}

func NewDependencyStore(s db.Pool) *DependencyStore {
	return &DependencyStore{store: s}
}

func (s *DependencyStore) ResolveUpgrades(ctx context.Context, payload types.DependencyCheckRequest, currentPackages []string) (upgrades []types.UpgradeResponse, err error) {
	rows, err := s.store.Query(ctx, getTransientPackageUpgradesQuery, payload.Package, payload.UpgradeVersion, payload.Ecosystem, currentPackages)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var upgradedPackage types.UpgradeResponse
		if err := rows.Scan(&upgradedPackage.Id, &upgradedPackage.PackageName, &upgradedPackage.FixVersion, &upgradedPackage.EnvName); err != nil {
			return nil, err
		}
		upgrades = append(upgrades, upgradedPackage)
	}
	return
}

// Given a fix version of a package, find all packages that exist in your requirements file return the required versions of these packages
// Return the versions of the required packages
const getTransientPackageUpgradesQuery = `
	SELECT p.id AS id, p.name AS package_name, p.version AS fix_version, e.name AS environment_name
	FROM dependencies d
	JOIN packages p ON p.internal_id = d.package_id
	JOIN environments e ON p.environment_id = e.internal_id
	JOIN packages dep ON d.dependency_id = dep.internal_id
	WHERE
		-- Query the version it needs to become
		dep.name =  $1-- 'werkzeug' -- VULNERABLE PACKAGE
		AND dep.version = $2 -- '3.1.3' -- FIX VERSION
		AND e.name = $3 -- 'pypi' -- PACKAGE ECOSYSTEM
		AND dep.environment_id = e.internal_id -- DEPENDENCY ECOSYSTEM
		AND p.name = ANY ($4)-- ('flask') -- These are the packages in my requirements.txt
`
