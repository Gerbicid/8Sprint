package main

import (
	"database/sql"
	"fmt"
)

const (
	addQuery = `
INSERT INTO parcel (client, status, address, created_at)
VALUES (:client, :status, :address, :created_at)
`

	getByNumberQuery = `
SELECT
	number,
	client,
	status,
	address,
	created_at
FROM parcel p
WHERE p.number = :number
`

	getByClientQuery = `
SELECT
	number,
	client,
	status,
	address,
	created_at
FROM parcel p
WHERE p.client = :client
`

	setStatusQuery = `
UPDATE parcel
SET status = :status
WHERE number = :number
`

	setAddressQuery = `
UPDATE parcel
SET address = :address
WHERE number = :number AND status = :status
`

	deleteQuery = `
DELETE FROM parcel
WHERE number = :number AND status = :status
`
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(addQuery,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, fmt.Errorf("add query execution error: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("can't get last insertion id: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow(getByNumberQuery, sql.Named("number", number))
	p := Parcel{}
	err := row.Scan(
		&p.Number,
		&p.Client,
		&p.Status,
		&p.Address,
		&p.CreatedAt,
	)
	if err != nil {
		return Parcel{}, fmt.Errorf("get parcel by id error: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(getByClientQuery, sql.Named("client", client))
	if err != nil {
		return nil, fmt.Errorf("get parcels by client id error: %w", err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}()

	res := make([]Parcel, 0)
	for rows.Next() {
		p := Parcel{}
		err = rows.Scan(
			&p.Number,
			&p.Client,
			&p.Status,
			&p.Address,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("get parcels by client id error: %w", err)
		}
		res = append(res, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec(setStatusQuery,
		sql.Named("number", number),
		sql.Named("status", status),
	)
	if err != nil {
		return fmt.Errorf("set parcel status error: %w", err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec(setAddressQuery,
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		return fmt.Errorf("set parcel address error: %w", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec(deleteQuery,
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		return fmt.Errorf("delete parcel error: %w", err)
	}

	return nil
}
