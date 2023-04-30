package paralixutils

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestRunCmdAndWaitForItToFinish(t *testing.T) {
	type args struct {
		cmd *exec.Cmd
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "successful command without output",
			args:    args{cmd: exec.Command("echo", "test")},
			wantErr: false,
		},
		{
			name:    "successful command with output",
			args:    args{cmd: exec.Command("ls", "-la")},
			wantErr: false,
		},
		{
			name:    "failing command with output",
			args:    args{cmd: exec.Command("ls", "nonexistentfile")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunCmdAndWaitForItToFinish(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("RunCmdAndWaitForItToFinish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetValuesBetweenDelimiters(t *testing.T) {
	type args struct {
		str       string
		openChar  string
		closeChar string
		splitChar string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "simple test case",
			args: args{
				str:       "Hello [world,how]",
				openChar:  "[",
				closeChar: "]",
				splitChar: ",",
			},
			want:    []string{"world", "how"},
			wantErr: false,
		},
		{
			name: "missing close delimiter",
			args: args{
				str:       "Hello [world",
				openChar:  "[",
				closeChar: "]",
				splitChar: ",",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValuesBetweenDelimiters(tt.args.str, tt.args.openChar, tt.args.closeChar, tt.args.splitChar)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValuesBetweenDelimiters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValuesBetweenDelimiters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadLinesFromFileReturnSliceOfLines(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "read file with multiple lines",
			args: args{
				filepath: "testfile.txt",
			},
			want:    []string{"line 1", "line 2", "line 3"},
			wantErr: false,
		},
		{
			name: "read non-existent file",
			args: args{
				filepath: "testdata/nonexistent.txt",
			},
			wantErr: true,
		},
	}

	// Create test file
	file, err := os.Create("testfile.txt")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer file.Close()

	// Write test data to file
	testData := "line 1\nline 2\nline 3\n"
	_, err = file.WriteString(testData)
	if err != nil {
		t.Fatalf("failed to write test data to file: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadLinesFromFileReturnSliceOfLines(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLinesFromFileReturnSliceOfLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadLinesFromFileReturnSliceOfLines() = %v, want %v", got, tt.want)
			}
		})
	}

	// Remove test file
	err = os.Remove("testfile.txt")
	if err != nil {
		t.Fatalf("failed to remove test file: %v", err)
	}
}

func TestIsStringInSlice(t *testing.T) {
	type args struct {
		s      []string
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success case",
			args: args{
				s:      []string{"apple", "banana", "cherry", "date", "elderberry"},
				target: "banana",
			},
			want: true,
		},
		{
			name: "fail case",
			args: args{
				s:      []string{"apple", "banana", "cherry", "date", "elderberry"},
				target: "fig",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStringInSlice(tt.args.s, tt.args.target); got != tt.want {
				t.Errorf("IsStringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMatchedRegexOccurencesFromString(t *testing.T) {
	type args struct {
		regex string
		input string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "successful match",
			args: args{
				regex: `<(.*?)>`,
				input: "My social security number is <123-45-6789>",
			},
			want: []string{"123-45-6789"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMatchedRegexOccurencesFromString(tt.args.regex, tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatchedRegexOccurencesFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
