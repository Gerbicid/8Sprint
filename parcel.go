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
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec(addQuery,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, fmt.Errorf("add query execution error: %w", err)
	}
	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("can't get last insertion id: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow(getByNumberQuery, sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
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
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
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

	// заполните срез Parcel данными из таблицы
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

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
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
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
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
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec(deleteQuery,
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		return fmt.Errorf("delete parcel error: %w", err)
	}

	return nil
}
