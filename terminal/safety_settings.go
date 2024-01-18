// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License
// Note: The safety settings listed below will not affect complexity (same as Command Registery), unlike 'if', 'for', 'case', '&&' or '||' which would increase complexity.
// For instance, you can introduce numerous settings (e.g., low, high, default, etc) here without impacting the complexity.

package terminal

import (
	genai "github.com/google/generative-ai-go/genai"
)

// SafetyOption is a function type that takes a pointer to a SafetySettings
// instance and applies a specific safety configuration to it. It is used
// to abstract the different safety level settings (e.g., low, high, default)
// and allows for a flexible and scalable way to manage safety configurations
// through function mapping.
type SafetyOption struct {
	Setter func(s *SafetySettings)
	Valid  bool
}

// SafetySettings encapsulates the content safety configuration for the AI model.
// It defines thresholds for various categories of potentially harmful content,
// allowing users to set the desired level of content filtering based on the
// application's requirements and user preferences.
type SafetySettings struct {
	// DangerousContentThreshold defines the threshold for filtering dangerous content.
	DangerousContentThreshold genai.HarmBlockThreshold
	// HarassmentContentThreshold defines the threshold for filtering harassment-related content.
	HarassmentContentThreshold genai.HarmBlockThreshold
	// SexuallyExplicitContentThreshold defines the threshold for filtering sexually explicit content.
	SexuallyExplicitContentThreshold genai.HarmBlockThreshold
	// MedicalThreshold defines the threshold for filtering medical-related content.
	MedicalThreshold genai.HarmBlockThreshold
	// ViolenceThreshold defines the threshold for filtering violent content.
	ViolenceThreshold genai.HarmBlockThreshold
	// HateSpeechThreshold defines the threshold for filtering hate speech.
	HateSpeechThreshold genai.HarmBlockThreshold
	// ToxicityThreshold defines the threshold for filtering toxic content.
	ToxicityThreshold genai.HarmBlockThreshold
}

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
}

// ApplyToModel applies the configured safety settings to a given generative AI model.
// This method updates the model's safety settings to match the thresholds specified
// in the SafetySettings instance, affecting how the model filters generated content.
func (s *SafetySettings) ApplyToModel(model *genai.GenerativeModel) {
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
	}
}
