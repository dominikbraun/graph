package graph

import "testing"

func TestStringHash(t *testing.T) {
	tests := map[string]struct {
		value        string
		expectedHash string
	}{
		"string value": {
			value:        "London",
			expectedHash: "London",
		},
	}

	for name, test := range tests {
		hash := StringHash(test.value)

		if hash != test.expectedHash {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, test.expectedHash, hash)
		}
	}
}

func TestIntHash(t *testing.T) {
	tests := map[string]struct {
		value        int
		expectedHash int
	}{
		"int value": {
			value:        3,
			expectedHash: 3,
		},
	}

	for name, test := range tests {
		hash := IntHash(test.value)

		if hash != test.expectedHash {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, test.expectedHash, hash)
		}
	}
}

func TestInt32Hash(t *testing.T) {
	tests := map[string]struct {
		value        int32
		expectedHash int32
	}{
		"int32 value": {
			value:        3,
			expectedHash: 3,
		},
	}

	for name, test := range tests {
		hash := Int32Hash(test.value)

		if hash != test.expectedHash {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, test.expectedHash, hash)
		}
	}
}

func TestInt64Hash(t *testing.T) {
	tests := map[string]struct {
		value        int64
		expectedHash int64
	}{
		"int64 value": {
			value:        3,
			expectedHash: 3,
		},
	}

	for name, test := range tests {
		hash := Int64Hash(test.value)

		if hash != test.expectedHash {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, test.expectedHash, hash)
		}
	}
}

func TestUint32Hash(t *testing.T) {
	tests := map[string]struct {
		value        uint32
		expectedHash uint32
	}{
		"uint32 value": {
			value:        3,
			expectedHash: 3,
		},
	}

	for name, test := range tests {
		hash := Uint32Hash(test.value)

		if hash != test.expectedHash {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, test.expectedHash, hash)
		}
	}
}

func TestUint64Hash(t *testing.T) {
	tests := map[string]struct {
		value        uint64
		expectedHash uint64
	}{
		"uint64 value": {
			value:        3,
			expectedHash: 3,
		},
	}

	for name, test := range tests {
		hash := Uint64Hash(test.value)

		if hash != test.expectedHash {
			t.Errorf("%s: hash expectancy doesn't match: expected %v, got %v", name, test.expectedHash, hash)
		}
	}
}
