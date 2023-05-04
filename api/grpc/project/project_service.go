package project

import (
	"context"
	connect_go "github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/project/projectpb"
	"github.com/loomi-labs/star-scope/grpc/project/projectpb/projectpbconnect"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type ProjectService struct {
	projectpbconnect.UnimplementedProjectServiceHandler
}

func NewProjectServiceHandler() projectpbconnect.ProjectServiceHandler {
	return &ProjectService{}
}

func (e ProjectService) ListProjects(ctx context.Context, _ *connect_go.Request[emptypb.Empty]) (*connect_go.Response[projectpb.ListProjectsResponse], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	projects, err := user.
		QueryProjects().
		WithChannels().
		All(ctx)
	if err != nil {
		log.Sugar.Errorf("failed to query projects: %v", err)
		return nil, err
	}

	var projectPbs []*projectpb.Project
	for _, project := range projects {
		var channelPbs []*projectpb.Channel
		for _, channel := range project.Edges.Channels {
			channelPbs = append(channelPbs, &projectpb.Channel{
				Name: channel.Name,
			})
		}
		projectPbs = append(projectPbs, &projectpb.Project{
			Name:     project.Name,
			Channels: channelPbs,
		})
	}

	return connect_go.NewResponse(&projectpb.ListProjectsResponse{
		Projects: projectPbs,
	}), nil
}
