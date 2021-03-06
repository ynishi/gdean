package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ynishi/gdean/gql/graph/generated"
	"github.com/ynishi/gdean/gql/graph/model"
	pb "github.com/ynishi/gdean/pb/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (r *mutationResolver) CreateIssue(ctx context.Context, userID string, input model.IssueInput) (*model.Issue, error) {
	var branches = make([]*pb.Branch, len(input.Branches))
	for i := range branches {
		istr := strconv.Itoa(i)
		branches[i] = &pb.Branch{
			Id:      &istr,
			Content: input.Branches[i],
		}
	}
	var desc string = ""
	if input.Desc != nil {
		desc = *input.Desc
	}
	req := &pb.CreateIssueRequest{
		UserId: userID,
		Issue: &pb.Issue{
			Title:    input.Title,
			Desc:     desc,
			Branches: branches,
		},
	}
	resp, err := r.IssueClient.CreateIssue(ctx, req)
	if err != nil {
		return nil, err
	}
	return IssueFromPB(resp.Issue), nil
}

func (r *mutationResolver) UpdateIssue(ctx context.Context, issueID string, input model.IssueInput) (*model.Issue, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateIssueComment(ctx context.Context, input *model.NewIssueComment) (*model.Comment, error) {
	resp, err := r.IssueClient.GetIssue(ctx, &pb.GetIssueRequest{IssueId: input.IssueID})
	if err != nil {
		return nil, err
	}
	issue := resp.Issue
	issue.Comments = append(
		issue.Comments,
		&pb.Comment{
			Author:  input.UserID,
			Content: input.Content,
		},
	)
	updatedResp, err := r.IssueClient.UpdateIssue(ctx, &pb.UpdateIssueRequest{
		Issue:     issue,
		FieldMask: &fieldmaskpb.FieldMask{Paths: []string{"Comments"}},
	},
	)
	if err != nil {
		return nil, err
	}

	return &model.Comment{
		ID:        *updatedResp.Issue.Comments[len(issue.Comments)].Id,
		Content:   input.Content,
		CreatedAt: updatedResp.Issue.Comments[len(issue.Comments)].CreateTime.AsTime(),
	}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	resp, err := r.UserClient.CreateUser(ctx, &pb.CreateUserRequest{User: &pb.User{Name: input.Name}})
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:   *resp.User.Id,
		Name: resp.User.Name,
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, userID string, input model.UserInput) (*model.User, error) {
	req := pb.UpdateUserRequest{
		User: &pb.User{
			Id:   &userID,
			Name: input.Name,
		},
		FieldMask: &fieldmaskpb.FieldMask{
			Paths: []string{"Name"},
		},
	}
	resp, err := r.UserClient.UpdateUser(ctx, &req)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:   *resp.User.Id,
		Name: resp.User.Name,
	}, nil
}

func (r *queryResolver) Issues(ctx context.Context, userID string) ([]*model.IssueSummary, error) {
	resp, err := r.IssueClient.ListIssues(ctx, &pb.ListIssuesRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	var issueSummaryList = make([]*model.IssueSummary, len(resp.Issues))

	for i := range issueSummaryList {
		issue := resp.Issues[i]
		var decidedTitle string
		for i := range issue.Branches {
			if *issue.Branches[i].Id == issue.Decided {
				decidedTitle = BranchToTitle(issue.Branches[i])
				break
			}
		}
		issueSummaryList[i] = &model.IssueSummary{
			ID:                 *issue.Id,
			Title:              issue.Title,
			Desc:               issue.Desc,
			AuthorName:         issue.Author.Name,
			ContributerCount:   len(issue.Contributers),
			DecidedBranchTitle: decidedTitle,
			AnalysisCount:      len(issue.Results),
			CreatedAt:          issue.CreateTime.AsTime(),
		}
	}

	return issueSummaryList, nil
}

func (r *queryResolver) Issue(ctx context.Context, id string) (*model.Issue, error) {
	resp, err := r.IssueClient.GetIssue(ctx, &pb.GetIssueRequest{IssueId: id})
	if err != nil {
		return nil, err
	}
	return IssueFromPB(resp.Issue), nil
}

func (r *queryResolver) User(ctx context.Context, userID string) (*model.User, error) {
	resp, err := r.UserClient.GetUser(ctx, &pb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:   *resp.User.Id,
		Name: resp.User.Name,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func BranchToTitle(branch *pb.Branch) string {
	return string([]rune(branch.Content)[:15])
}
func IssueFromPB(issue *pb.Issue) *model.Issue {
	var branchesRes = make([]*model.Branch, len(issue.Branches))
	for i := range branchesRes {
		branchesRes[i] = &model.Branch{
			IssueID:   *issue.Id,
			Title:     BranchToTitle(issue.Branches[i]),
			CreatedAt: issue.CreateTime.AsTime(),
		}
	}

	return &model.Issue{
		ID:       *issue.Id,
		Title:    issue.Title,
		Desc:     issue.Desc,
		Branches: branchesRes,
		Author: &model.User{
			ID:   *issue.Author.Id,
			Name: issue.Author.Name,
		},
		CreatedAt: issue.CreateTime.AsTime(),
	}
}
func (r *queryResolver) Comments(ctx context.Context, issueID string) ([]*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *queryResolver) Branches(ctx context.Context, issueID string) ([]*model.Branch, error) {
	panic(fmt.Errorf("not implemented"))
}
