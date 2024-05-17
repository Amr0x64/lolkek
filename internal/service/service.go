package service

import (
	"context"

	"encoding/json"

	"log"

	"wb-l0/internal/cache"
	"wb-l0/internal/models"
	"wb-l0/internal/repository"

	"github.com/nats-io/stan.go"
)

type Service struct {
	sc    stan.Conn
	repo  *repository.Repo
	Cache cache.Cache
}

func New(sc stan.Conn, repository *repository.Repo) *Service {
	return &Service{
		sc:    sc,
		Cache: cache.New(),
		repo:  repository,
	}
}

func (s *Service) ListenNats(ctx context.Context, channel string) error {
	sub, err := s.sc.Subscribe(channel, func(m *stan.Msg) {
		data := models.Order{}
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			log.Printf("Nats: unmarshal %v", err)
		}
		log.Println(data)
		err = data.Validate()
		if err != nil {
			log.Println(err)
			return
		}

		for i := range data.Items {
			data.Items[i].OrderUid = data.OrderUid
		}

		err = s.repo.Save(ctx, data)
		if err != nil {
			log.Println(err)
		}

		s.Cache.Set(data)
		log.Println(s.Cache.Get(data.OrderUid))
	})

	if err != nil {
		return err
	}

	defer func() {
		err := sub.Unsubscribe()
		if err != nil {
			log.Println(err)
		}
		log.Println("Отписка от канала")
		err = s.sc.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()

	return nil
}
//Иницализируем кеш из бд
func (s *Service) StartCache(ctx context.Context) error {
	cacheDb, err := s.repo.LoadCache(ctx)
	if err != nil {
		return err
	}
	
	for _, o := range cacheDb {
		s.Cache.OrderCache[o.OrderUid] = o
		log.Println(o)
	}

	return nil
}