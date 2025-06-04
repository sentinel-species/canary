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

func (s *DependencyStore) ResolveUpgrades(ctx context.Context, payload types.DependencyCheckRequest) (upgrades []types.UpgradeResponse, err error) {
	// TODO decode filecontent
	rows, err := s.store.Query(ctx, getTransientPackageUpgradesQuery, payload.Package, payload.UpgradeVersion, payload.Ecosystem, "flask")
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

//func (s *DependencyStore) InsertPackage(ctx context.Context, pkg types.Package, envID int) (int, error) {
//	var id int
//	// Try insert, fallback to select if already exists
//	err := s.store.QueryRow(ctx, `
//        INSERT INTO packages (name, version, environment_id)
//        VALUES ($1, $2, $3)
//        ON CONFLICT (name, version, environment_id) DO NOTHING
//        RETURNING internal_id
//    `, pkg.Name, pkg.Version, envID).Scan(&id)
//
//	if errors.Is(err, sql.ErrNoRows) {
//		// Already exists, fetch the id
//		err = s.store.QueryRow(ctx, `
//            SELECT internal_id FROM packages
//            WHERE name = $1 AND version = $2 AND environment_id = $3
//        `, pkg.Name, pkg.Version, envID).Scan(&id)
//	}
//
//	if err != nil {
//		return 0, err
//	}
//
//	// Recursively insert dependencies
//	for _, dep := range pkg.Dependencies {
//		depID, err := s.InsertPackage(ctx, dep, envID)
//		if err != nil {
//			return 0, err
//		}
//
//		// Insert dependency relation
//		_, err = s.store.Exec(ctx, `
//            INSERT INTO dependencies (package_id, dependency_id)
//            VALUES ($1, $2)
//            ON CONFLICT DO NOTHING
//        `, id, depID)
//
//		if err != nil {
//			return 0, err
//		}
//	}
//
//	return id, nil
//}

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
		AND p.name IN ($4)-- ('flask') -- These are the packages in my requirements.txt
`
