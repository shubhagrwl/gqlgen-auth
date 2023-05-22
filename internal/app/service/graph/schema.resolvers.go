package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"todo/internal/app/service/graph/generated"
	"todo/internal/app/service/graph/model"
)

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.LoginResponse, error) {
	return r.Services.Auth().SignIn(ctx, input)
}

func (r *mutationResolver) Signup(ctx context.Context, input model.UserInput) (*model.LoginResponse, error) {
	return r.Services.Auth().SignUp(ctx, input)
}

func (r *mutationResolver) SendCode(ctx context.Context, input model.SendCode) (*model.Response, error) {
	return r.Services.Auth().SendCode(ctx, input)
}

func (r *mutationResolver) VerifyCode(ctx context.Context, input model.Code) (*model.Success, error) {
	return r.Services.Auth().VerifyCode(ctx, input)
}

func (r *mutationResolver) ResetPassword(ctx context.Context, input model.ResetPassword) (*model.Success, error) {
	return r.Services.Auth().ResetPassword(ctx, input)
}

func (r *queryResolver) GetProfile(ctx context.Context) (*model.User, error) {
	return r.Services.User().GetProfile(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
