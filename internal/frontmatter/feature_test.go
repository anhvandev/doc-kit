package frontmatter

import (
	"bytes"
	"testing"
)

const featureDoc = "# dk:\n# type: test-case\n# title: Bộ lọc\n# source: F-001\n\n@F-001\nFeature: Bộ lọc\n"

func TestFeatureCommentRoundTrip(t *testing.T) {
	fm, body, ok := SplitFile("x.feature", []byte(featureDoc))
	if !ok || GetString(fm, "source") != "F-001" || string(body) != "\n@F-001\nFeature: Bộ lọc\n" {
		t.Fatalf("tách sai: ok=%v body=%q", ok, body)
	}
	out, err := JoinFile("x.feature", fm, body)
	if err != nil || string(out) != featureDoc {
		t.Fatalf("round-trip đổi nội dung (%v):\n%s", err, out)
	}
	// Thân bắt đầu bằng chú thích # phải được cách một dòng trống để không bị nuốt.
	out, _ = JoinFile("x.feature", fm, []byte("# ghi chú\nFeature: x\n"))
	fm2, body2, ok := SplitFile("x.feature", out)
	if !ok || GetString(fm2, "title") != "Bộ lọc" || !bytes.HasPrefix(body2, []byte("\n# ghi chú")) {
		t.Fatalf("thân bắt đầu bằng #:\n%s", out)
	}
	if _, _, ok := SplitFile("x.feature", []byte("Feature: không metadata\n")); ok {
		t.Fatal("thiếu # dk: phải ok=false")
	}
	if fm, body, ok := SplitFile("x.feature", []byte("# dk:\n# type: x")); !ok || GetString(fm, "type") != "x" || len(body) != 0 {
		t.Fatal("file chỉ có metadata phải ok=true với thân rỗng")
	}
}

func TestFeatureCommentEmptyBlock(t *testing.T) {
	if _, _, ok := SplitFile("x.feature", []byte("# dk:\nFeature: x\n")); ok {
		t.Fatal("khối # dk: rỗng phải ok=false")
	}
}
