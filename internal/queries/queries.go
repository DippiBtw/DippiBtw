package queries

const id = 115025256
const name = "DippiBtw"

// Get basic user information
const UserQuery = `query($login: String!) {
  user(login: $login) {
		email
		createdAt
		followers {
			totalCount
		}
		location
		name
		id
  }
}`

// Get all commits made by the user since the given date
const CommitsQuery = `query($login: String!, $from: DateTime, $to: DateTime) {
  user(login: $login) {
    contributionsCollection(from: $from, to: $to) {
      contributionCalendar {
        totalContributions
      }
    }
  }
}`

const ContributionQuery = `query($login: String!, $from: DateTime, $to: DateTime) {
  user(login: $login) {
    contributionsCollection(from: $from, to: $to) {
      commitContributionsByRepository {
        repository {
          nameWithOwner
        }
      }
    }
  }
}`

// Get all repositories of a user
// Able to specify affiliation of the user to the repository
const ReposQuery = `query($login: String!, $affiliations: [RepositoryAffiliation], $cursor: String) {
  user(login: $login) {
    repositories(first: 100, after: $cursor, ownerAffiliations: $affiliations) {
			totalCount
			edges {
				node {
					nameWithOwner
					stargazerCount
				}
			}
			pageInfo {
				endCursor
				hasNextPage
			}
		}
  }
}`

const LinesOfCodeQuery = `query ($repo_name: String!, $owner: String!, $cursor: String, $author: CommitAuthor) {
		repository(name: $repo_name, owner: $owner) {
				defaultBranchRef {
						target {
								... on Commit {
										history(first: 100, after: $cursor, author: $author) {
												totalCount
												edges {
														node {
																... on Commit {
																		committedDate
																}
																author {
																		user {
																				id
																		}
																}
																deletions
																additions
														}
												}
												pageInfo {
														endCursor
														hasNextPage
												}
										}
								}
						}
				}
		}
}
`
