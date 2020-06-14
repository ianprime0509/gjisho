package tatoeba

import "testing"

func TestParseIndex(t *testing.T) {
	tests := []struct {
		input string
		index Index
	}{
		{"は", Index{Word: "は"}},
		{"直ぐに{すぐに}", Index{Word: "直ぐに", SentenceForm: "すぐに"}},
		{"為る(する){する}", Index{Word: "為る", Disambiguation: "する", SentenceForm: "する"}},
		{"出来る[01]{できない}~", Index{Word: "出来る", Sense: "01", SentenceForm: "できない", Good: true}},
		{"申し訳ございません[01]", Index{Word: "申し訳ございません", Sense: "01"}},
		{"口(くち)[05]", Index{Word: "口", Disambiguation: "くち", Sense: "05"}},
	}

	for _, test := range tests {
		parsed, err := parseIndex(test.input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if parsed != test.index {
			t.Errorf("indices not equal: want %v, got %v", test.index, parsed)
		}
	}
}
