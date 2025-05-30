package repo

import (
	"context"

	gen "github.com/vadim8q258475/store-user-microservice/internal/repo/ent"
	"github.com/vadim8q258475/store-user-microservice/internal/repo/ent/user"
)

type Repo interface {
	Create(ctx context.Context, email, password string) (*gen.User, error)
	List(ctx context.Context) ([]*gen.User, error)
	GetByEmail(ctx context.Context, email string) (*gen.User, error)
	GetByID(ctx context.Context, id uint32) (*gen.User, error)
	Update(ctx context.Context, id uint32, email, password string) (*gen.User, error)
	Delete(ctx context.Context, id uint32) error
}

type repo struct {
	client *gen.Client
}

func NewRepo(client *gen.Client) Repo {
	return &repo{
		client: client,
	}
}

func (r *repo) Create(ctx context.Context, email, password string) (*gen.User, error) {
	return r.client.User.Create().SetEmail(email).SetPassword(password).Save(context.Background())
}

func (r *repo) List(ctx context.Context) ([]*gen.User, error) {
	return r.client.User.Query().All(ctx)
}

func (r *repo) GetByEmail(ctx context.Context, email string) (*gen.User, error) {
	return r.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
}

func (r *repo) GetByID(ctx context.Context, id uint32) (*gen.User, error) {
	return r.client.User.Get(context.Background(), int(id))
}

func (r *repo) Update(ctx context.Context, id uint32, email, password string) (*gen.User, error) {
	return r.client.User.UpdateOneID(int(id)).SetEmail(email).SetPassword(password).Save(ctx)
}

func (r *repo) Delete(ctx context.Context, id uint32) error {
	return r.client.User.DeleteOneID(int(id)).Exec(ctx)
}
