---
layout: page
title: The Gemara Model
---

**Status**: <span class="badge badge-stable">Stable</span>

The Gemara Model describes the six categorical layers of GRC (Governance, Risk, Compliance) 
activities. These layers represent the longstanding reality of how GRC activities are 
organized and interact.

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
    <p>Prevention or remediation based on assessment findings.</p>
  </div>

  <div>
    <h3>Layer 6: Audit</h3>
    <p>Review of organizational policy and conformance</em>.</p>
  </div>
</div>

## Model Stability

This model is intentionally stable. Changes to the model are rare and require significant 
community discussion, as the model reflects fundamental organizational patterns in GRC engineering 
activities.

**Why Stability Matters:**
- Provides a consistent foundation for all Gemara work
- Allows the Lexicon and Implementation to evolve without model changes
- Ensures long-term compatibility and understanding

If you're interested in proposing changes to the model, please consider opening an Architecture Decision 
Record (ADR). See the **[ADR page](../adr.html)** for more information.

## Relationship to Other Components

As one of Gemara's three core components, the Model works alongside:

### [The Lexicon Component](../lexicon)
The Lexiconcomponent provides specific definitions for terms used 
within each layer. While the Model describes the structure, the Lexicon provides the 
shared vocabulary that teams can agree on.

### [The Implementation Component](../implementation)
The Implementation component provides schemas and SDK libraries 
based on the Model. The Model describes the conceptual layers, while the Implementation 
provides machine-readable formats and programming APIs.

Together, these three components support the entire Gemara ecosystem.

