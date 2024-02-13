// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License
//
// Note: The safety settings listed below will not affect complexity (same as Command Registry), unlike 'if', 'for', 'case', '&&' or '||' which would increase complexity.
// For instance, you can introduce numerous settings (e.g., low, high, default, etc) here without impacting the complexity.

package terminal

import (
	genai "github.com/google/generative-ai-go/genai"
)

// DefaultSafetySettings returns a SafetySettings instance with a default
// configuration where all categories are set to block medium and above levels
// of harmful content. This default setting provides a balanced approach to
// content safety, suitable for general use cases.
func DefaultSafetySettings() *SafetySettings {
	return &SafetySettings{
		DangerousContentThreshold:        genai.HarmBlockMediumAndAbove,
		HarassmentContentThreshold:       genai.HarmBlockMediumAndAbove,
		SexuallyExplicitContentThreshold: genai.HarmBlockMediumAndAbove,
		MedicalThreshold:                 genai.HarmBlockMediumAndAbove,
		ViolenceThreshold:                genai.HarmBlockMediumAndAbove,
		HateSpeechThreshold:              genai.HarmBlockMediumAndAbove,
		ToxicityThreshold:                genai.HarmBlockMediumAndAbove,
		DerogatoryThershold:              genai.HarmBlockMediumAndAbove,
	}
}

// SetLowSafety adjusts the safety settings to a lower threshold, allowing more
// content through the filter. This setting may be appropriate for environments
// where content restrictions can be more relaxed, or where users are expected
// to handle a wider range of content types.
func (s *SafetySettings) SetLowSafety() {
	s.DangerousContentThreshold = genai.HarmBlockLowAndAbove
	s.HarassmentContentThreshold = genai.HarmBlockLowAndAbove
	s.SexuallyExplicitContentThreshold = genai.HarmBlockLowAndAbove
	s.MedicalThreshold = genai.HarmBlockLowAndAbove
	s.ViolenceThreshold = genai.HarmBlockLowAndAbove
	s.HateSpeechThreshold = genai.HarmBlockLowAndAbove
	s.ToxicityThreshold = genai.HarmBlockLowAndAbove
	s.DerogatoryThershold = genai.HarmBlockLowAndAbove
}

// SetHighSafety raises the safety settings to a higher threshold, providing
// stricter content filtering. This setting is useful in environments that
// require a high degree of content moderation to ensure user safety or to
// comply with strict regulatory standards.
func (s *SafetySettings) SetHighSafety() {
	s.DangerousContentThreshold = genai.HarmBlockOnlyHigh
	s.HarassmentContentThreshold = genai.HarmBlockOnlyHigh
	s.SexuallyExplicitContentThreshold = genai.HarmBlockOnlyHigh
	s.MedicalThreshold = genai.HarmBlockOnlyHigh
	s.ViolenceThreshold = genai.HarmBlockOnlyHigh
	s.HateSpeechThreshold = genai.HarmBlockOnlyHigh
	s.ToxicityThreshold = genai.HarmBlockOnlyHigh
	s.DerogatoryThershold = genai.HarmBlockOnlyHigh
}

// SetUnspecifiedSafety sets the safety settings to an unspecified threshold, which means no specific blocking level is set.
// This setting may be appropriate for environments where the default AI model behavior is desired without additional filters.
func (s *SafetySettings) SetUnspecifiedSafety() {
	s.DangerousContentThreshold = genai.HarmBlockUnspecified
	s.HarassmentContentThreshold = genai.HarmBlockUnspecified
	s.SexuallyExplicitContentThreshold = genai.HarmBlockUnspecified
	s.MedicalThreshold = genai.HarmBlockUnspecified
	s.ViolenceThreshold = genai.HarmBlockUnspecified
	s.HateSpeechThreshold = genai.HarmBlockUnspecified
	s.ToxicityThreshold = genai.HarmBlockUnspecified
	s.DerogatoryThershold = genai.HarmBlockUnspecified
}

// SetNoneSafety sets the safety settings to no threshold, which means no content is blocked.
// This setting may be appropriate for environments where content moderation is not required or where users are expected
// to handle all types of content without any filtering.
func (s *SafetySettings) SetNoneSafety() {
	s.DangerousContentThreshold = genai.HarmBlockNone
	s.HarassmentContentThreshold = genai.HarmBlockNone
	s.SexuallyExplicitContentThreshold = genai.HarmBlockNone
	s.MedicalThreshold = genai.HarmBlockNone
	s.ViolenceThreshold = genai.HarmBlockNone
	s.HateSpeechThreshold = genai.HarmBlockNone
	s.ToxicityThreshold = genai.HarmBlockNone
	s.DerogatoryThershold = genai.HarmBlockNone
}

// ApplyToModel applies the configured safety settings to a given generative AI model.
// This method updates the model's safety settings to match the thresholds specified
// in the SafetySettings instance, affecting how the model filters generated content.
func (s *SafetySettings) ApplyToModel(model *genai.GenerativeModel, modelName string) {
	// fix 400 error lmao, should be work now
	// Note: This is subject to change to avoid stupid unnecessary complexity, especially when dealing with numerous models.
	// For instance, simplify the process by breaking down the logic into smaller components.
	// Keeping cyclomatic complexity under 5 is a secret key hahaha in Go programming. It leads to reusable, easy-to-maintain code that boosts performance and minimizes bugs.
	switch modelName {
	case GeminiPro:
		// Apply a specific set of safety settings for the "gemini-pro" model
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: s.DangerousContentThreshold,
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: s.HarassmentContentThreshold,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: s.SexuallyExplicitContentThreshold,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: s.HateSpeechThreshold,
			},
		}
	// TODO: Implement Perfect System in Terminal For Gemini-Pro Vision Capabilities
	case GeminiProVision:
		// Apply a specific set of safety settings for the "gemini-pro-vision" model
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: s.DangerousContentThreshold,
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: s.HarassmentContentThreshold,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: s.SexuallyExplicitContentThreshold,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: s.HateSpeechThreshold,
			},
		}
		// TODO: This for other model
		// would be implemented in next years maybe lmao
	default:
		// Apply a different set of safety settings for other models
		model.SafetySettings = []*genai.SafetySetting{
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: s.DangerousContentThreshold,
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: s.HarassmentContentThreshold,
			},
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: s.SexuallyExplicitContentThreshold,
			},
			{
				Category:  genai.HarmCategoryMedical,
				Threshold: s.MedicalThreshold,
			},
			{
				Category:  genai.HarmCategoryViolence,
				Threshold: s.ViolenceThreshold,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: s.HateSpeechThreshold,
			},
			{
				Category:  genai.HarmCategoryToxicity,
				Threshold: s.ToxicityThreshold,
			},
			{
				Category:  genai.HarmCategoryDerogatory,
				Threshold: s.DerogatoryThershold,
			},
		}
	}
}
