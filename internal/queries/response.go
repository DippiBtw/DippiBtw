package queries

// ---------- User Query ----------
type UserQueryResponse struct {
	User struct {
		Email     string `json:"email"`
		CreatedAt string `json:"createdAt"`
		Followers struct {
			TotalCount int `json:"totalCount"`
		} `json:"followers"`
		Location string `json:"location"`
		Name     string `json:"name"`
		Id       string `json:"id"`
	} `json:"user"`
}

// ---------- Commits Query ----------
type CommitsQueryResponse struct {
	User struct {
		ContributionsCollection struct {
			ContributionCalendar struct {
				TotalContributions int `json:"totalContributions"`
			} `json:"contributionCalendar"`
		} `json:"contributionsCollection"`
	} `json:"user"`
}

// ---------- Contributions Query ----------
type ContributionsQueryResponse struct {
	User struct {
		ContributionsCollection struct {
			CommitContributionsByRepository []struct {
				Repository struct {
					NameWithOwner string `json:"nameWithOwner"`
				} `json:"repository"`
			} `json:"commitContributionsByRepository"`
		} `json:"contributionsCollection"`
	} `json:"user"`
}

// ---------- Repositories Query ----------
type RepositoriesQueryResponse struct {
	User struct {
		Repositories struct {
			TotalCount int `json:"totalCount"`
			Edges      []struct {
				Node struct {
					NameWithOwner  string `json:"nameWithOwner"`
					StargazerCount int    `json:"stargazerCount"`
				} `json:"node"`
			} `json:"edges"`
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
		} `json:"repositories"`
	} `json:"user"`
}

// ---------- Lines of Code Query ----------
type LinesOfCodeQueryResponse struct {
	Repository struct {
		DefaultBranchRef struct {
			Target struct {
				History struct {
					TotalCount int `json:"totalCount"`
					Edges      []struct {
						Node struct {
							CommittedDate string `json:"committedDate"`
							Author        struct {
								User *struct {
									ID string `json:"id"`
								} `json:"user"`
							} `json:"author"`
							Deletions int `json:"deletions"`
							Additions int `json:"additions"`
						} `json:"node"`
					} `json:"edges"`
					PageInfo struct {
						EndCursor   string `json:"endCursor"`
						HasNextPage bool   `json:"hasNextPage"`
					} `json:"pageInfo"`
				} `json:"history"`
			} `json:"target"`
		} `json:"defaultBranchRef"`
	} `json:"repository"`
}

type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type PageInfoExtractor[T any] func(*T) PageInfo

func extractRepoPageInfo(resp *RepositoriesQueryResponse) PageInfo {
	return resp.User.Repositories.PageInfo
}

func extractLocPageInfo(resp *LinesOfCodeQueryResponse) PageInfo {
	return resp.Repository.DefaultBranchRef.Target.History.PageInfo
}
