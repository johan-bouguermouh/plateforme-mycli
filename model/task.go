package model

import (
	"encoding/binary"
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

type Task struct {
    ID        int
    Name      string
    Completed bool
}

type TaskStore struct {
    db *bolt.DB
}

func UseTaskStore(db *bolt.DB) *TaskStore {
    store := &TaskStore{db: db}
    store.init()
    return store
}

func (store *TaskStore) init() {
    store.db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte("Task"))
        return err
    })
}

func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}

func (store *TaskStore) CreateTask(task Task) error {
    return store.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Task"))
        id, _ := b.NextSequence()
        task.ID = int(id)
        encoded, err := json.Marshal(task)
        if err != nil {
            return err
        }
        return b.Put(itob(task.ID), encoded)
    })
}

func (store *TaskStore) ReadTask(id int) (Task, error) {
    var task Task
    err := store.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Task"))
        v := b.Get(itob(id))
        if v == nil {
            return nil
        }
        return json.Unmarshal(v, &task)
    })
    return task, err
}

func (store *TaskStore) UpdateTask(task Task) error {
    return store.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Task"))
        encoded, err := json.Marshal(task)
        if err != nil {
            return err
        }
        return b.Put(itob(task.ID), encoded)
    })
}

func (store *TaskStore) DeleteTask(id int) error {
    return store.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Task"))
        return b.Delete(itob(id))
    })
}

func (store *TaskStore) ListTasks() ([]Task, error) {
	var tasks []Task
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Task"))
		return b.ForEach(func(k, v []byte) error {
			var task Task
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			tasks = append(tasks, task)
			return nil
		})
	})
	return tasks, err
}