package parser

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testdata    = "testdata"
	goldenFile1 = "report1.golden.json"
	goldenFile2 = "report2.golden.json"
	testReport1 = "report1.txt"
	testReport2 = "report2.txt"
)

func getReaderFromFile(tb testing.TB, fPath string) io.Reader {
	tb.Helper()

	file, err := os.Open(fPath)
	require.NoError(tb, err)
	return file
}

func getReportFromFile(tb testing.TB, fPath string) Report {
	tb.Helper()

	b, err := ioutil.ReadFile(fPath)
	require.NoError(tb, err)

	var rp Report
	err = json.Unmarshal(b, &rp)
	require.NoError(tb, err)
	return rp
}

func TestProcessReport(t *testing.T) {
	type args struct {
		filepath string
	}
	type expected struct {
		filepath string
		wantErr  bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "process report type 1",
			args: args{
				filepath: filepath.Join(testdata, testReport1),
			},
			expected: expected{
				filepath: filepath.Join(testdata, goldenFile1),
				wantErr:  false,
			},
		},
		{
			name: "process report type 2",
			args: args{
				filepath: filepath.Join(testdata, testReport2),
			},
			expected: expected{
				filepath: filepath.Join(testdata, goldenFile2),
				wantErr:  false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			data := getReaderFromFile(t, tt.args.filepath)
			got, err := ProcessReport(data)
			if tt.expected.wantErr {
				assert.Error(t, err)
				return
			}

			wantReport := getReportFromFile(t, tt.expected.filepath)
			assert.NoError(t, err)
			assert.Equal(t, wantReport, got)
		})
	}
}

func Test_processReport(t *testing.T) {
	type args struct {
		filepath string
		rt       reportType
	}
	type expected struct {
		filepath string
		wantErr  bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "process report type 1",
			args: args{
				filepath: filepath.Join(testdata, testReport1),
				rt:       reportType1,
			},
			expected: expected{
				filepath: filepath.Join(testdata, goldenFile1),
				wantErr:  false,
			},
		},
		{
			name: "process report type 2",
			args: args{
				filepath: filepath.Join(testdata, testReport2),
				rt:       reportType2,
			},
			expected: expected{
				filepath: filepath.Join(testdata, goldenFile2),
				wantErr:  false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			data := getReaderFromFile(t, tt.args.filepath)
			got, err := processReport(data, tt.args.rt)
			if tt.expected.wantErr {
				assert.Error(t, err)
				return
			}
			wantReport := getReportFromFile(t, tt.expected.filepath)
			assert.NoError(t, err)
			assert.Equal(t, wantReport, got)
		})
	}
}

func Test_getReportType(t *testing.T) {
	type args struct {
		filepath string
	}
	type expected struct {
		rt reportType
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "report type 1",
			args: args{
				filepath: filepath.Join(testdata, testReport1),
			},
			expected: expected{
				rt: reportType1,
			},
		},
		{
			name: "report type 2",
			args: args{
				filepath: filepath.Join(testdata, testReport2),
			},
			expected: expected{
				rt: reportType2,
			},
		},
		{
			name: "unknown type",
			args: args{
				filepath: filepath.Join(testdata, "empty.json"),
			},
			expected: expected{
				rt: reportTypeUnknown,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			data := getReaderFromFile(t, tt.args.filepath)
			got := getReportType(data)
			assert.Equal(t, tt.expected.rt.String(), got.String())
		})
	}
}