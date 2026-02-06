package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client,status,address,created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	r, err := res.LastInsertId()
	return int(r), err
}

func (s ParcelStore) Get(number int) (p Parcel, err error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row, err := s.db.Query("SELECT * FROM parcel WHERE number= ?", number)
	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}

	defer row.Close()

	// заполните объект Parcel данными из таблицы
	if row.Next() {
		err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return Parcel{}, err
		}
		return p, nil
	}

	return Parcel{}, errors.New("строка не найдена")

}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Parcel
	for rows.Next() {
		r := Parcel{}
		err := rows.Scan(&r.Number, &r.Client, &r.Status, &r.Address, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel

	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	r, err := s.Get(number)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if r.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	} else {
		return errors.New(fmt.Sprintf("нельзя изменить адрес посылки № %d, т.к. статус отличен от зарегистрирован (%s)", number, r.Status))
	}
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	r, err := s.Get(number)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if r.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("Delete from parcel WHERE number = ?", number)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	} else {
		return errors.New(fmt.Sprintf("нельзя удалить посылку № %d, т.к. статус отличен от зарегистрирован (%s)", number, r.Status))
	}
}
