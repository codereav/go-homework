package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем исходный файл
	input, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer func() {
		err = input.Close()
		if err != nil {
			fmt.Printf("can't close input file: %s", err)
		}
	}()

	// Пытаемся открыть для записи файл, в который нужно копировать и очищаем его, либо создаём, если не существует
	output, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return fmt.Errorf("can't create output file: %w", err)
	}
	defer func() {
		err = output.Close()
		if err != nil {
			fmt.Printf("can't close output file: %s", err)
		}
	}()

	// Сдвигаем указатель в исходном файле на offset от начала файла
	_, err = input.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("can't seek input file: %w", err)
	}

	// Определяем total
	total, err := detectTotal(limit, offset, input)
	if err != nil {
		return err
	}

	// Конфигурируем и запускаем progress bar
	progress := pb.New(int(total)).
		SetRefreshRate(100 * time.Millisecond).
		SetUnits(pb.U_BYTES).
		Start()

	outputWriter := io.MultiWriter(output, progress)

	_, err = io.CopyN(outputWriter, input, total)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("unable to copy file: %w", err)
	}

	progress.Finish()

	return nil
}

func detectTotal(limit, offset int64, input *os.File) (int64, error) {
	var total, fileSize int64

	// Получаем размер файла
	inputStat, err := input.Stat()
	if err != nil {
		return 0, fmt.Errorf("unable to get file stat: %w", err)
	}
	fileSize = inputStat.Size()
	// offset не может быть больше размера файла
	if offset > fileSize {
		return 0, ErrOffsetExceedsFileSize
	}

	if limit == 0 || offset+limit > fileSize {
		total = fileSize - offset
	} else {
		total = limit
	}

	return total, nil
}
