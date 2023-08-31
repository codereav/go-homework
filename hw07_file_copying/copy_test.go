package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testdataDir    = "./testdata/"
	inputFileName  = "input.txt"
	outFileNameTpl = "out_offset%d_limit%d.txt"
)

type testCase struct {
	offset        int64
	limit         int64
	error         error
	inputFileName string `default:"input.txt"`
}

func TestCopy(t *testing.T) {
	tests := []testCase{
		{offset: 0, limit: 0},
		{offset: 0, limit: 1000},
		{offset: 0, limit: 10000},
		{offset: 100, limit: 1000},
		{offset: 6000, limit: 1000},
		{offset: 10000, limit: 1000, error: ErrOffsetExceedsFileSize},
		{inputFileName: "file not exists", error: ErrUnsupportedFile},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("limit: %d, offset: %d", tc.limit, tc.offset), func(t *testing.T) {
			outFileName := fmt.Sprintf(outFileNameTpl, tc.offset, tc.limit)
			outFileRealPath := "/tmp/" + outFileName

			var inputFile string

			if tc.inputFileName != "" {
				inputFile = tc.inputFileName
			} else {
				inputFile = inputFileName
			}

			err := Copy(testdataDir+inputFile, outFileRealPath, tc.offset, tc.limit)

			if tc.error != nil {
				require.ErrorIs(t, err, tc.error)
				return
			}

			require.NoError(t, err)

			actualFile, err := os.Open(outFileRealPath)
			require.NoError(t, err)
			defer func() {
				err := actualFile.Close()
				require.NoError(t, err)
			}()

			expectedFile, err := os.Open(testdataDir + outFileName)
			require.NoError(t, err)
			defer func() {
				err := expectedFile.Close()
				require.NoError(t, err)
			}()

			_, err = compareFiles(t, expectedFile, actualFile)
			require.NoError(t, err)
		})
	}
}

func compareFiles(t *testing.T, file1, file2 *os.File) (bool, error) {
	t.Helper()
	// Кол-во байт, которые будем сравнивать за одну итерацию
	const bufferSize = 4096

	buffer1 := make([]byte, bufferSize)
	buffer2 := make([]byte, bufferSize)

	for {
		// читаем данные из файлов в буферы
		_, err1 := file1.Read(buffer1)
		_, err2 := file2.Read(buffer2)

		// ловим любые ошибки, кроме EOF
		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}

		// Сравниваем буферы
		require.Equal(t, buffer1, buffer2)

		// Если оба файла закончились - выходим из цикла
		if err1 == io.EOF && err2 == io.EOF {
			break
		}
	}

	return true, nil
}
