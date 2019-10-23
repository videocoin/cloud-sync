package service

import "testing"

func TestParseReqPath(t *testing.T) {
	streamIDExpect := "5abd3f6b-4d3e-4c56-6b18-3ddb4f6f0370"
	segmentNumExpect := 1
	streamID, segmentNum, err := parseReqPath("5abd3f6b-4d3e-4c56-6b18-3ddb4f6f0370/1.ts")

	if err != nil {
		t.Fatalf("failed to parse request path: %s", err)
	}

	if streamID != streamIDExpect {
		t.Errorf("wrong stream id (actual: %s, expected: %s)", streamID, streamIDExpect)
	}

	if segmentNum != segmentNumExpect {
		t.Errorf("wrong segment num (actual: %d, expected: %d)", segmentNum, segmentNumExpect)
	}

	_, _, err = parseReqPath("5abd3f6b-4d3e-4c56-6b18-3ddb4f6f0370")
	if err == nil {
		t.Fatal("failed to parse request path")
	}

	_, _, err = parseReqPath("5abd3f6b-4d3e-4c56-6b18-3ddb4f6f0370/a.ts")
	if err == nil {
		t.Fatal("failed to parse request path")
	}
}
