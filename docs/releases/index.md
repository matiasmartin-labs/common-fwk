---
title: Releases
nav_order: 4
has_children: true
permalink: /releases/
---

# Releases

This section documents every published release of `common-fwk`, reconstructed from
tag-to-tag commit comparison.

| Version | Date | Summary |
|---|---|---|
| [v0.9.0](v0.9.0/) | TBD | RSA private key accessor |
| [v0.8.0](v0.8.0/) | 2026-05-02 | Default JSON 404/405 handlers |
| [v0.7.0](v0.7.0/) | 2026-05-01 | slog logger registry with scoped controls |
| [v0.6.0](v0.6.0/) | 2026-05-01 | Opt-in health/readiness endpoint presets |
| [v0.5.0](v0.5.0/) | 2026-05-01 | Read-only app runtime accessors |
| [v0.4.0](v0.4.0/) | 2026-05-01 | RS256 keypair security mode |
| [v0.3.0](v0.3.0/) | 2026-05-01 | HTTP server runtime limits |
| [v0.2.0](v0.2.0/) | 2026-05-01 | Kebab-case config keys (Viper adapter) |
| [v0.1.0](v0.1.0/) | 2026-04-25 | Initial stable release |

## Release Labels

Release automation uses the following labels on PRs:

| Label | Effect |
|---|---|
| `release-type:patch` | Patch version bump (bug fixes) |
| `release-type:minor` | Minor version bump (new features) |
| `release-type:major` | Major version bump (breaking changes) |
| `release:skip` | Skip release generation |
