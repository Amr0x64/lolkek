package repository

import (
	"context"

	"wb-l0/internal/models"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{
		db: conn,
	}
}

func (r *Repo) Save(ctx context.Context, o models.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	var payment_id, delivery_id int
	err = tx.QueryRow(ctx, "insert into delivery(name,phone,zip,city,address,region,email) values($1,$2,$3,$4,$5,$6,$7) returning id",
		o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip, o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email).Scan(&delivery_id)
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, `insert into payment(transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) 
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id`,
		o.Payment.Transaction, o.Payment.RequestId, o.Payment.Currency, o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank, o.Payment.DeliveryCost,
		o.Payment.GoodsTotal, o.Payment.CustomFee).Scan(&payment_id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `insert into orders(order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard,
		delivery_id,payment_id) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		o.OrderUid, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature, o.CustomerId, o.DeliveryService, o.ShardKey, o.SmId, o.DateCreated, o.OofShard, delivery_id, payment_id)
	if err != nil {
		return err
	}

	for i := 0; i < len(o.Items); i++ {
		_, err = tx.Exec(ctx, "insert into item(name,sale,size,total_price,nm_id,brand,status,chrt_id,track_number,price,rid,order_uid) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
			o.Items[i].Name, o.Items[i].Sale, o.Items[i].Size, o.Items[i].TotalPrice, o.Items[i].NmId, o.Items[i].Brand, o.Items[i].Status, o.Items[i].ChrtId, o.Items[i].TrackNumber,
			o.Items[i].Price, o.Items[i].Rid, o.OrderUid)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) LoadCache(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	rows, err := r.db.Query(ctx, `select orders.order_uid,orders.track_number,entry,locale,internal_signature,customer_id,
		delivery_service,shardkey,sm_id,date_created,oof_shard,delivery.name,phone,zip,city,address,region,email,transaction,
		request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee from orders 
		join delivery ON delivery.id = orders.delivery_id join payment ON payment.id = orders.payment_id`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o models.Order
		err = rows.Scan(&o.OrderUid, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerId, &o.DeliveryService,
			&o.ShardKey, &o.SmId, &o.DateCreated, &o.OofShard, &o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City,
			&o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email, &o.Payment.Transaction, &o.Payment.RequestId,
			&o.Payment.Currency, &o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank, &o.Payment.DeliveryCost,
			&o.Payment.GoodsTotal, &o.Payment.CustomFee)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	rows.Close()
	if rows.Err() != nil {
		return nil, err
	}
	for i := 0; i < len(orders); i++ {
		rows, err := r.db.Query(ctx, `select name,sale,size,total_price,nm_id,brand,status,chrt_id,item.track_number,price,
		rid from item where order_uid=$1`, orders[i].OrderUid)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var item models.Item
			err = rows.Scan(&item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status, &item.ChrtId,
				&item.TrackNumber, &item.Price, &item.Rid)
			if err != nil {
				return nil, err
			}
			orders[i].Items = append(orders[i].Items, item)
		}
		rows.Close()
		if rows.Err() != nil {
			return nil, err
		}

	}
	return orders, nil
}
