package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"github.com/ynishi/gdean/gql/client"
)

type Resolver struct {
	IssueClient *client.IssueClient
	UserClient  *client.UserClient
}
