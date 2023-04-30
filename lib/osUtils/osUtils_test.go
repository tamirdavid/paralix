package osUtils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFilesContentAndWriteToOneFile(t *testing.T) {
	type args struct {
		output *os.File
		files  []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test With Existing Files",
			args: args{
				output: createTempFile(t),
				files:  []string{"testdata/file1.txt", "testdata/file2.txt"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create testdata directory if it doesn't exist
			if _, err := os.Stat("testdata"); os.IsNotExist(err) {
				os.Mkdir("testdata", os.ModePerm)
			}

			if tt.name == "Test With Existing Files" {
				err := ioutil.WriteFile(filepath.Join("testdata", "file1.txt"), []byte("content of file1"), os.ModePerm)
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile(filepath.Join("testdata", "file2.txt"), []byte("content of file2"), os.ModePerm)
				if err != nil {
					panic(err)
				}
			}
			readError := ReadFilesContentAndWriteToOneFile(tt.args.output, tt.args.files)
			if (readError != nil) != tt.wantErr {
				t.Errorf("ReadFilesContentAndWriteToOneFile() error = %v, wantErr %v", readError, tt.wantErr)
			}

			// Check output file content
			outputContent, err := ioutil.ReadFile(tt.args.output.Name())
			if err != nil {
				t.Errorf("error reading output file content: %v", err)
			}
			expectedOutput := "file1.txt\ncontent of file1\nfile2.txt\ncontent of file2\n"
			if string(outputContent) != expectedOutput {
				t.Errorf("output file content = %s, want %s", string(outputContent), expectedOutput)
			}
		})
	}
}

func createTempFile(t *testing.T) *os.File {
	file, err := ioutil.TempFile("", "test_output")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return file
}

func TestPrintFileContent(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		fileData string
		wantErr  bool
	}{
		{
			name:     "valid file",
			filePath: "testdata/valid_file.txt",
			fileData: "This is some test data.",
			wantErr:  false,
		},
		{
			name:     "invalid file",
			filePath: "testdata/invalid_file.txt",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test file
			if tc.fileData != "" {
				err := ioutil.WriteFile(tc.filePath, []byte(tc.fileData), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				defer os.Remove(tc.filePath)
			}

			// Call function being tested
			err := PrintFileContent(tc.filePath)

			// Check if error matches expectation
			if (err != nil) != tc.wantErr {
				t.Errorf("PrintFileContent() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
func TestMakeDirectoryDeleteIfExists(t *testing.T) {
	type args struct {
		dirName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create directory that doesn't exist",
			args:    args{dirName: "testdata/mydir"},
			wantErr: false,
		},
		{
			name:    "Delete and recreate directory that exists",
			args:    args{dirName: "testdata/mydir"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Delete the directory before running the test
			os.RemoveAll(tt.args.dirName)

			// Create the directory before running the test if needed
			if _, err := os.Stat(tt.args.dirName); os.IsNotExist(err) {
				if err := os.MkdirAll(tt.args.dirName, 0755); err != nil {
					t.Fatalf("error creating directory: %v", err)
				}
			}

			// Run the function and check if it returns an error if not expected
			if err := MakeDirectoryDeleteIfExists(tt.args.dirName); (err != nil) != tt.wantErr {
				t.Errorf("MakeDirectoryDeleteIfExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		{
			name:    "Test valid file path",
			args:    args{filePath: "./testfile.txt"},
			want:    &os.File{},
			wantErr: false,
		},
		{
			name:    "Test no valid file path",
			args:    args{filePath: "./novalid/testfile.txt"},
			want:    &os.File{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
