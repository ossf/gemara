---
layout: default
title: Creating Guidance Documents
---

# Creating and Extending Guidance Documents

This tutorial demonstrates how to create Layer 1 Guidance documents and extend existing guidelines.

## Use Case: Extending NIST Controls

An organization needs to extend specific NIST 800-53 controls with domain-specific requirements for cloud-based access management.
This demonstrates how to create extensions that can be reused in multiple policies.

### Step 1: Create Extension Guidance Document

Create a new Guidance document that extends NIST controls:

```yaml
title: Guidance for Cloud Access Management
metadata:
  id: Cloud-Access-Guidance
  description: Domain-specific extensions to NIST 800-53 controls for cloud-based access management
  author:
    id: owner
    name: Owner
    type: Human
  version: "v0.1.0"
  date: "2025-12-29"
  applicability-categories:
    - id: cloud-services
      title: Cloud Services
      description: Applies to cloud-hosted applications and services
  mapping-references:
    - id: NIST-800-53-Rev5
      title: NIST Special Publication 800-53 Revision 5
      version: "Rev 5"
      description: Complete NIST 800-53 Catalog
      url: "https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final"

families:
  - id: Access-Control
    title: Access Control
    description: Controls for managing access to information systems

guidelines:
  # Add cloud-specific statements that supplement the base guideline for the domain
  - id: ORG-AC-3-EXT
    title: Cloud Access Enforcement Extension
    family: Access-Control
    extends:
      reference-id: NIST-800-53-Rev5
      entry-id: AC-3
    applicability: 
      - cloud-services
    statements:
      - id: AC-3-EXT.1
        title: Multi-Factor Authentication Requirement
        text: All cloud service access must require multi-factor authentication using organization-approved authenticators.
  # Independent controls don't use `extends`
  - id: ORG-AC-CLOUD-01
    title: Cloud Session Management
    family: Access-Control
    objective: Ensure cloud sessions are properly managed and terminated
    statements:
      - id: ORG-AC-CLOUD-01.1
        title: Session Timeout Requirement
        text: Cloud service sessions must automatically timeout after 30 minutes of inactivity and require re-authentication.
    applicability:
      - cloud-services
```

### Step 2: Create Policy Consuming the Guidance

When creating a Policy, you can import both the guidance documentations. 

The resolution mechanism automatically merges extensions with their base guidelines:
```yaml
title: Enterprise Security Policy
metadata:
  id: Enterprise-Security-Policy
  description: Security policy for enterprise systems
  author:
    id: security-team
    name: Security Team
    type: Human
  version: "v1.0.0"
  date: "2025-12-29"
  mapping-references:
    - id: NIST-800-53-Rev5
      title: NIST Special Publication 800-53 Revision 5
      version: "Rev 5"
      description: Complete NIST 800-53 Catalog
      url: "https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final"
    - id: Cloud-Access-Guidance
      title: Guidance for Cloud Access Management
      version: "v0.1.0"
      description: Domain-specific extensions to NIST 800-53 controls for cloud-based access management

imports:
  guidance:
    # Import upstream NIST guidance
    - reference-id: NIST-800-53-Rev5
      exclusions: [AC-19, AC-20]  # Exclude controls not applicable to this policy
    # Import cloud access extensions and independent guidelines
    - reference-id: Cloud-Access-Guidance
      constraints:
        - id: mfa-no-sms
          target-id: ORG-AC-3-EXT
          text: Multi-factor authentication must not use SMS-based methods. Only authenticator apps or hardware tokens are permitted.
```

**Key Points:**
- Import multiple guidance documents: broad (NIST) and domain-specific (Cloud-Access-Guidance)
- Extensions merge: `ORG-AC-3-EXT` extends `AC-3` from NIST-800-53-Rev5, so when both are imported, `AC-3` will include both the base NIST requirements and the cloud access extension
- Use `exclusions` to define which controls from guidance do not apply to your policy
- Constraints allow you to add organization-specific context to imported guidance
- Constraints reference the guideline ID (`target-id`) they apply to