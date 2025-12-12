---
layout: home
title: Home
---

# Gemara <span class="pronunciation">(Juh-MAH-ruh)</span>

<img src="{{ '/assets/gemara-logo.png' | relative_url }}" alt="Gemara Logo" class="gemara-logo" />

**GRC Engineering Model for Automated Risk Assessment**

Gemara provides a logical model to describe the categories of compliance activities, how they interact, and the schemas to enable automated interoperability between them.

In order to better facilitate cross-functional communication, the Gemara Model seeks to outline the categorical layers of activities related to automated governance.

<!--
## Quick Start

- **New to Gemara?** Start with our [About page](/about) to understand the model
- **Want to dive deeper?** Explore the [Six Layers](/layers) of the model
- **Ready to build?** Check out our [Tutorial](/tutorial) for a hands-on example
- **Want to contribute?** See our [Contributing Guide](/contributing)
-->

## The Six Layers

Gemara organizes compliance activities into six categorical layers, each building upon the previous:

<div class="layer-grid">
  <div>
    <h3>Layer 1: Guidance</h3>
    <p>High-level guidance on cybersecurity measures from industry groups and standards bodies.</p>
  </div>

  <div>
    <h3>Layer 2: Controls</h3>
    <p>Technology-specific, threat-informed security controls for protecting information systems.</p>
  </div>

  <div>
    <h3>Layer 3: Policy</h3>
    <p>Risk-informed guidance tailored to your organization's specific needs and risk appetite.</p>
  </div>

  <div>
    <h3>Layer 4: Evaluation</h3>
    <p>Inspection of code, configurations, and deployments against policies and controls.</p>
  </div>

  <div>
    <h3>Layer 5: Enforcement</h3>
    <p>Prevention or remediation based on assessment findings <em>(Coming Soon)</em>.</p>
  </div>

  <div>
    <h3>Layer 6: Audit</h3>
    <p>Review of organizational policy and conformance <em>(Coming Soon)</em>.</p>
  </div>
</div>

## Real-World Usage

Gemara is being used today in production environments:

- **[FINOS Common Cloud Controls](https://www.finos.org/common-cloud-controls-project)** - Layer 2 controls for cloud environments
- **[Open Source Project Security Baseline](https://baseline.openssf.org/)** - Layer 2 security baseline for open source projects
- **[Privateer](https://github.com/privateerproj/privateer)** - Layer 4 evaluation framework with plugins like the [OSPS Baseline Plugin](https://github.com/revanite-io/pvtr-github-repo)

## Get Started

Use the CUE schemas directly for content validation.

Refer to the official CUE documentation for [installation instructions](https://cuelang.org/docs/introduction/installation/).

```bash
cue vet ./your-data.yaml ./schemas/layer-2.cue
```

If you are building automated tools, we maintain a Go module to help you manipulate Gemara data.

The consolidated details and documentation can be quickly navigated on the Go registry.

[![Go Reference](https://pkg.go.dev/badge/github.com/ossf/gemara.svg)](https://pkg.go.dev/github.com/ossf/gemara)

```bash
go get github.com/ossf/gemara
```

## Community

Join the conversation:

- **Slack:** [#gemara](https://openssf.slack.com/archives/C09A9PP765Q) on OpenSSF Slack
- **Meetings:** Bi-weekly on alternate Thursdays - see the [OpenSSF calendar](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
- **GitHub:** [ossf/gemara](https://github.com/ossf/gemara)