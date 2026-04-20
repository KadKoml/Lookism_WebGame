package usecase

import (
	"math"
	"testing"
)

func TestCalculateFinalStat_Level1_NoStars(t *testing.T) {
	result := CalculateFinalStat(100, 1, 0)
	if result != 100.0 {
		t.Errorf("Expected 100.0, got %f", result)
	}
}

func TestCalculateFinalStat_Level10_NoStars(t *testing.T) {
	// 100 * 1.0 * (1 + 9*0.03) = 100 * 1.27 = 127
	result := CalculateFinalStat(100, 10, 0)
	expected := 127.0
	if result != expected {
		t.Errorf("Expected %f, got %f", expected, result)
	}
}

func TestCalculateFinalStat_Level1_3Stars(t *testing.T) {
	// 100 * (1 + 3*0.05) * 1.0 = 100 * 1.15 = 115
	result := CalculateFinalStat(100, 1, 3)
	expected := 115.0
	if math.Abs(result-expected) > 0.001 {
		t.Errorf("Expected %f, got %f", expected, result)
	}
}

func TestCalculateFinalStat_Level20_5Stars(t *testing.T) {
	// 100 * (1 + 5*0.05) * (1 + 19*0.03) = 100 * 1.25 * 1.57 = 196.25
	result := CalculateFinalStat(100, 20, 5)
	expected := 196.25
	if math.Abs(result-expected) > 0.001 {
		t.Errorf("Expected %f, got %f", expected, result)
	}
}

func TestCalculateFinalStat_ClampMinLevel(t *testing.T) {
	// Level < 1 should be clamped to 1
	result := CalculateFinalStat(100, -5, 0)
	if result != 100.0 {
		t.Errorf("Expected 100.0 for negative level, got %f", result)
	}
}

func TestCalculateFinalStat_ClampMaxLevel(t *testing.T) {
	// Level > 60 should be clamped to 60
	result := CalculateFinalStat(100, 999, 0)
	expected := CalculateFinalStat(100, 60, 0)
	if result != expected {
		t.Errorf("Expected %f for clamped level, got %f", expected, result)
	}
}

func TestCalculateFinalIntStat(t *testing.T) {
	result := CalculateFinalIntStat(430, 12, 0)
	// 430 * 1.0 * (1 + 11*0.03) = 430 * 1.33 = 571.9 -> 572
	expected := int32(572)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
