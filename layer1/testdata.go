package layer1

import "time"

func goodAIGFExample() GuidanceDocument {
	return GuidanceDocument{
		Metadata: Metadata{
			Id:      "FINOS-AIR",
			Title:   "AI Governance Framework",
			Version: "0.1.0",
			MappingReferences: []MappingReference{
				{
					Id:      "NIST-800-53",
					Title:   "NIST SP 800-53r5",
					Version: "rev5",
					Url:     "https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-53r5.pdf#%5B%7B%22num%22%3A342%2C%22gen%22%3A0%7D%2C%7B%22name%22%3A%22XYZ%22%7D%2C88%2C310%2C0%5D",
				},
				{
					Id:      "AIR-PRIN",
					Title:   "Example Principles Document for the Framework",
					Version: "0.1.0",
				},
			},
			DocumentType: "Framework",
			Applicability: &Applicability{
				TechnologyDomains: []string{
					"artificial-intelligence",
				},
				IndustrySectors: []string{
					"financial-services",
				},
			},
			LastModified: time.Now().Format(time.RFC3339),
		},
		FrontMatter: "The following framework has been developed by FINOS (Fintech Open Source Foundation).",
		Categories: []Category{
			{
				Id:          "DET",
				Title:       "Detective",
				Description: "Detection and Continuous Improvement",
				Guidelines: []Guideline{
					{
						Id:    "AIR-DET-011",
						Title: "Human Feedback Loop for AI Systems",
						Objective: "A Human Feedback Loop is a critical detective and continuous improvement mechanism that involves systematically collecting, " +
							"analyzing, and acting upon feedback provided by human users, subject matter experts (SMEs), or reviewers regarding an AI system’s performance, outputs, or behavior.",
						Rationale: &Rationale{
							Risks: []Risk{},
							Outcomes: []Outcome{
								{
									Title:       "Governance Support",
									Description: "Provides data for AI governance bodies to monitor impact and make decisions",
								},
							},
						},
						GuidelineParts: []Part{
							{
								Id:    "AIR-DET-011.1",
								Title: "Designing the Feedback Mechanism",
								Prose: "Implementing an effective human feedback loop involves careful design of the mechanism.",
								Recommendations: []string{
									"Define Intended Use and KPIs:\n" +
										"Objectives: Clearly document how feedback data will be utilized, such as for prompt fine-tuning, RAG document updates," +
										"model/data drift detection, or more advanced uses like Reinforcement Learning from Human Feedback (RLHF).\nKPI Alignment: Design feedback questions and metrics " +
										"to align with the solution’s key performance indicators (KPIs). For example, if accuracy is a KPI, feedback might involve users or SMEs annotating if an answer was correct.",
								},
							},
							{
								Id:    "AIR-DET-011.2",
								Title: "Types of Feedback and Collection Methods",
								Prose: "Implementing an effective human feedback loop involves clear collection processes.",
								Recommendations: []string{
									"Quantitative Feedback:\n" +
										"Description: Involves collecting structured responses that can be easily aggregated and measured, such as numerical ratings " +
										"(e.g., “Rate this response on a scale of 1-5 for helpfulness”), categorical choices (e.g., “Was this answer: Correct/Incorrect/Partially Correct”), " +
										"or binary responses (e.g., thumbs up/down).\nUse Cases: Effective for tracking trends, measuring against KPIs, and quickly identifying areas of high or low performance.",
								},
							},
						},
						GuidelineMappings: []Mapping{
							{
								ReferenceId: "NIST-800-53",
								Entries: []MappingEntry{
									{
										ReferenceId: "CA-7",
										Strength:    7,
										Remarks:     "This control is closely related to CA-7.",
									},
									{
										ReferenceId: "IR-6",
										Strength:    5,
										Remarks:     "This control has some relevance to IR-6.",
									},
									{
										ReferenceId: "PM-26",
										Strength:    3,
										Remarks:     "This control is loosely related to PM-26.",
									},
									{
										ReferenceId: "RA-5",
										Strength:    7,
										Remarks:     "This control is closely related to RA-5.",
									},
									{
										ReferenceId: "SI-2",
										Strength:    5,
										Remarks:     "This control has some relevance to SI-2.",
									},
								},
							},
						},
						PrincipleMappings: []Mapping{
							{
								ReferenceId: "AIR-PRIN",
								Entries: []MappingEntry{
									{
										ReferenceId: "TIMELINESS",
										Strength:    7,
										Remarks:     "This principle emphasizes the importance of timely feedback.",
									},
								},
							},
						},
						SeeAlso: []string{
							"AIR-DET-015",
							"AIR-DET-004",
							"AIR-PREV-005",
						},
					},
				},
			},
		},
	}
}
