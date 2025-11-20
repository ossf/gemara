# Finding Determination Specification

<!-- TOC -->
* [Finding Determination Specification](#finding-determination-specification)
  * [Abstract](#abstract)
  * [Overview](#overview)
  * [Notations and Terminology](#notations-and-terminology)
    * [Notational Conventions](#notational-conventions)
    * [Terminology](#terminology)
      * [Finding](#finding)
      * [Evaluator](#evaluator)
      * [Conflict Resolution Strategy](#conflict-resolution-strategy)
      * [Authoritative Evaluator](#authoritative-evaluator)
  * [Result Types](#result-types)
  * [Conflict Resolution Strategies](#conflict-resolution-strategies)
    * [Strict Strategy](#strict-strategy)
    * [ManualOverride Strategy](#manualoverride-strategy)
    * [AuthoritativeConfirmation Strategy](#authoritativeconfirmation-strategy)
  * [Use Cases](#use-cases)
    * [Strategy Selection Anti-Patterns](#strategy-selection-anti-patterns)
<!-- TOC -->

## Abstract

This specification defines how findings are determined when multiple assessment evaluators and procedures provide results for the same assessment requirement in Layer 4 Evaluation Plans.
It specifies conflict resolution strategies and result type semantics to ensure consistent finding determination across implementations.

## Overview

Layer 4 Evaluation Plans support multiple assessment evaluators running assessment procedures to evaluate control requirements. 
When multiple evaluators run the same procedure, or multiple procedures evaluate the same requirement, conflict resolution strategies determine how their results are combined to determine if there is a finding.
This specification provides a formal definition of these strategies to ensure consistent and predictable behavior across implementations.

## Notations and Terminology

### Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119).

### Terminology

This specification defines the following terms:

#### Finding

A finding is a documented observation that a Layer 2 control requirement, as referenced by a Layer 3 policy, is not being met or is not implemented correctly.

A finding is determined when one or more evaluator results indicate non-compliance with the assessed control requirement. Only results of `Failed`, `Unknown`, or `NeedsReview` constitute findings; 
`Passed` and `NotApplicable` results indicate compliance or inapplicability and do not produce findings.
Findings **MUST** be supported by evidence from assessment logs that document the evaluation results, and serve as the basis for Layer 5 enforcement actions and remediation activities.

#### Evaluator

An evaluator is a tool, process, or person that executes an assessment procedure and produces a result.

#### Conflict Resolution Strategy

A conflict resolution strategy is an algorithm that determines how multiple evaluator results are combined to produce a single finding determination.

#### Authoritative Evaluator

Evaluator authoritativeness determines how an evaluator participates in conflict resolution when using the `AuthoritativeConfirmation` strategy. 
Evaluators can be marked as `authoritative` (can trigger findings independently) or `non-authoritative` (requires confirmation from authoritative evaluators to trigger findings). 

The distinction between authoritative and non-authoritative evaluators maps directly to enforcement readiness:

- **Authoritative evaluators** represent procedures that are ready for enforcement. When authoritative evaluators report failures, findings are determined and can trigger enforcement actions (blocking, remediation, alerts).
- **Non-authoritative evaluators** represent procedures that collect risk data but are not ready to trigger enforcement. Non-authoritative evaluators provide risk visibility and inform decision-making, but their failures do not independently trigger enforcement actions.

**Default Behavior**: If the `authoritative` field is not explicitly set, evaluators default to `authoritative: false` (non-authoritative).
However, since the default conflict resolution strategy is `Strict`, this default has no effect on finding determination - `Strict` ignores the `authoritative` field and treats all evaluators equally.

## Result Types

The following result types are used in Layer 4:

- **NotRun**: The assessment was not executed
- **Passed**: The assessment passed successfully
- **Failed**: The assessment failed
- **NeedsReview**: The assessment requires manual review
- **NotApplicable**: The assessment is not applicable to the current context
- **Unknown**: The assessment result is unknown or indeterminate

## Conflict Resolution Strategies

For aggregating results within a single log (e.g., multiple steps within an assessment), implementations MUST use severity-based determination, where the most severe result takes precedence according to the hierarchy: `Failed > Unknown > NeedsReview > Passed > NotApplicable`.
If all results are `NotRun`, no finding **MUST** be determined.  This severity-based aggregation strategy **MUST** be used if a finding determination is inconclusive using other strategies.

Three multi-source conflict resolution strategies are defined in this specification. Implementations **MUST** support all strategies.

### Strict Strategy

The Strict strategy determines that a finding exists if ANY evaluator reports a failure, regardless of other evaluator results. This strategy provides zero tolerance for failures and is the simplest conflict resolution approach.

**Security-First Design**: `Strict` applies uniform zero-tolerance logic to all non-passing results (`Failed`, `Unknown`, `NeedsReview`). This makes `Strict` ideal for organizations that want predictable, consistent behavior and absolute zero-tolerance for security violations, ensuring that any evaluator reporting a problem triggers a finding.

**Process**: When using the Strict strategy, a finding **MUST** be determined according to the following priority order:

1. If ANY evaluator reports `Failed`, then **Finding exists** (Failed)
2. Else if ANY evaluator reports `Unknown`, then **Finding exists** (Unknown)
3. Else if ANY evaluator reports `NeedsReview`, then **Finding exists** (NeedsReview)
4. Else if ALL evaluator results are `Passed`, then **No finding** (Passed)
5. Else if ALL evaluator results are `NotApplicable`, then **No finding** (NotApplicable)

Evaluators with `NotRun` results MUST be excluded from the determination process. All evaluators are treated equally when using the Strict strategy.

### ManualOverride Strategy

The ManualOverride strategy gives precedence to manual review evaluators over automated evaluators when determining findings from conflicting results.

**Process**: When using the ManualOverride strategy:

1. Separate results into manual and automated evaluator results based on the evaluator's ExecutionType (`Manual` vs `Automated`).
2. If manual evaluators exist:
   - If any manual evaluator reports `Failed`: **Finding exists** (Failed)
   - Else if any manual evaluator reports `Unknown`: **Finding exists** (Unknown)
   - Else if any manual evaluator reports `NeedsReview`: **Finding exists** (NeedsReview)
   - Else if all manual evaluators report `Passed`: **No finding** (Passed)
   - Else: **No finding** (NotApplicable)

Evaluators with `NotRun` results **MUST** be excluded from the determination process.

### AuthoritativeConfirmation Strategy

The AuthoritativeConfirmation strategy treats non-authoritative evaluators as requiring confirmation from authoritative evaluators before triggering findings.

**Process**: When using the AuthoritativeConfirmation strategy:

1. Separate evaluators into authoritative and non-authoritative groups based on their `authoritative` field. Evaluators without an explicit `authoritative` field default to `authoritative: false` (non-authoritative).
2. Authoritative evaluators trigger findings using Strict logic:
   - If any authoritative evaluator reports `Failed`: **Finding exists** (Failed)
   - Else if any authoritative evaluator reports `Unknown`: **Finding exists** (Unknown)
   - Else if any authoritative evaluator reports `NeedsReview`: **Finding exists** (NeedsReview)
   - Else if all authoritative evaluators report `Passed`: Continue to step 3
3. Non-authoritative evaluators require confirmation:
   - If only non-authoritative evaluators report `Failed`: **No finding** (non-authoritative cannot trigger alone)
   - If non-authoritative evaluator reports `Failed` AND authoritative evaluator reports `Passed`: **No finding** (contradicted)
   - If non-authoritative evaluator reports `Failed` AND authoritative evaluator reports `Failed`: **Finding exists** (confirmed)
   - If non-authoritative evaluator reports `Failed` AND authoritative evaluator reports `Unknown` or `NeedsReview`: **Finding exists** (escalated for investigation)
   - If all evaluators (authoritative and non-authoritative) report `Passed`: **No finding** (Passed)
   - If all evaluators report `NotApplicable`: **No finding** (NotApplicable)

Evaluators with `NotRun` results MUST be excluded from the determination process.

**Key Behaviors**:

- **Authoritative evaluators** can trigger findings independently using zero-tolerance logic. If any authoritative evaluator reports a non-passing result, a finding is immediately determined, regardless of non-authoritative evaluator results.
- **Non-authoritative evaluators** can only trigger findings when:
  - They confirm an authoritative evaluator's failure (both report `Failed`)
  - They escalate an unclear authoritative result (authoritative reports `Unknown`/`NeedsReview` and non-authoritative reports `Failed`)
- **Non-authoritative evaluators cannot**:
  - Trigger findings independently (if only non-authoritative evaluators report failures, no finding is determined)
  - Override authoritative evaluators (if authoritative evaluators report `Passed`, non-authoritative failures are ignored)

## Use Cases

Layer 4 Evaluation Plans serve as the foundation for Layer 5 enforcement decisions. The evaluation results inform what enforcement actions should be taken, such as blocking deployments, triggering remediation, or generating alerts. 
However, not all procedures in an Evaluation Plan need to trigger enforcement actions.

Use the following decision matrix to select the appropriate conflict resolution strategy:

| Scenario                                                             | Recommended Strategy        | Rationale                                                                                                                                             |
|----------------------------------------------------------------------|-----------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| **All evaluators are equally trusted and validated**                 | `Strict`                    | Zero tolerance for failures. Any evaluator reporting a problem triggers a finding. Simplest and most predictable behavior.                            |
| **Multiple evaluators, need human judgment for final determination** | `ManualOverride`            | Manual review takes precedence. Automated tools provide initial screening, but human reviewers make final decisions.                                  |
| **Gradual rollout of new enforcement**                               | `AuthoritativeConfirmation` | Start with evaluators as non-authoritative to collect baseline data. Promote to authoritative once validated.                                         |
| **Experimental or unvalidated evaluators**                           | `AuthoritativeConfirmation` | Mark experimental evaluators as non-authoritative. They provide visibility but don't trigger enforcement until confirmed by authoritative evaluators. |
| **Simple, predictable behavior needed**                              | `Strict`                    | No complex logic. Any failure = finding. Easiest to understand and maintain.                                                                          |
| **Automated tools need human verification**                          | `ManualOverride`            | Automated tools flag issues, but require manual confirmation before triggering findings.                                                              |
| **Collecting metrics before enforcing**                              | `AuthoritativeConfirmation` | Run evaluators non-authoritatively to understand violation patterns, then promote to authoritative for enforcement.                                   |

**Quick Decision Guide:**

```
Do you need human judgment for final decisions?
├─ YES → Use ManualOverride
└─ NO → Do you need to distinguish enforcement-ready from experimental evaluators?
    ├─ YES → Use AuthoritativeConfirmation
    └─ NO → Use Strict (default)
```

A typical workflow for introducing promoting procedures from non-authoritative to authoritative:

1. **Add to Evaluation Plan (Non-Authoritative)**
   - Add new procedure with evaluators (default is `authoritative: false`, so no explicit setting needed)
   - Use `AuthoritativeConfirmation` strategy to enable authoritative/non-authoritative behavior
   - Run evaluations to collect baseline data
   - Understand violation patterns and impact
   - Validate that the evaluator produces accurate results

2. **Assess and Remediate**
   - Analyze non-authoritative evaluator results
   - Fix critical violations discovered during baseline collection
   - Validate procedure accuracy and false positive rates
   - Ensure evaluator is ready for enforcement

3. **Promote to Authoritative**
   - Change evaluator `authoritative` field to `true` (explicit opt-in required)
   - Now triggers findings and enforcement actions independently
   - Monitor enforcement impact
   - Verify that enforcement actions are appropriate

4. **Maintain Some Procedures as Non-Authoritative**
   - Keep informational/audit-only procedures as non-authoritative (default `authoritative: false`)
   - Maintain experimental procedures as non-authoritative until validated
   - Use non-authoritative for risk visibility without enforcement
   - With `AuthoritativeConfirmation` strategy, non-authoritative evaluators require confirmation

**Example: Understanding Defaults**

```yaml
# Example 1: Default behavior (security-first)
procedures:
  - id: check-branch-protection
    evaluators:
      - id: security-scanner
        # No authoritative field = defaults to false (but ignored by Strict)
        # No strategy = defaults to Strict
        # Result: Any failure from security-scanner triggers a finding (Strict ignores authoritative)
```

```yaml
# Example 2: Explicit non-authoritative (opt-out)
procedures:
  - id: experimental-check
    evaluators:
      - id: new-tool
        authoritative: false  # Explicitly opt-out
    strategy:
      conflict-rule-type: AuthoritativeConfirmation
      # Result: new-tool provides visibility but doesn't trigger findings alone
```

### Strategy Selection Anti-Patterns

**Avoid these patterns:**

- ❌ **Using `AuthoritativeConfirmation` with all evaluators as non-authoritative** - No findings will ever be triggered. At least one authoritative evaluator is required.
- ❌ **Using `ManualOverride` when all evaluators have ExecutionType `Automated`** - Falls back to `Strict`, defeating the purpose of manual override. At least one evaluator with ExecutionType `Manual` is required for ManualOverride to function as intended.
- ❌ **Mixing strategies inconsistently** - Use the same strategy at both Assessment and Procedure levels unless there's a clear reason for different behavior.
- ❌ **Setting `authoritative: false` without using `AuthoritativeConfirmation`** - The field is ignored by other strategies, which can be confusing.
