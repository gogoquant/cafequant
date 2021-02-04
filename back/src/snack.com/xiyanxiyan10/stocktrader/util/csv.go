package util

import (
	"encoding/csv"
	"os"
)

// ReadCSV 读取csv文件
func ReadCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	w := csv.NewReader(f)
	data, err := w.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteCSV 写入csv文件
func WriteCSV(path string, vec [][]string, add bool) error {
	var file *os.File
	var err error
	if add {
		file, err = os.OpenFile(path, os.O_APPEND, 0644)
	} else {
		file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入UTF-8 BOM，防止中文乱码
	_, err = file.WriteString("\xEF\xBB\xBF")
	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	for _, v := range vec {
		err = w.Write(v)
		if err != nil {
			return err
		}
	}
	// 写文件需要flush，不然缓存满了，后面的就写不进去了，只会写一部分
	w.Flush()
	return nil
}
