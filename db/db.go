// SPDX-License-Identifier: MIT

// Package db 提供区域数据库的相关操作
package db

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/issue9/cnregion/id"
	"github.com/issue9/errwrap"
)

// DB 区域数据库信息
type DB struct {
	*Region
	Versions []int // 支持的版本
}

// Load 从数据库文件加载数据
func Load(file string) (*DB, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return Unmarshal(data)
}

// Dump 输出到文件
func (db *DB) Dump(file string) error {
	data, err := Marshal(db)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, os.ModePerm)
}

// Marshal 将 DB 转换成 []byte
func Marshal(db *DB) ([]byte, error) {
	return db.marshal()
}

// Unmarshal 解码 data 至 DB
func Unmarshal(data []byte) (*DB, error) {
	db := &DB{}
	if err := db.unmarshal(data); err != nil {
		return nil, err
	}
	return db, nil
}

// Find 查找指定 ID 对应的信息
func (db *DB) Find(regionID string) *Region {
	province, city, county, town, village := id.Split(regionID)
	return db.findItem(province, city, county, town, village)
}

func (db *DB) marshal() ([]byte, error) {
	vers := make([]string, 0, len(db.Versions))
	for _, v := range db.Versions {
		vers = append(vers, strconv.Itoa(v))
	}

	buf := errwrap.Buffer{Buffer: bytes.Buffer{}}
	buf.WByte('[')
	buf.WriteString(strings.Join(vers, ","))
	buf.WByte(']').WByte(':')

	err := db.Region.marshal(&buf)
	if err != nil {
		return nil, err
	}

	if buf.Err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (db *DB) unmarshal(data []byte) error {
	index := bytes.IndexByte(data, ':')
	arr := bytes.Trim(data[:index], "[]")
	arrs := strings.Split(string(arr), ",")
	db.Versions = make([]int, 0, len(arrs))
	for _, item := range arrs {
		v, err := strconv.Atoi(item)
		if err != nil {
			return err
		}
		db.Versions = append(db.Versions, v)
	}

	data = data[index+1:]
	db.Region = &Region{}
	return db.Region.unmarshal(data)
}