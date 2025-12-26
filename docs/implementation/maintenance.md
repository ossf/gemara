---
layout: page
title: Release & Maintenance
---

Project release deliverables are divided into the **Core Specification** and language-specific **SDKs**.

* **Core Specification Release:** Bundles the Model, Lexicon, and CUE Schemas together. These are versioned and released as a single unit because they are tightly coupled.
* **SDK Releases:** Language-specific implementations that provide tooling, types, and helpers to work with Gemara documents. SDK types are generated from the CUE schemas.

Each **maintains its own** independent [SemVer](https://semver.org/) lifecycle.

## Specification Release Versioning

The core specification release bundles the Model, Lexicon, and CUE Schemas together and versions them as a single unit.

| Change Type | Version Bump | Examples                                                     |
|:------------|:-------------|:-------------------------------------------------------------|
| Major       | v2.0.0       | Breaking changes to the Model, Lexicon, or Stable schemas.   |
| Minor       | v1.(x+1).0   | Additive changes, schema promotions, or new optional fields. |
| Patch       | v1.x.(y+1)   | Bug fixes in schema logic or documentation.                  |

## Schema Lifecycle and Major Version Example

Possible schema states include: **Experimental** → **Stable** → **Deprecated**. These are denoted on each layer schema with a `@status(experimental|stable|deprecated)` attribute.

The following table illustrates how schemas progress through their lifecycle and how major version changes are handled:

| Version | Status              | Change Type | Example Scenario                        |
|:--------|:--------------------|:------------|:----------------------------------------|
| v1.0.0  | Experimental        | Initial     | Schema first published                  |
| v1.1.0  | Experimental        | Minor       | Optional fields added                   |
| v1.2.0  | Stable              | Minor       | Promoted to Stable                      |
| v1.3.0  | Stable              | Minor       | Additive changes                        |
| v1.3.1  | Stable              | Patch       | Bug fix                                 |
| v1.4.0  | Stable              | Minor       | Field deprecated, replacement added     |
| v1.5.0  | Stable → Deprecated | Minor       | Original deprecated, replacement Stable |
| v2.0.0  | Stable              | Major       | Deprecated schema removed               |
| v2.1.0  | Stable              | Minor       | Additive changes                        |

### Experimental Status
* All new schemas start as Experimental.
* Adding Experimental schemas or making breaking changes triggers minor version increments.
* Breaking changes and performance issues may occur during this phase.
* These schemas may not be feature-complete.

### Stable Status

* Layers promote independently. Each layer only requires its direct dependencies to be Stable (e.g., Layer 2 requires Layer 1, but not Layer 6).
* Layers can be promoted to Stable at different times. Layer 2 can be Stable while Layer 6 remains Experimental.
* Stable schemas may only reference other Stable schemas.
* Stable schemas maintain backward and forward compatibility within major versions, allowing **additive optional changes** only.
* Breaking changes require major version increments and should be avoided in all normal circumstances.

### Deprecated Status

* Schemas or fields within schemas may be deprecated when replaced.
* Replacement schema or field must be added in Experimental status and promoted to Stable before deprecation.
* Deprecated schemas and fields maintain the same support guarantees as Stable schemas and remain functional.
* Deprecated schemas are removed in the next major version release.

## SDK Release Versioning

SDK releases are versioned independently to allow for rapid iteration on tooling, bug fixes, and utilities without requiring core specification version bumps.

* The CUE schemas in the core specification release are the source of truth for validation.
* SDKs must explicitly document which core specification release version they support.
* When a new core specification release is published, SDKs regenerate their types from the updated schemas and release new versions that support the updated specification.

## Questions or Feedback

For questions about versioning strategy or to propose changes:

* Open an issue on [GitHub](https://github.com/ossf/gemara/issues)
* Discuss in the [OpenSSF Slack #gemara channel](https://openssf.slack.com/archives/C09A9PP765Q)
* Attend the [biweekly Gemara meeting](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
