// Package sarif provides conversion functions to transform Gemara evaluation
// results into SARIF (Static Analysis Results Interchange Format) format.
//
// SARIF is a standard format for static analysis tool output, enabling
// integration with code scanning platforms like GitHub Code Scanning,
// Azure DevOps, and other security analysis tools.
//
// This package converts EvaluationLog entries into SARIF v2.1.0 format,
// where each AssessmentLog becomes a SARIF result. The conversion supports
// optional catalog enrichment to include control and requirement details
// in the SARIF output.
package sarif
