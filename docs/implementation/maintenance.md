---
layout: page
title: Maintenance
---

## Versioning Strategy

Gemara follows [Semantic Versioning](https://semver.org/) (SemVer) with a **single release cycle** for both the Go module and CUE schemas.
A single version tag (`v1.2.3`) applies to both the Go module (`github.com/ossf/gemara`) and the CUE schemas (`github.com/ossf/gemara/schemas`).

Schema changes almost always affect the Go library because Go types are generated from the CUE schemas.
The unified release cycle ensures consistency across both components.
All schemas are versioned together as a single bundle (schema set).

| Change Type | Version Bump | Examples                                                                                          |
|:------------|:-------------|:--------------------------------------------------------------------------------------------------|
| Major       | v2.0.0       | Breaking changes violating backward/forward compatibility                                         |
| Minor       | v1.x+1.0     | Additive changes, schema promotions, schema deprecations, field deprecations, new optional fields |
| Patch       | v1.x.y+1     | Bug fixes                                                                                         |

## Schema Lifecycle

Possible schema states include: **Experimental** → **Stable** → **Deprecated**.

### Experimental Status

* Schemas start Experimental.
* Adding Experimental schemas triggers minor increments.
* Breaking changes and performance issues may occur.
* These schemas may not be feature-complete.

### Stable Status

* Promoting to Stable triggers minor increments and requires release notes and documentation updates.
* Layers promote independently. Each layer only requires its direct dependencies to be Stable (e.g., Layer 2 requires Layer 1, but not Layer 6).
* Layers can be promoted to Stable at different times. Layer 2 can be Stable while Layer 6 remains Experimental.
* For v2.0.0 release, at least some schemas must be Stable, respecting dependency order.
* Stable schemas may only reference other Stable schemas.
* Stable schemas maintain backward and forward compatibility within major versions, allowing additive changes.
* Forward compatibility means new required fields are not added, ensuring older tooling continues to work with new data.
* Breaking changes require major version increments and should be avoided in all normal circumstances.

### Deprecated Status

* Schemas or fields within schemas may be deprecated when replaced.
* Schemas or fields MUST NOT be marked as deprecated unless the replacement is Stable.
* The replacement schema or field must be added in Experimental status and promoted to Stable before deprecation.
* Deprecating a schema or field triggers a minor version increment.
* Deprecated schemas and fields maintain the same support guarantees as Stable schemas and remain functional.
* Go types generated from deprecated schemas remain available in the same major version's Go module.
* Deprecated fields remain available in the same major version's Go types.
* When v2.0.0 is released, deprecated v1 schemas and fields are excluded from v2 but remain available in v1.x releases.

## Versioning for Go

* Release Branching used for major version changes (v2.0.0+). 
* Go supports multiple major versions via `/v2` subdirectories or branches.

**Breaking changes trigger this workflow**:

| Step | Action                               | Branch State              |
|:-----|:-------------------------------------|:--------------------------|
| 1    | Identify schemas to be replaced in v2 | Main branch               |
| 2a   | Create v1 maintenance branch         | Isolated                  |
| 2b   | Continue development on main         | v2 Experimental           |
| 3    | Promote v2 schemas to Stable         | Respect dependency order  |
| 4    | Deprecate corresponding v1 schemas   | When v2 replacements Stable |
| 5    | Tag v2.0.0 release                   | At least some v2 Stable   |

* v1 remains on the maintenance branch for bug fixes.
* New features occur only in v2.
* As v2 schemas stabilize, corresponding v1 schemas are deprecated.
* Deprecated v1 schemas and their Go types remain available in v1.x releases but are excluded from v2.0.0+.

## Questions or Feedback

For questions about versioning strategy or to propose changes:

* Open an issue on [GitHub](https://github.com/ossf/gemara/issues)
* Discuss in the [OpenSSF Slack #gemara channel](https://openssf.slack.com/archives/C09A9PP765Q)
* Attend the [biweekly Gemara meeting](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
