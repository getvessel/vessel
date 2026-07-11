# Release Guide

## Versioning

Vessel follows [Semantic Versioning](https://semver.org/) (`vMAJOR.MINOR.PATCH`). Until `v1.0.0`, breaking changes may occur in minor versions.

Current version: `v0.1.0-alpha`.

## Release Process

1. **Branch**: Develop on `next`. Open a PR into `main`.
2. **Version bump**: Update `const vesselVersion` in `cmd/main.go`.
3. **Changelog**: Summarise changes in the GitHub release body.
4. **Tag & release**:

   ```bash
   git tag -a v0.1.1 -m "v0.1.1"
   git push origin v0.1.1
   ```

5. **CI**: GitHub Actions builds the binary and Docker image, pushes to GHCR, and creates a release draft.
6. **Publish**: Review the draft and publish.

## Upgrade

```bash
curl -fsSL https://get.vessl.dev | sh
```

Or in the dashboard: `Settings → Updates → Check for Updates`.

## Downgrade

Downgrading is supported but not recommended (DB schema migrations may not be backward-compatible).

```bash
# Restore previous binary + data snapshot
sudo vesselctl rollback
```

A manual downgrade requires restoring the previous binary and a database backup taken before the upgrade.
