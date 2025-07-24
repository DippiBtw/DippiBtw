package queries

import (
	"context"
	"strings"
	"time"
)

type Info struct {
	Commits       int
	Followers     int
	Additions     int
	Deletions     int
	Repos         int
	ContributedTo int
	Stars         int
}

func GetInfo(ctx context.Context) (*Info, error) {
	info := &Info{}
	client := New("https://api.github.com/graphql")
	variables := make(map[string]interface{})
	variables["login"] = "DippiBtw"
	variables["affiliations"] = []string{"OWNER"}
	variables["cursor"] = nil

	user, err := Execute[UserQueryResponse](ctx, client, UserQuery, variables)
	if err != nil {
		return nil, err
	}

	variables["author"] = map[string]interface{}{
		"id": user.User.Id,
	}

	from, _ := time.Parse(time.RFC3339, user.User.CreatedAt)
	commits, err := QueryYearly[CommitsQueryResponse](ctx, client, CommitsQuery, variables, from, time.Now())
	if err != nil {
		return nil, err
	}

	contributions, err := QueryYearly[ContributionsQueryResponse](ctx, client, ContributionQuery, variables, from, time.Now())
	if err != nil {
		return nil, err
	}

	myRepos, err := QueryRecursively[RepositoriesQueryResponse](ctx, client, ReposQuery, variables, extractRepoPageInfo)
	if err != nil {
		return nil, err
	}

	for _, r := range myRepos {
		info.Repos += r.User.Repositories.TotalCount
		for _, edge := range r.User.Repositories.Edges {
			variables["repo_name"] = strings.Split(edge.Node.NameWithOwner, "/")[1]
			variables["owner"] = strings.Split(edge.Node.NameWithOwner, "/")[0]
			loc, err := QueryRecursively[LinesOfCodeQueryResponse](ctx, client, LinesOfCodeQuery, variables, extractLocPageInfo)
			if err != nil {
				return nil, err
			}

			for _, l := range loc {
				for _, h := range l.Repository.DefaultBranchRef.Target.History.Edges {
					info.Additions += h.Node.Additions
					info.Deletions += h.Node.Deletions
				}
			}

			info.Stars += edge.Node.StargazerCount
		}
	}

	for _, c := range commits {
		info.Commits += c.User.ContributionsCollection.ContributionCalendar.TotalContributions
	}
	for _, c := range contributions {
		for _, r := range c.User.ContributionsCollection.CommitContributionsByRepository {
			if !strings.Contains(r.Repository.NameWithOwner, "DippiBtw") {
				info.ContributedTo++
			}
		}
	}
	info.Followers = user.User.Followers.TotalCount
	return info, nil
}
