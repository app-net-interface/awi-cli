// Copyright (c) 2023 Cisco Systems, Inc. and its affiliates
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"

	"github.com/app-net-interface/awi-cli/types"
)

const (
	crTableName  = "connection_requests"
	aclTableName = "acls"
)

var (
	tableNames = [...]string{crTableName, aclTableName}
)

type Client interface {
	Open(filename string) error
	Close() error
	UpdateConnectionRequest(cr *types.ConnectionRequest, id string) error
	GetConnectionRequest(id string) (*types.ConnectionRequest, error)
	ListConnectionRequests() ([]types.ConnectionRequest, error)
	DeleteConnectionRequest(id string) error
	UpdateACL(acl *ACL, aclID string) error
	GetACL(aclID string) (*ACL, error)
	ListACLs() ([]ACL, error)
	DeleteACL(aclID string) error
}

type client struct {
	db *bolt.DB
}

func NewClient() Client {
	return &client{}
}

func (client *client) Open(filename string) error {
	options := &bolt.Options{Timeout: time.Second}
	var err error
	client.db, err = bolt.Open(filename, 0600, options)
	if err != nil {
		return err
	}

	return client.db.Update(func(tx *bolt.Tx) error {
		for _, tableName := range tableNames {
			_, err := tx.CreateBucketIfNotExists([]byte(tableName))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (client *client) Close() error {
	return client.db.Close()
}

func (client *client) UpdateConnectionRequest(cr *types.ConnectionRequest, id string) error {
	return update(client, cr, id, crTableName)
}

func (client *client) GetConnectionRequest(id string) (*types.ConnectionRequest, error) {
	return get[types.ConnectionRequest](client, id, crTableName)
}

func (client *client) ListConnectionRequests() ([]types.ConnectionRequest, error) {
	return list[types.ConnectionRequest](client, crTableName)
}

func (client *client) DeleteConnectionRequest(id string) error {
	return client.delete(id, crTableName)
}

func (client *client) UpdateACL(acl *ACL, aclID string) error {
	return update(client, acl, aclID, aclTableName)
}

func (client *client) GetACL(aclID string) (*ACL, error) {
	return get[ACL](client, aclID, aclTableName)
}

func (client *client) ListACLs() ([]ACL, error) {
	return list[ACL](client, aclTableName)
}

func (client *client) DeleteACL(aclID string) error {
	return client.delete(aclID, aclTableName)
}

func (client *client) delete(id, tableName string) error {
	return client.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		return bucket.Delete([]byte(id))
	})
}

func update[T any](client *client, t T, id, tableName string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return client.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))

		if err := bucket.Put([]byte(id), data); err != nil {
			return err
		}
		return nil
	})
}

func get[T any](client *client, id, tableName string) (*T, error) {
	var data []byte
	if err := client.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		data = bucket.Get([]byte(id))
		return nil
	}); err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func list[T any](client *client, tableName string) ([]T, error) {
	var ts []T
	if err := client.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(tableName))
		return bucket.ForEach(func(k, v []byte) error {
			var t T
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			ts = append(ts, t)
			return nil
		})
	}); err != nil {
		return nil, err
	}

	return ts, nil
}
