package insertremove

import "testing"

var implementations = []struct {
	descr string
	f     func(n int) error
}{
	{
		descr: "slice",
		f:     InsertRemoveSlice,
	},
	{
		descr: "list",
		f:     InsertRemoveList,
	},
}

func TestInsertRemove(t *testing.T) {
	testCases := []struct {
		descr       string
		input       int
		expectError bool
	}{
		{
			descr:       "negative value",
			input:       -5,
			expectError: true,
		},
		{
			descr:       "single element",
			input:       1,
			expectError: false,
		},
		{
			descr:       "small number of elements",
			input:       10,
			expectError: false,
		},
		{
			descr:       "large number of elements",
			input:       50000,
			expectError: false,
		},
	}

	for _, impl := range implementations {
		t.Run(impl.descr, func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.descr, func(t *testing.T) {
					err := impl.f(tc.input)
					if err == nil && tc.expectError {
						t.Fatalf("InsertRemove(%d) expected to fail, but returned without error", tc.input)
					}
					if err != nil && !tc.expectError {
						t.Fatalf("InsertRemove(%d) returned unexpected error", tc.input)
					}
				})
			}
		})
	}
}
