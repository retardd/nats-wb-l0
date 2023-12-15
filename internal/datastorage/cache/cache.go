package cache

import (
	"context"
	ds "l0/internal/datastorage"
	"l0/internal/datastorage/structure"
	"l0/utils"
	"sync"
)

type Cache struct {
	Cch    map[int]structure.Model
	Mtx    *sync.RWMutex
	Size   int
	Psql   *ds.Db
	LastId int
}

func InitCache(ctx context.Context) (error, *Cache) {
	cch := Cache{}
	cch.Size = utils.GetIntEnv("SIZE_CCH", 10)
	cch.Mtx = &sync.RWMutex{}
	cc := ds.ConnectionConfig{}
	cc.GettingEnv()
	err, client := ds.CreateClient(context.TODO(), cc)

	if err != nil {
		return err, nil
	}

	cch.Psql = ds.NewDB(client)
	cch.Mtx.Lock()
	err, cch.Cch, cch.LastId = cch.Psql.FillCache(ctx, cch.Size)

	cch.Mtx.Unlock()

	if err != nil {
		return err, nil
	}

	return err, &cch
}

func (cch *Cache) FindData(ctx context.Context, mdId int) (error, *structure.Model) {

	cch.Mtx.RLock()
	model, ok := cch.Cch[mdId]
	cch.Mtx.RUnlock()
	if ok {
		return nil, &model
	}

	err, model := cch.Psql.SelectOne(ctx, mdId)

	if err != nil {
		return err, &model
	}

	cch.Mtx.Lock()
	if len(cch.Cch) > cch.Size-1 {
		for k, _ := range cch.Cch {
			delete(cch.Cch, k)
			break
		}
	}

	cch.Cch[mdId] = model
	cch.Mtx.Unlock()

	return nil, &model
}

func (cch *Cache) AddData(ctx context.Context, model *structure.Model) (error, int) {
	cch.Mtx.Lock()
	if len(cch.Cch) > cch.Size {
		for k, _ := range cch.Cch {
			delete(cch.Cch, k)
			break
		}
	}

	mdId, err := cch.Psql.InsertAll(ctx, model)

	cch.LastId = mdId

	if err != nil {
		return err, mdId
	}

	cch.Cch[mdId] = *model
	cch.Mtx.Unlock()

	return nil, mdId
}
