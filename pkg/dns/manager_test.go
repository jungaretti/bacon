package dns

import "testing"

type diffRecordTestCase struct {
	from         []Record
	to           []Record
	addedCount   int
	removedCount int
}

func TestDiffRecordsLength(t *testing.T) {
	testCases := []diffRecordTestCase{
		{
			from: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
			},
			to: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
			},
			addedCount:   0,
			removedCount: 0,
		},
		{
			from: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
			},
			to: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
				ConfigRecord{
					Name: "mockName2",
				},
			},
			addedCount:   1,
			removedCount: 0,
		},
		{
			from: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
			},
			to: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
				ConfigRecord{
					Name: "mockName2",
				},
				ConfigRecord{
					Name: "mockName3",
				},
			},
			addedCount:   2,
			removedCount: 0,
		},
		{
			from: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
				ConfigRecord{
					Name: "mockName2",
				},
			},
			to: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
			},
			addedCount:   0,
			removedCount: 1,
		},
		{
			from: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
				ConfigRecord{
					Name: "mockName2",
				},
				ConfigRecord{
					Name: "mockName3",
				},
			},
			to: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
			},
			addedCount:   0,
			removedCount: 2,
		},
		{
			from: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
				ConfigRecord{
					Name: "mockName2",
				},
			},
			to: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
				ConfigRecord{
					Name: "mockName3",
				},
			},
			addedCount:   1,
			removedCount: 1,
		},
		{
			from: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
				ConfigRecord{
					Name: "mockName2",
				},
				ConfigRecord{
					Name: "mockName4",
				},
			},
			to: []Record{
				ConfigRecord{
					Name: "mockName1",
				},
				ConfigRecord{
					Name: "mockName3",
				},
				ConfigRecord{
					Name: "mockName5",
				},
			},
			addedCount:   2,
			removedCount: 2,
		},
	}

	for _, testCase := range testCases {
		added, removed := DiffRecords(testCase.from, testCase.to)
		if len(added) != testCase.addedCount {
			t.Error("added is not empty", added)
		}
		if len(removed) != testCase.removedCount {
			t.Error("removed is not empty", removed)
		}
	}
}

func TestDiffRecordsContent(t *testing.T) {
	oldRecord := ConfigRecord{
		Name: "mockName2",
		Type: "A",
	}
	newRecord := ConfigRecord{
		Name: "mockName2",
		Type: "AAAA",
	}

	from := []Record{
		ConfigRecord{
			Name: "mockName1",
		},
		oldRecord,
	}
	to := []Record{
		ConfigRecord{
			Name: "mockName1",
		},
		newRecord,
	}

	added, removed := DiffRecords(from, to)
	if !added[0].Equals(newRecord) {
		t.Error("new record is not equal", added[0], newRecord)
	}
	if !removed[0].Equals(oldRecord) {
		t.Error("old record is not equal", removed[0], oldRecord)
	}
}
