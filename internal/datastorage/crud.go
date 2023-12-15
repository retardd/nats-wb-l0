package datastorage

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"l0/internal/datastorage/structure"
	"log"
)

type Db struct {
	client Client
}

func NewDB(cln Client) *Db {
	dtb := Db{
		client: cln,
	}
	return &dtb
}

func errCheck(fn func() error) error {
	if err := fn(); err != nil {
		var psqlErr *pgconn.PgError
		if errors.As(err, &psqlErr) {
			log.Printf("SQLERROR: %s", psqlErr.Message)
			return psqlErr
		}
		log.Printf("ERROR: %s", err)
		return err
	}
	return nil
}

func errCheckWithRows(fn func() (error, pgx.Rows)) (error, pgx.Rows) {
	err, rows := fn()
	if err != nil {
		var psqlErr *pgconn.PgError
		if errors.As(err, &psqlErr) {
			log.Printf("SQLERROR: %s", psqlErr.Message)
			return psqlErr, rows
		}
		log.Printf("ERROR: %s", err)
		return err, rows
	}
	return nil, rows
}

func (db *Db) InsertAll(ctx context.Context, model *structure.Model) (temp int, err error) {

	query1 := `INSERT INTO model 
					(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service,
					shardkey, sm_id, date_created, oof_shard) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

	err = errCheck(func() error {
		tempErr := db.client.QueryRow(ctx, query1, model.OrderUid,
			model.TrackNumber, model.Entry, model.Locale, model.InternalSignature, model.CustomerId,
			model.DeliveryService, model.ShardKey, model.SmId, model.DateCreated, model.OofShard).Scan(&temp)
		return tempErr
	})

	if err != nil {
		return 0, err
	}

	query2 := `INSERT INTO payment 
					(fk_c_id, transcation, request_id, currency, provider, amount, payment_dt, bank,
					 delivery_cost, goods_total, custom_fee) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING fk_c_id`

	err = errCheck(func() error {
		tempErr := db.client.QueryRow(ctx, query2, temp, model.Payment.Transcation, model.Payment.RequestId,
			model.Payment.Currency, model.Payment.Provider, model.Payment.Amount, model.Payment.PaymentDt,
			model.Payment.Bank, model.Payment.DeliveryCost, model.Payment.GoodsTotal, model.Payment.CustomFee).Scan(&temp)
		return tempErr
	})

	if err != nil {
		return 0, err
	}

	query3 := `INSERT INTO delivery 
					(fk_c_id, name, phone, zip, city, address, region, email) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING fk_c_id`

	err = errCheck(func() error {
		tempErr := db.client.QueryRow(ctx, query3, temp, model.Delivery.Name, model.Delivery.Phone, model.Delivery.Zip,
			model.Delivery.City, model.Delivery.Address, model.Delivery.Region, model.Delivery.Email).Scan(&temp)
		return tempErr
	})

	if err != nil {
		return 0, err
	}

	query4 := `INSERT INTO items
					(fk_c_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING fk_c_id`

	for _, item := range model.Items {
		err = errCheck(func() error {
			tempErr := db.client.QueryRow(ctx, query4, temp, item.ChrtId, item.TrackNumber, item.Price, item.Rid,
				item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status).Scan(&temp)
			return tempErr
		})
		if err != nil {
			return 0, err
		}
	}

	return
}

func (db *Db) SelectOne(ctx context.Context, mdId int) (err error, model structure.Model) {
	query1 := `SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service,
					shardkey, sm_id, date_created, oof_shard FROM model WHERE id = $1`

	err = errCheck(func() error {
		tempErr := db.client.QueryRow(ctx, query1, mdId).Scan(&model.OrderUid,
			&model.TrackNumber, &model.Entry, &model.Locale, &model.InternalSignature, &model.CustomerId,
			&model.DeliveryService, &model.ShardKey, &model.SmId, &model.DateCreated, &model.OofShard)
		return tempErr
	})

	if err != nil {
		return err, model
	}

	query2 := `SELECT transcation, request_id, currency, provider, amount, payment_dt, bank,
					 delivery_cost, goods_total, custom_fee FROM payment WHERE fk_c_id = $1`

	err = errCheck(func() error {
		tempErr := db.client.QueryRow(ctx, query2, mdId).Scan(&model.Payment.Transcation, &model.Payment.RequestId,
			&model.Payment.Currency, &model.Payment.Provider, &model.Payment.Amount, &model.Payment.PaymentDt,
			&model.Payment.Bank, &model.Payment.DeliveryCost, &model.Payment.GoodsTotal, &model.Payment.CustomFee)
		return tempErr
	})

	if err != nil {
		return err, model
	}

	query3 := `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE fk_c_id = $1`

	err = errCheck(func() error {
		tempErr := db.client.QueryRow(ctx, query3, mdId).Scan(&model.Delivery.Name, &model.Delivery.Phone, &model.Delivery.Zip,
			&model.Delivery.City, &model.Delivery.Address, &model.Delivery.Region, &model.Delivery.Email)
		return tempErr
	})

	if err != nil {
		return err, model
	}

	query4 := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, 
       nm_id, brand, status FROM items WHERE fk_c_id = $1`

	var items pgx.Rows

	err, items = errCheckWithRows(func() (error, pgx.Rows) {
		rows, tempErr := db.client.Query(ctx, query4, mdId)
		return tempErr, rows
	})

	if err != nil {
		return err, model
	}

	for items.Next() {
		var item structure.Item
		err = errCheck(func() error {
			tempErr := items.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid,
				&item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
			return tempErr
		})

		if err != nil {
			return err, model
		}

		model.Items = append(model.Items, item)
	}

	return nil, model
}

func (db *Db) FillCache(ctx context.Context, size int) (error, map[int]structure.Model, int) {
	temp := make(map[int]structure.Model, size)
	query := `SELECT id FROM model ORDER BY id DESC LIMIT $1`
	err, rows := errCheckWithRows(func() (error, pgx.Rows) {
		rows, tempErr := db.client.Query(ctx, query, size)
		return tempErr, rows
	})

	if err != nil {
		return err, temp, 0
	}
	var modelId int
	for rows.Next() {
		err = rows.Scan(&modelId)
		if err != nil {
			return err, temp, 0
		}
		err, model := db.SelectOne(ctx, modelId)

		if err != nil {
			return err, temp, 0
		}

		temp[modelId] = model
	}

	return nil, temp, modelId
}
