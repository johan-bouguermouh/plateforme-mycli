package model

import (
	"encoding/json"
	"errors"

	color "bucketool/utils/colorPrint"

	bolt "go.etcd.io/bbolt"
)

type Alias struct {
	ID 	  int
    HOST        string
	Name       string
	Port 	 int
	KeyName   string
	SecretKey string
	Current    bool
}

type AliasStore struct {
    db *bolt.DB
}

func UseAliasStore(db *bolt.DB) *AliasStore {
    store := &AliasStore{db: db}
    store.init()
    return store
}

func (store *AliasStore) init() {
    store.db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte("Alias"))
        return err
    })
}

func (store *AliasStore) SaveAlias(Alias *Alias) error {
	var suffixMessage string = ""
	if(Alias.Current){
		err := store.UnsetCurrentAlias()
		if err != nil {
			return err
		}
		suffixMessage = " and set as current"
	}
    return store.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Alias"))
        encoded, err := json.Marshal(Alias)
        if err != nil {
            return err
        }
		println(color.ColorPrint("Grey",
		"Alias has "+Alias.Name+" been saved" + suffixMessage,
		&color.Options{
			Italic: true,
		}))
        return b.Put([]byte(Alias.Name), encoded)
    })
}

func (store *AliasStore) ReadAlias(Name string) (Alias, error) {
    var Alias Alias
    err := store.db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Alias"))
        v := b.Get([]byte(Name))
        if v == nil {
            return nil
        }
        return json.Unmarshal(v, &Alias)
    })
    return Alias, err
}

func (store *AliasStore) UpdateAlias(Alias Alias) error {
    return store.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Alias"))
        encoded, err := json.Marshal(Alias)
        if err != nil {
            return err
        }

        return b.Put([]byte(Alias.Name), encoded)
    })
}

func (store *AliasStore) DeleteAlias(id int) error {
    return store.db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Alias"))
        return b.Delete(itob(id))
    })
}

func (store *AliasStore) ListAliass() ([]Alias, error) {
	var Aliass []Alias
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Alias"))
		return b.ForEach(func(k, v []byte) error {
			var Alias Alias
			if err := json.Unmarshal(v, &Alias); err != nil {
				return err
			}
			Aliass = append(Aliass, Alias)
			return nil
		})
	})
	return Aliass, err
}

func (store *AliasStore) DeleteAliasByName(Name string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Alias"))
		return b.Delete([]byte(Name))
	})
}

func (store *AliasStore) DeleteAllAlias() error {
	return store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Alias"))
		return b.ForEach(func(k, v []byte) error {
			return b.Delete(k)
		})
	})
}

func (store *AliasStore) GetCurrentAlias() (Alias, error){
	var currentAlias Alias
	found := false
	err := store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Alias"))
		return bucket.ForEach(func(key, value []byte) error {
			var Alias Alias
			if err := json.Unmarshal(value, &Alias); err != nil {
				return err
			}
			if Alias.Current {
				currentAlias  = Alias
				found = true
                return nil

			}
			return nil
		})
	})
	// On verife que l'Alias ne comporte qu'un seul element
	if err != nil {
        return Alias{}, err
    }
	if !found {
		return Alias{}, errors.New(color.RedP("no Alias is marked as current"))
	}

	return currentAlias, err
}

func (store *AliasStore) UnsetCurrentAlias() error {
	return store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Alias"))
		return bucket.ForEach(func(key, value []byte) error {
			var Alias Alias
			if err := json.Unmarshal(value, &Alias); err != nil {
				return err
			}
			if Alias.Current {
				Alias.Current = false
				println(color.ColorPrint("Grey",
				"Current Alias has been unset :",
				&color.Options{
					Italic: true,
				}), 
				color.ColorPrint("Grey",
				Alias.Name,
				&color.Options{
					Italic: true,
					Strikethrough: true,
				}))
			}
			encoded, err := json.Marshal(Alias)
			if err != nil {
				return err
			}
			return bucket.Put(key, encoded)
		})
	})
}

func (store *AliasStore) SetCurrentAlias(Name string) error {
	err := store.UnsetCurrentAlias()
	if err != nil {
		return err
	}

	return store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Alias"))
		return bucket.ForEach(func(key, value []byte) error {
			var Alias Alias
			if err := json.Unmarshal(value, &Alias); err != nil {
				return err
			}
			if Alias.Name == Name {
				Alias.Current = true
				encoded, err := json.Marshal(Alias)
				if err != nil {
					return err
				}
				println(color.ColorPrint("Grey",
				"Switch Alias to " + color.GreenP(Name),
				&color.Options{
					Italic: true,
				}))
				return bucket.Put(key, encoded)

			}
			return nil
		})
	})
}

func (store *AliasStore) IsAliasExist(Name string) bool {
	suspectedAlias, err := store.ReadAlias(Name)
	if(err != nil){
		return false
	}
	return !store.IsEmptyAlias(suspectedAlias)
}

func (store *AliasStore) IsEmptyAlias(alias Alias) bool {
    return alias == Alias{}
}
	