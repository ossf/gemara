# Finding Determination Specification

<!-- TOC -->
* [Finding Determination Specification](#finding-determination-specification)
  * [Abstract](#abstract)
  * [Overview](#overview)
  * [Notations and Terminology](#notations-and-terminology)
    * [Notational Conventions](#notational-conventions)
    * [Terminology](#terminology)
      * [Finding](#finding)
      * [Executor](#executor)
      * [Conflict Resolution Strategy](#conflict-resolution-strategy)
      * [Executor Role](#executor-role)
  * [Result Types](#result-types)
  * [Conflict Resolution Strategies](#conflict-resolution-strategies)
    * [Strict Strategy](#strict-strategy)
    * [ManualOverride Strategy](#manualoverride-strategy)
    * [AdvisoryRequiresConfirmation Strategy](#advisoryrequiresconfirmation-strategy)
  * [Use Cases](#use-cases)
    * [Evaluation Plans Guide Enforcement Decisions](#evaluation-plans-guide-enforcement-decisions)
    * [Primary vs Advisory: Enforcement-Ready vs Data Collection](#primary-vs-advisory-enforcement-ready-vs-data-collection)
    * [Workflow: Promoting Procedures from Advisory to Primary](#workflow-promoting-procedures-from-advisory-to-primary)
  * [Default Behavior](#default-behavior)
<!-- TOC -->

## Abstract

This specification defines how findings are determined when multiple assessment executors and procedures provide results for the same assessment requirement in Layer 4 Evaluation Plans.
It specifies conflict resolution strategies and result type semantics to ensure consistent finding determination across implementations.

## Overview

Layer 4 Evaluation Plans support multiple assessment executors running assessment procedures to evaluate control requirements. 
When multiple executors run the same procedure, or multiple procedures evaluate the same requirement, conflict resolution strategies determine how their results are combined to determine if there is a finding.
This specification provides a formal definition of these strategies to ensure consistent and predictable behavior across implementations.

## Notations and Terminology

### Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119).

### Terminology

This specification defines the following terms:

#### Finding

A finding is a documented observation that a Layer 2 control requirement, as referenced by a Layer 3 policy, is not being met or is not implemented correctly.

A finding is determined when one or more executor results indicate non-compliance with the assessed control requirement. Only results of `Failed`, `Unknown`, or `NeedsReview` constitute findings; 
`Passed` and `NotApplicable` results indicate compliance or inapplicability and do not produce findings.
Findings MUST be supported by evidence from assessment logs that document the evaluation results, and serve as the basis for Layer 5 enforcement actions and remediation activities.

#### Executor

An executor is a tool, process, or person that executes an assessment procedure and produces a result.

#### Conflict Resolution Strategy

A conflict resolution strategy is an algorithm that determines how multiple executor results are combined to produce a single finding determination.

#### Executor Role

An executor role determines how an executor participates in conflict resolution when using the `AdvisoryRequiresConfirmation` strategy. Executors can be assigned the role of `Primary` (can trigger findings independently) or `Advisory` (requires confirmation from Primary executors to trigger findings). If no role is explicitly assigned, executors default to `Primary` role.

## Result Types

The following result types are used in Layer 4:

- **NotRun**: The assessment was not executed
- **Passed**: The assessment passed successfully
- **Failed**: The assessment failed
- **NeedsReview**: The assessment requires manual review
- **NotApplicable**: The assessment is not applicable to the current context
- **Unknown**: The assessment result is unknown or indeterminate

## Conflict Resolution Strategies

Three conflict resolution strategies are defined in this specification. Implementations MUST support all strategies.

### Strict Strategy

The Strict strategy determines that a finding exists if ANY executor reports a failure, regardless of other executor results. This strategy provides zero tolerance for failures and is the simplest conflict resolution approach.

**Security-First Design**: `Strict` applies uniform zero-tolerance logic to all non-passing results (`Failed`, `Unknown`, `NeedsReview`). This makes `Strict` ideal for organizations that want predictable, consistent behavior and absolute zero-tolerance for security violations, ensuring that any executor reporting a problem triggers a finding.

**Process**: When using the Strict strategy, a finding MUST be determined according to the following priority order:

1. If ANY executor reports `Failed`, then **Finding exists** (Failed)
2. Else if ANY executor reports `Unknown`, then **Finding exists** (Unknown)
3. Else if ANY executor reports `NeedsReview`, then **Finding exists** (NeedsReview)
4. Else if ALL executor results are `Passed`, then **No finding** (Passed)
5. Else if ALL executor results are `NotApplicable`, then **No finding** (NotApplicable)

Executors with `NotRun` results MUST be excluded from the determination process. All executors are treated equally when using the Strict strategy.

### ManualOverride Strategy

The ManualOverride strategy gives precedence to manual review executors over automated executors when determining findings from conflicting results.

**Process**: When using the ManualOverride strategy:

1. Separate results into manual and automated executor results based on executor type.
2. If manual executors exist:
   - If any manual executor reports `Failed`: **Finding exists** (Failed)
   - Else if any manual executor reports `Unknown`: **Finding exists** (Unknown)
   - Else if any manual executor reports `NeedsReview`: **Finding exists** (NeedsReview)
   - Else if all manual executors report `Passed`: **No finding** (Passed)
   - Else: **No finding** (NotApplicable)
3. If no manual executors exist, determine finding from automated results using severity hierarchy: `Failed > Unknown > NeedsReview > Passed`

Executors with `NotRun` results MUST be excluded from the determination process.

### AdvisoryRequiresConfirmation Strategy

The AdvisoryRequiresConfirmation strategy treats Advisory executors as requiring confirmation from Primary executors before triggering findings.

**Process**: When using the AdvisoryRequiresConfirmation strategy:

1. Separate executors into Primary and Advisory groups based on their `role` field. Executors without an explicit role default to Primary.
2. Primary executors trigger findings using Strict logic:
   - If any Primary executor reports `Failed`: **Finding exists** (Failed)
   - Else if any Primary executor reports `Unknown`: **Finding exists** (Unknown)
   - Else if any Primary executor reports `NeedsReview`: **Finding exists** (NeedsReview)
   - Else if all Primary executors report `Passed`: Continue to step 3
3. Advisory executors require confirmation:
   - If only Advisory executors report `Failed`: **No finding** (advisory cannot trigger alone)
   - If Advisory executor reports `Failed` AND Primary executor reports `Passed`: **No finding** (contradicted)
   - If Advisory executor reports `Failed` AND Primary executor reports `Failed`: **Finding exists** (confirmed)
   - If Advisory executor reports `Failed` AND Primary executor reports `Unknown` or `NeedsReview`: **Finding exists** (escalated for investigation)
   - If all executors (Primary and Advisory) report `Passed`: **No finding** (Passed)
   - If all executors report `NotApplicable`: **No finding** (NotApplicable)

Executors with `NotRun` results MUST be excluded from the determination process.

**Key Behaviors**:

- **Primary executors** can trigger findings independently using zero-tolerance logic. If any Primary executor reports a non-passing result, a finding is immediately determined, regardless of Advisory executor results.
- **Advisory executors** can only trigger findings when:
  - They confirm a Primary executor's failure (both report `Failed`)
  - They escalate an unclear Primary result (Primary reports `Unknown`/`NeedsReview` and Advisory reports `Failed`)
- **Advisory executors cannot**:
  - Trigger findings independently (if only Advisory executors report failures, no finding is determined)
  - Override Primary executors (if Primary executors report `Passed`, Advisory failures are ignored)

## Use Cases

### Evaluation Plans Guide Enforcement Decisions

Layer 4 Evaluation Plans serve as the foundation for Layer 5 enforcement decisions. The evaluation results inform what enforcement actions should be taken, such as blocking deployments, triggering remediation, or generating alerts. However, not all procedures in an Evaluation Plan need to trigger enforcement actions.

### Primary vs Advisory: Enforcement-Ready vs Data Collection

The distinction between Primary and Advisory executors maps directly to enforcement readiness:

- **Primary executors** represent procedures that are ready for enforcement. When Primary executors report failures, findings are determined and can trigger enforcement actions (blocking, remediation, alerts).

- **Advisory executors** represent procedures that collect risk data but are not ready to trigger enforcement. Advisory executors provide risk visibility and inform decision-making, but their failures do not independently trigger enforcement actions.

### Workflow: Promoting Procedures from Advisory to Primary

A typical workflow for introducing new enforcement:

1. **Add to Evaluation Plan (Advisory)**
   - Add new procedure with executors in `Advisory` role
   - Run evaluations to collect baseline data
   - Understand violation patterns and impact

2. **Assess and Remediate**
   - Analyze Advisory executor results
   - Fix critical violations
   - Validate procedure accuracy

3. **Promote to Primary**
   - Change executor roles from `Advisory` to `Primary`
   - Now triggers findings and enforcement actions
   - Monitor enforcement impact

4. **Maintain Some Procedures as Advisory**
   - Keep informational/audit-only procedures in Advisory
   - Maintain experimental procedures in Advisory until validated
   - Use Advisory for risk visibility without enforcement

## Default Behavior

When aggregating results from multiple executors or procedures (log aggregation), if no conflict resolution strategy is explicitly specified, implementations MUST use the `Strict` strategy.

For aggregating results within a single log (e.g., multiple steps within an assessment), implementations MUST use severity-based determination, where the most severe result takes precedence according to the hierarchy: `Failed > Unknown > NeedsReview > Passed > NotApplicable`.

If all results are `NotRun`, no finding MUST be determined.
