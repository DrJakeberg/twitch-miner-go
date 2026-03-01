package version

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		want    Version
		wantErr bool
	}{
		{"1.2.3", Version{1, 2, 3}, false},
		{"0.1.0", Version{0, 1, 0}, false},
		{"v1.2.3", Version{1, 2, 3}, false},
		{"10.20.30", Version{10, 20, 30}, false},
		{"bad", Version{}, true},
		{"1.2", Version{}, true},
		{"1.2.3.4", Version{}, true},
	}
	for _, tt := range tests {
		got, err := Parse(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"0.1.0", "0.2.0", -1},
	}
	for _, tt := range tests {
		a, _ := Parse(tt.a)
		b, _ := Parse(tt.b)
		got := Compare(a, b)
		if got != tt.want {
			t.Errorf("Compare(%s, %s) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestString(t *testing.T) {
	Number = "1.2.3"
	GitCommit = "abc1234"
	want := "1.2.3 (abc1234)"
	if got := String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	Number = "dev"
	GitCommit = "unknown"
	if got := String(); got != "dev" {
		t.Errorf("String() = %q, want %q", got, "dev")
	}
}
