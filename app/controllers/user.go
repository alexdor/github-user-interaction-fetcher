package controllers

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type (
	ResponseError struct {
		Error error  `json:"error"`
		User  string `json:"user"`
	}

	response struct {
		Res    []string        `json:"response"`
		Errors []ResponseError `json:"errors"`
	}

	urlsType map[string]uint8

	/**
	* GraphQL Types
	 */
	Node struct {
		Repository struct {
			Url       string
			IsPrivate bool
		}
	}
	PageInfo struct {
		EndCursor   string
		HasNextPage bool
	}
	Issues struct {
		Nodes    []Node
		PageInfo PageInfo
	}
	RepoUrls struct {
		Url       string
		IsPrivate bool
	}
	Repositories struct {
		Nodes    []RepoUrls
		PageInfo PageInfo
	}
	IssueComments struct {
		Nodes    []Node
		PageInfo PageInfo
	}
	RepositoriesContributedTo struct {
		Nodes    []RepoUrls
		PageInfo PageInfo
	}
	initialUserQuery struct {
		User struct {
			Issues                    Issues                    `graphql:"issues(first: 100)"`
			IssueComments             IssueComments             `graphql:"issueComments(first: 100)"`
			Repositories              Repositories              `graphql:"repositories(first: 100)"`
			RepositoriesContributedTo RepositoriesContributedTo `graphql:"repositoriesContributedTo(first: 100)"`
		} `graphql:"user(login: $login)"`
	}
	userQueryWithAfter struct {
		User struct {
			Issues                    Issues                    `graphql:"issues(first: 100, after: $issuesAfter)"`
			IssueComments             IssueComments             `graphql:"issueComments(first: 100, after: $issueComAfter)"`
			Repositories              Repositories              `graphql:"repositories(first: 100, after: $repoAfter)"`
			RepositoriesContributedTo RepositoriesContributedTo `graphql:"repositoriesContributedTo(first: 100, after: $repoConAfter)"`
		} `graphql:"user(login: $login)"`
	}
)

var (
	client *githubv4.Client
)

func Init() error {
	token := os.Getenv("GITHUB_TOKEN")
	if len(token) < 1 {
		return errors.New("No github token provided")
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client = githubv4.NewClient(httpClient)
	return nil
}

/**
 * Controllers
 */

func GetUserInfo(c *gin.Context) {

	var users struct {
		Users []string `binding:"required" json:"users"`
	}

	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	res := response{}
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	urls := urlsType{}
	usernames := map[string]int8{}
	for _, user := range users.Users {
		if usernames[user] != 1 {
			wg.Add(1)
			go collectUserData(ctx, strings.Trim(strings.TrimSpace(user), "\""), &res, urls, mutex, wg)
			usernames[user] = 1
		}
	}
	wg.Wait()
	c.JSON(http.StatusOK, res)
}

/**
 * Helpers
 */

/**
 * Fetch Initial Data for user
 */
func collectUserData(ctx context.Context, user string, res *response, urls urlsType, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	variables := map[string]interface{}{"login": githubv4.String(user)}
	userData := initialUserQuery{}
	err := client.Query(ctx, &userData, variables)
	initialData := userQueryWithAfter(userData)
	wg.Add(1)
	go writeToResponse(&err, res, &initialData, variables, urls, mutex, wg)

	if userData.User.Issues.PageInfo.HasNextPage || userData.User.IssueComments.PageInfo.HasNextPage || userData.User.Repositories.PageInfo.HasNextPage || userData.User.RepositoriesContributedTo.PageInfo.HasNextPage {
		fetchAdditionalData(ctx, res, urls, &initialData, variables, mutex, wg)
	}

}

/**
 * Fetch additional data for user
 */
func fetchAdditionalData(ctx context.Context, res *response, urls urlsType, userRes *userQueryWithAfter, variables map[string]interface{}, mutex *sync.Mutex, wg *sync.WaitGroup) {
	userData := userQueryWithAfter{}
	userStruct := &userRes.User

	issueCursor := &(*userStruct).Issues.PageInfo.EndCursor
	if len(*issueCursor) > 0 {
		variables["issuesAfter"] = githubv4.String(*issueCursor)

	}
	issueComCursor := &(*userStruct).IssueComments.PageInfo.EndCursor
	if len(*issueComCursor) > 0 {
		variables["issueComAfter"] = githubv4.String(*issueComCursor)
	}
	repoCursor := &(*userStruct).Repositories.PageInfo.EndCursor
	if len(*repoCursor) > 0 {
		variables["repoAfter"] = githubv4.String(*repoCursor)
	}
	repoConCursror := &(*userStruct).RepositoriesContributedTo.PageInfo.EndCursor
	if len(*repoConCursror) > 0 {
		variables["repoConAfter"] = githubv4.String(*repoConCursror)
	}

	err := client.Query(ctx, &userData, variables)
	wg.Add(1)
	go writeToResponse(&err, res, &userData, variables, urls, mutex, wg)

	if userData.User.Issues.PageInfo.HasNextPage || userData.User.IssueComments.PageInfo.HasNextPage || userData.User.Repositories.PageInfo.HasNextPage || userData.User.RepositoriesContributedTo.PageInfo.HasNextPage {
		fetchAdditionalData(ctx, res, urls, &userData, variables, mutex, wg)
	}

}

/**
 * Parse responses and write them to http res
 */
func writeToResponse(err *error, res *response, userRes *userQueryWithAfter, variables map[string]interface{}, urls urlsType, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mutex.Lock()
	defer mutex.Unlock()

	if *err != nil {
		user := string(variables["login"].(githubv4.String))
		res.Errors = append(res.Errors, ResponseError{Error: *err, User: user})
	}
	for _, node := range userRes.User.Issues.Nodes {
		if !node.Repository.IsPrivate && urls[node.Repository.Url] != 1 {
			res.Res = append(res.Res, node.Repository.Url)
			urls[node.Repository.Url] = 1
		}

	}
	for _, node := range userRes.User.IssueComments.Nodes {
		if !node.Repository.IsPrivate && urls[node.Repository.Url] != 1 {
			res.Res = append(res.Res, node.Repository.Url)
			urls[node.Repository.Url] = 1
		}

	}

	for _, node := range userRes.User.Repositories.Nodes {
		if !node.IsPrivate && urls[node.Url] != 1 {
			res.Res = append(res.Res, node.Url)
			urls[node.Url] = 1
		}

	}

	for _, node := range userRes.User.RepositoriesContributedTo.Nodes {
		if !node.IsPrivate && urls[node.Url] != 1 {
			res.Res = append(res.Res, node.Url)
			urls[node.Url] = 1
		}

	}
}
