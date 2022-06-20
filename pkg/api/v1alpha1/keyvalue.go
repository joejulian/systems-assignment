package v1alpha1

import (
	"bufio"
	"bytes"
	"os"

	"github.com/google/uuid"
	"github.com/savingoyal/systems-assignment/pkg/flatfile"
)

type KeyValue struct {
	Key   *uuid.UUID `flatfile:"1" json:"key"`
	Value string     `flatfile:"2," json:"value"`
}

type KeyValueStore struct {
	Data  []KeyValue
	index map[string]*KeyValue
}

func (kvs *KeyValueStore) UnmarshalBuffer(data bytes.Buffer) error {
	scanner := bufio.NewScanner(&data)
	return kvs.Unmarshal(scanner)
}

func (kvs *KeyValueStore) Unmarshal(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		kv := KeyValue{}
		if err := flatfile.Unmarshal([]byte(line), &kv); err != nil {
			return err
		}
		if kv.Key.String() == "" {
			return ErrInvalidKey
		}
		if kv.Value == "" {
			return ErrInvalidValue
		}
		if _, ok := kvs.index[kv.Key.String()]; ok {
			return ErrDuplicateKey
		}
		kvs.Data = append(kvs.Data, kv)
		if kvs.index == nil {
			kvs.index = make(map[string]*KeyValue)
		}
		kvs.index[kv.Key.String()] = &kv
	}

	return nil
}

func (kvs *KeyValueStore) Get(key string) (*KeyValue, error) {
	if kvs.index == nil {
		return nil, ErrCacheEmpty
	}

	kv, ok := kvs.index[key]
	if !ok {
		return nil, ErrInvalidKey
	}
	return kv, nil
}

func LoadKeyValueStore(filename string) (*KeyValueStore, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	kvs := &KeyValueStore{}
	scanner := bufio.NewScanner(f)
	if err := kvs.Unmarshal(scanner); err != nil {
		return nil, err
	}
	return kvs, nil
}
