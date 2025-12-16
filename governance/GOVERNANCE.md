# Gemara Project Governance

As a developing project, Gemara aims to have a quick development cycle where decisions and community issues are resolved promptly while capturing the input of interested stakeholders.

Gemara has no formal collegiate body in charge of steering. Decisions are guided by the consensus of community members who have achieved maintainer status.

While maintainer consensus shall be the process for decision making, all issues and proposals shall be governed by the project's [Guiding Governance Principles](#guiding-governance-principles).

## Guiding Governance Principles

Any issues or proposals brought to the project's maintainers shall be framed in the [Guiding Governance Principles](#guiding-governance-principles).
Proposals not adhering to said principles shall not be considered for consensus.

### Follow Layer-Based Architecture

The six-layer model (Guidance, Controls, Policy, Evaluation, Enforcement, Audit) provides a clear structure for organizing compliance activities. Changes respect this architectural model and relationships between layers.

### Ensure Engineering-Centric Design

Schemas and models must prioritize practical, implementable solutions aligned with how GRC professionals apply engineering practices to compliance work, ensuring relevance and usability.

### Use Schema-Driven Development

Machine-readable schemas (CUE) form the foundation for all compliance activities, ensuring consistency, validation, and automation throughout the GRC lifecycle.

### Incremental and Backward-Compatible

Changes prioritize backward compatibility. Breaking changes are rare and require careful consideration, community input, and clear migration paths.

## Maintainer Consensus

To reach a decision on an issue or proposal, the proponents must seek maintainer consensus.

In the context of this document, "maintainer consensus" means collecting approvals from at least 51% of the current maintainer body, with enough time for all maintainers to review (usually 2 business days), and without a dissenting maintainer opinion.

This document does not prescribe a method of voting. Any mechanism that enables the collection of positive/negative votes associated with an identity may be used. Examples of this include voting through "thumbs up/down" emojis or with "+1" comments in issues.

## Maintainer Status

Community members may become maintainer candidates through:

- Nomination by a [sponsoring committee] at any time
- Self-nomination after actively contributing to Gemara monthly for six months or more

Nominations are submitted via pull request to update Gemara's [MAINTAINERS.md]. After validation, [maintainer consensus] is sought. Upon consensus, the PR is merged to confirm the new maintainer.

### Sponsoring Committees

A sponsoring committee must have at least two members. When maintainers represent three or more organizations, committee members must be from different organizations.

### Continued Maintainer Status

Maintainer status requires regular activity and adherence to the [OpenSSF Code of Conduct](https://openssf.org/community/code-of-conduct/).

### Emeritus Maintainers

Emeritus maintainers are listed in a separate section on Gemara's [MAINTAINERS.md].
A maintainer may be given Emeritus status after six months of inactivity (e.g., no pull request or issue interactions) or may self-assign Emeritus status via pull request.
A maintainer may return from Emeritus status through [maintainer consensus] and a pull request.

## Revisions to the Governance Model

The governance model is revisited every six months to address community needs. At any point, a GitHub issue may propose changes to governance. Proposals require approval from at least 66% of active maintainers.

## Acknowledgements

This document was adapted from the Security Baseline Governance [documentation](https://github.com/ossf/security-baseline/blob/main/governance/GOVERNANCE.md).

[MAINTAINERS.md]: ./MAINTAINERS.md
[maintainer consensus]: #maintainer-consensus
[Sponsoring Committee]: #sponsoring-committees