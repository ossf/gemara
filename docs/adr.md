---
layout: page
nav-title: ADR
---

# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for the Gemara project.

## What are ADRs?

Architecture Decision Records document important architectural decisions made in the project. They capture:
- The context that led to the decision
- The decision itself
- The consequences (positive, negative, and neutral)
- Alternatives that were considered

## ADR Index

- [ADR-0001: Unified Package Structure](./adr/0001-unified-package-structure.html) - Decision to consolidate layer-based packages into a single unified `gemara` package

## Format

ADRs follow the format described at [adr.github.io](https://adr.github.io/), with the following structure:

1. **Status** - Proposed, Accepted, Deprecated, or Superseded
2. **Context** - The situation and problem that led to this decision
3. **Decision** - The architectural decision and its rationale
4. **Consequences** - Positive, negative, and neutral outcomes
5. **Alternatives Considered** - Other options that were evaluated

## When to Create an ADR

Create an ADR when:
- Making a significant architectural decision that affects the public API
- Choosing between multiple viable approaches
- The decision will impact future development or maintenance
- The decision needs to be communicated to stakeholders

Don't create an ADR for:
- Routine implementation details
- Temporary workarounds
- Decisions that are clearly the only viable option

