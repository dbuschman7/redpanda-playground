package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	var seq Seq[string, int]

	seq = Seq[string, int]{}
	assert.Equal(t, "", seq.First)
	assert.Equal(t, 0, seq.Second)

	seq = Seq[string, int]{First: "42"}
	assert.Equal(t, "42", seq.First)
	assert.Equal(t, 0, seq.Second)

	seq = Seq[string, int]{Second: 42}
	assert.Equal(t, "", seq.First)
	assert.Equal(t, 42, seq.Second)

}

func TestSequenceSequence(t *testing.T) {
	var seq Seq[Seq[string, int], bool]

	seq = Seq[Seq[string, int], bool]{}
	assert.Equal(t, "", seq.First.First)
	assert.Equal(t, 0, seq.First.Second)
	assert.Equal(t, false, seq.Second)

	seq = Seq[Seq[string, int], bool]{
		First: Seq[string, int]{First: "42"},
	}
	assert.Equal(t, "42", seq.First.First)
	assert.Equal(t, 0, seq.First.Second)
	assert.Equal(t, false, seq.Second)

	seq = Seq[Seq[string, int], bool]{
		First:  Seq[string, int]{Second: 42},
		Second: true,
	}
	assert.Equal(t, "", seq.First.First)
	assert.Equal(t, 42, seq.First.Second)
	assert.Equal(t, true, seq.Second)
}

func TestSequenceSequenceSequence(t *testing.T) {
	var seq Seq[Seq[Seq[string, int], bool], string]

	seq = Seq[Seq[Seq[string, int], bool], string]{}
	assert.Equal(t, "", seq.First.First.First)
	assert.Equal(t, 0, seq.First.First.Second)
	assert.Equal(t, false, seq.First.Second)
	assert.Equal(t, "", seq.Second)

	seq = Seq[Seq[Seq[string, int], bool], string]{
		First: Seq[Seq[string, int], bool]{
			First: Seq[string, int]{First: "42"},
		},
	}
	assert.Equal(t, "42", seq.First.First.First)
	assert.Equal(t, 0, seq.First.First.Second)
	assert.Equal(t, false, seq.First.Second)
	assert.Equal(t, "", seq.Second)

	seq = Seq[Seq[Seq[string, int], bool], string]{
		First: Seq[Seq[string, int], bool]{
			First:  Seq[string, int]{Second: 42},
			Second: true,
		},
		Second: "hello",
	}
	assert.Equal(t, "", seq.First.First.First)
	assert.Equal(t, 42, seq.First.First.Second)
	assert.Equal(t, true, seq.First.Second)
	assert.Equal(t, "hello", seq.Second)
}

func TestSequenceSequenceSequenceSequence(t *testing.T) {
	var seq Seq[Seq[Seq[Seq[string, int], bool], string], float64]

	seq = Seq[Seq[Seq[Seq[string, int], bool], string], float64]{}
	assert.Equal(t, "", seq.First.First.First.First)
	assert.Equal(t, 0, seq.First.First.First.Second)
	assert.Equal(t, false, seq.First.First.Second)
	assert.Equal(t, "", seq.First.Second)
	assert.Equal(t, 0.0, seq.Second)

	seq = Seq[Seq[Seq[Seq[string, int], bool], string], float64]{
		First: Seq[Seq[Seq[string, int], bool], string]{
			First: Seq[Seq[string, int], bool]{
				First: Seq[string, int]{First: "42"},
			},
		},
	}
	assert.Equal(t, "42", seq.First.First.First.First)
	assert.Equal(t, 0, seq.First.First.First.Second)
	assert.Equal(t, false, seq.First.First.Second)
	assert.Equal(t, "", seq.First.Second)
	assert.Equal(t, 0.0, seq.Second)

	seq = Seq[Seq[Seq[Seq[string, int], bool], string], float64]{
		First: Seq[Seq[Seq[string, int], bool], string]{
			First: Seq[Seq[string, int], bool]{
				First:  Seq[string, int]{Second: 42},
				Second: true,
			},
			Second: "hello",
		},
		Second: 3.14,
	}
	assert.Equal(t, "", seq.First.First.First.First)
	assert.Equal(t, 42, seq.First.First.First.Second)
	assert.Equal(t, true, seq.First.First.Second)
	assert.Equal(t, "hello", seq.First.Second)
	assert.Equal(t, 3.14, seq.Second)
}

