package helper

import "testing"

var (
	varStrNoAppendix = "${gateway:test-gateway}"
	varStrAppendix   = "${gateway:test-gateway}/test"
)

func TestVarNoAppendix(t *testing.T) {
	v := ExtractVar(varStrNoAppendix)

	resource := "gateway"
	id := "test-gateway"

	if v.Resource != resource {
		t.Fatalf("expected %s, got %s", resource, v.Resource)
	}

	if v.ID != id {
		t.Fatalf("expected %s, got %s", id, v.ID)
	}
}

func TestVarAppendixt(t *testing.T) {
	v := ExtractVar(varStrAppendix)

	resource := "gateway"
	id := "test-gateway"
	appendix := "/test"

	if v.Resource != resource {
		t.Fatalf("expected %s, got %s", resource, v.Resource)
	}

	if v.ID != id {
		t.Fatalf("expected %s, got %s", id, v.ID)
	}

	if v.Appendix != appendix {
		t.Fatalf("expected %s, got %s", appendix, v.Appendix)
	}
}
