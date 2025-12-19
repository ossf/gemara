---
layout: page
title: Maintenance
---

This document outlines the versioning and maintenance strategy for Gemara's implementation components.

## Versioning Strategy

Gemara follows [Semantic Versioning](https://semver.org/) (SemVer) with a **single release cycle** for both the Go module and CUE schemas. A single version tag (`v1.2.3`) applies to both the Go module (`github.com/ossf/gemara`) and the CUE schemas (`github.com/ossf/gemara/schemas`).

Schema changes almost always affect the Go library because Go types are generated from the CUE schemas. While Go library changes don't always require schema changes, using a unified release cycle simplifies versioning and ensures consistency.

Version increments follow SemVer:
* **Major version** increments for breaking changes (schema or API)
* **Minor version** increments for additive changes (new features, schema status promotions, or new fields/types in Stable schemas)
* **Patch version** increments for bug fixes

### Backward Compatible Changes to Stable Schemas

Backward compatible changes to Stable schemas trigger **minor version increments**.

Like:
* Adding new optional fields to existing types
* Adding new types
* Adding new optional properties

## Release Process

Changes are reviewed and tested before release. Git tags follow SemVer format (`v1.2.3`), and release notes document changes and migration paths.

## Schema Promotion

Schemas follow this lifecycle: **Experimental** → **Stable** → **Deprecated**. Status promotions trigger a **minor version increment**, as they are additive and do not break backward compatibility.

### Experimental Status

* Schemas start in **Experimental** status.
* Adding a new Experimental schema triggers a **minor version increment**.
* Breaking changes and performance issues MAY occur.
* Components **SHOULD NOT** be expected to be feature-complete.

### Stable Status

* Promoting a schema from Experimental to Stable triggers a **minor version increment** and involves a stabilization announcement in release notes, documentation updates, and tracking schema maturity.
* Individual layers can be promoted independently (e.g., Layer 2 can be stable while Layer 1 remains in Experimental).
* Once Stable, schemas **can still evolve** but maintain backward compatibility within major versions. Stable schemas allow **additive changes** such as new optional fields or new types.
<<<<<<< Updated upstream
* Breaking changes to Stable schemas require a major version increment. This should be avoided in all normal circumstances.
* Stable schemas represent a long-term commitment and will continue to be supported.
=======
* Breaking changes to Stable schemas require a major version increment.

#### Stable Schema Support

Stable schemas represent a long-term commitment and will continue to be supported until explicitly deprecated.

**Support includes:**

| Support Type           | Description                                   |
|:-----------------------|:----------------------------------------------|
| Backward compatibility | Maintain compatibility within major version.  |
| Bug fixes              | Fix bugs and issues without breaking changes. |
| Documentation          | Keep documentation current and accurate.      |
| Migration guidance     | Provide migration paths when deprecating.     |

Deprecation follows a clear process with migration paths.

### Deprecated Status

* When a Stable schema needs to be replaced, the deprecation process follows these steps:
  1. Add the replacing schema in **Experimental** status (triggers a **minor version increment**).
  2. Promote the replacing schema from Experimental to **Stable** (triggers a **minor version increment**).
  3. Remove the old deprecated schema (triggers a **major version increment**).
* Deprecated schemas are marked for removal and users are notified through release notes and documentation.
* The deprecation period provides time for users to migrate to the replacement schema before removal.

## Compatibility

* The Go library maintains backward compatibility within major versions.
* Schemas maintain backward compatibility within major versions (additive changes only).

## Questions or Feedback

For questions about versioning strategy or to propose changes:

* Open an issue on [GitHub](https://github.com/ossf/gemara/issues)
* Discuss in the [OpenSSF Slack #gemara channel](https://openssf.slack.com/archives/C09A9PP765Q)
* Attend the [biweekly Gemara meeting](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