func TestApply(t *testing.T) {
	parser := Apply(
		Succeed(Seq[Empty, int]{Second: 42}),
		func(int) string { return "hello" },
	)

	result, _, err := parser(State{})
	assert.Nil(t, err)
	assert.Equal(t, "hello", result)
}

func TestApply2(t *testing.T) {
	parser := Apply2(
		Succeed(Seq[Seq[Empty, int], string]{
			First:  Seq[Empty, int]{Second: 42},
			Second: "world",
		}),
		func(int, string) string { return "hello" },
	)

	result, _, err := parser(State{})
	assert.Nil(t, err)
	assert.Equal(t, "hello", result)
}

func TestApply3(t *testing.T) {
	parser := Apply3(
		Succeed(Seq[Seq[Seq[Empty, int], string], float64]{
			First: Seq[Seq[Empty, int], string]{
				Second: "hello",
			},
			Second: 3.14,
		}),
		func(int, string, float64) string { return "hello" },
	)

	result, _, err := parser(State{})
	assert.Nil(t, err)
	assert.Equal(t, "hello", result)

}

func TestApply4(t *testing.T) {
	parser := Apply4(
		Succeed(Seq[Seq[Seq[Seq[Empty, int], string], float64], string]{
			Second: "hello",
		}),
		func(int, string, float64, string) string { return "hello" },
	)

	result, _, err := parser(State{})
	assert.Nil(t, err)
	assert.Equal(t, "hello", result)
}

func TestAppendSkipping(t *testing.T) {
	parser := AppendSkipping(
		Succeed(Seq[Empty, int]{Second: 42}),
		Succeed(Seq[Empty, bool]{Second: true}),
	)

	result, _, err := parser(State{})
	assert.Nil(t, err)
	assert.Equal(t, 42, result.Second)
}

func TestAppendSkippingFail(t *testing.T) {
	parser := AppendSkipping(
		Fail[int],
		Succeed(Seq[Empty, bool]{Second: true}),
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestAppendSkippingFailSecond(t *testing.T) {
	parser := AppendSkipping(
		Succeed(Seq[Empty, int]{Second: 42}),
		Fail[bool],
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestAppendSkippingFailBoth(t *testing.T) {
	parser := AppendSkipping(
		Fail[int],
		Fail[bool],
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestAppendSkippingFailBothSecond(t *testing.T) {
	parser := AppendSkipping(
		Succeed(Seq[Empty, int]{Second: 42}),
		Fail[bool],
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestAppendSkippingFailBothFirst(t *testing.T) {
	parser := AppendSkipping(
		Fail[int],
		Succeed(Seq[Empty, bool]{Second: true}),
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestApplyFail(t *testing.T) {
	parser := Apply(
		Fail[Seq[Empty, int]],
		func(int) string { return "hello" },
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestApply2Fail(t *testing.T) {
	parser := Apply2(
		Fail[Seq[Seq[Empty, int], string]],
		func(int, string) string { return "hello" },
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestApply3Fail(t *testing.T) {
	parser := Apply3(
		Fail[Seq[Seq[Seq[Empty, int], string], float64]],
		func(int, string, float64) string { return "hello" },
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestSequenceFail(t *testing.T) {
	var seq = Seq[string, int]{First: "42", Second: 42}
	assert.Equal(t, "42", seq.First)
	assert.Equal(t, 42, seq.Second)
}

func TestSequenceSequenceFail(t *testing.T) {
	var seq = Seq[Seq[string, int], bool]{
		First:  Seq[string, int]{Second: 42},
		Second: true,
	}
	assert.Equal(t, "", seq.First.First)
	assert.Equal(t, 42, seq.First.Second)
	assert.Equal(t, true, seq.Second)
}

func TestStartKeepingFail(t *testing.T) {
	parser := StartKeeping(
		Fail[Seq[Empty, int]],
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestStartKeepingFailBoth(t *testing.T) {
	parser := StartKeeping(
		Fail[Seq[Empty, int]],
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestStartKeepingFailBothFirst(t *testing.T) {
	parser := StartKeeping(
		Fail[Seq[Empty, int]],
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}

func TestStartKeepingFailBothFirstSecond(t *testing.T) {
	parser := StartKeeping(
		Fail[Seq[Empty, int]],
	)

	_, _, err := parser(State{})
	assert.NotNil(t, err)
}
