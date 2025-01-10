package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "./data/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	testParcel, err := store.Get(id)
	require.NoError(t, err)

	// Исключаем поле Number из проверки
	parcel.Number = testParcel.Number
	assert.Equal(t, parcel, testParcel)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "./data/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	testParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, testParcel.Address)

	err = store.Delete(id)
	require.NoError(t, err)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "./data/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	testParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, testParcel.Status)

	err = store.Delete(id)
	require.NoError(t, err)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "./data/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.Greater(t, id, 0)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	for _, parcel := range storedParcels {
		mappedParcel, ok := parcelMap[parcel.Number]
		assert.True(t, ok)
		assert.Equal(t, mappedParcel, parcel)
	}
}
