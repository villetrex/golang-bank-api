package tutorial

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"tutorial.sqlc.dev/app/util"
)

func TestCreateAuthor(t *testing.T) Author {
	arg := CreateAuthorParams{
		Name: util.RandomOwner(),
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	}
	author, err := testQueries.CreateAuthor(context.Background(), arg)
	require.NoError(t, err) // note, t comes from the 'testing' package
	require.NotEmpty(t, author)

	require.Equal(t, arg.Name, author.Name)
	require.Equal(t, arg.Bio, author.Bio)

	require.NotZero(t, author.ID)
	// require.WithinDuration(author.CreatedAt)
	// require.EqualError(t, err, sql.ErrNoRows.Error())
	// require.Len(t, author, 5)
	return author
}

func TestGetAuthor(t *testing.T) {

}
