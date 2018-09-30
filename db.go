package main

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"strconv"
)

const (
	// KLastHeight key last height
	KLastHeight = "OBS_Last_height"
	// KTaskPrefix task list item prefix
	KTaskPrefix = "OBS_Task_"
)

// Database leveldb database
type Database struct {
	db *leveldb.DB
}

// Close close source
func (s *Database) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

//GetLastHeight get last height
func (s *Database) GetLastHeight() (int, error) {
	value, err := s.db.Get([]byte(KLastHeight), nil)
	if err == leveldb.ErrNotFound {
		err = s.UpdataLastHeight(0)
		if err != nil {
			return -1, err
		}
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(value))
}

// UpdataLastHeight updata lastheight
func (s *Database) UpdataLastHeight(height int) error {
	return s.db.Put([]byte(KLastHeight), []byte(strconv.Itoa(height)), nil)
}

// GetTaskList 获取任务列表
func (s *Database) GetTaskList() ([]DepositTask, error) {
	iter := s.db.NewIterator(util.BytesPrefix([]byte(KTaskPrefix)), nil)
	result := make([]DepositTask, 0)
	for iter.Next() {
		task := DepositTask{}
		if err := json.Unmarshal(iter.Value(), &task); err != nil {
			return nil, err
		}
		result = append(result, task)
	}
	iter.Release()
	return result, iter.Error()
}

// DeleteTask 删除任务
func (s *Database) DeleteTask(task DepositTask) error {
	key := KTaskPrefix + task.BlockHash
	return s.db.Delete([]byte(key), nil)
}

// AddTask 添加任务
func (s *Database) AddTask(task DepositTask) error {
	key := KTaskPrefix + task.BlockHash
	value, err := json.Marshal(&task)
	if err != nil {
		return err
	}
	return s.db.Put([]byte(key), value, nil)
}

// Exist 任务是否存在
func (s *Database) Exist(task DepositTask) (bool, error) {
	key := KTaskPrefix + task.BlockHash
	return s.db.Has([]byte(key), nil)
}

// NewDatabase new database
func NewDatabase(path string) (*Database, error) {
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	db := &Database{
		db: ldb,
	}
	return db, nil
}
