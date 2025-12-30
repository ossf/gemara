---
layout: page
title: The Gemara Model
---

**Status**: <span class="badge badge-stable">Stable</span>

The Gemara Model describes six categorical layers of GRC (Governance, Risk, Compliance) activities, representing how GRC activities are organized and interact.

## The Six Layers

Gemara organizes compliance activities into six categorical layers, each building upon the previous:

<div class="gemara-layer-diagram">
  <div class="layer-banner layer-6">
    <span class="layer-number">6</span>
    <div class="layer-content">
      <div class="layer-title">Audit</div>
      <div class="layer-description">Quality & Efficacy Review of all GRC Outputs</div>
    </div>
  </div>
  <div class="layer-banner layer-5">
    <span class="layer-number">5</span>
    <div class="layer-content">
      <div class="layer-title">Enforcement</div>
      <div class="layer-description">Remediation or Deployment Prevention</div>
    </div>
  </div>
  <div class="layer-banner layer-4">
    <span class="layer-number">4</span>
    <div class="layer-content">
      <div class="layer-title">Evaluation</div>
      <div class="layer-description">Inspection of Sensitive Activity Results</div>
    </div>
  </div>
  <div class="layer-banner layer-sensitive">
    <div class="layer-content">
      <div class="layer-title">Sensitive Activities</div>
      <div class="layer-description">e.g. Infrastructure & Application Development</div>
    </div>
  </div>
  <div class="layer-banner layer-3">
    <span class="layer-number">3</span>
    <div class="layer-content">
      <div class="layer-title">Policy</div>
      <div class="layer-description">Organizational-specific; Risk-informed</div>
    </div>
  </div>
  <div class="layer-banner layer-2">
    <span class="layer-number">2</span>
    <div class="layer-content">
      <div class="layer-title">Objectives</div>
      <div class="layer-description">Technology-specific; Threat-informed</div>
    </div>
  </div>
  <div class="layer-banner layer-1">
    <span class="layer-number">1</span>
    <div class="layer-content">
      <div class="layer-title">Guidance</div>
      <div class="layer-description">High-level Goals, Regulations, or Best Practices</div>
    </div>
  </div>
</div>

## Model Stability

This model is intentionally stable. Changes are rare and require significant community discussion, as the model reflects fundamental organizational patterns in GRC activities.

**Why Stability Matters:**
- Provides a consistent foundation for all Gemara work
- Allows the Lexicon and Implementation to evolve independently
- Ensures long-term compatibility

## Relationship to Other Components

### [The Lexicon](../lexicon)
Provides definitions for terms used within each layer. The Model describes structure; the Lexicon provides shared vocabulary.

### [The Implementation](../implementation)
Provides schemas and SDKs based on the Model. The Model describes conceptual layers; the Implementation provides machine-readable formats and APIs.

