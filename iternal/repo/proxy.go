package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vadim8q258475/store-user-microservice/iternal/cacher"
	gen "github.com/vadim8q258475/store-user-microservice/iternal/repo/ent"
)

const listKey = "list"
const idKeyPrefix = "id:"
const emailKeyPrefix = "email:"

type proxy struct {
	repo   Repo
	cacher cacher.Cacher
}

func NewProxy(repo Repo, cacher cacher.Cacher) Repo {
	return &proxy{
		repo:   repo,
		cacher: cacher,
	}
}

func (p *proxy) Create(ctx context.Context, email, password string) (*gen.User, error) {
	result, err := p.repo.Create(ctx, email, password)
	if err != nil {
		return nil, err
	}
	err = p.cacher.Delete(ctx, listKey)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (p *proxy) List(ctx context.Context) ([]*gen.User, error) {
	value, err := p.cacher.Get(ctx, listKey)
	if err != nil {
		if err == redis.Nil {
			result, err := p.repo.List(ctx)
			if err != nil {
				return nil, err
			}
			bytes, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			err = p.cacher.Set(ctx, listKey, bytes)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, err
	}
	var result []*gen.User
	err = json.Unmarshal(value, &result)
	return result, err
}
func (p *proxy) GetByEmail(ctx context.Context, email string) (*gen.User, error) {
	key := fmt.Sprintf("%s%s", emailKeyPrefix, email)
	value, err := p.cacher.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			result, err := p.repo.GetByEmail(ctx, email)
			if err != nil {
				return nil, err
			}
			bytes, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			err = p.cacher.Set(ctx, key, bytes)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, err
	}
	var result *gen.User
	err = json.Unmarshal(value, &result)
	return result, err
}
func (p *proxy) GetByID(ctx context.Context, id int) (*gen.User, error) {
	key := fmt.Sprintf("%s%d", idKeyPrefix, id)
	value, err := p.cacher.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			result, err := p.repo.GetByID(ctx, id)
			if err != nil {
				return nil, err
			}
			bytes, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			err = p.cacher.Set(ctx, key, bytes)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, err
	}
	var result *gen.User
	err = json.Unmarshal(value, &result)
	return result, err
}
func (p *proxy) Update(ctx context.Context, id int, email, password string) (*gen.User, error) {
	result, err := p.repo.Update(ctx, id, email, password)
	if err != nil {
		return nil, err
	}
	idKey := fmt.Sprintf("%s%d", idKeyPrefix, result.ID)
	emailKey := fmt.Sprintf("%s%s", emailKeyPrefix, result.Email)
	err = p.cacher.Delete(ctx, listKey, idKey, emailKey)
	if err != nil {
		return nil, err
	}
	return result, err
}
func (p *proxy) Delete(ctx context.Context, id int) error {
	result, err := p.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	idKey := fmt.Sprintf("%s%d", idKeyPrefix, result.ID)
	emailKey := fmt.Sprintf("%s%s", emailKeyPrefix, result.Email)

	err = p.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return p.cacher.Delete(ctx, listKey, idKey, emailKey)
}
