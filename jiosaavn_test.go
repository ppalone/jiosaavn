package jiosaavn_test

import (
	"context"
	"testing"

	"github.com/ppalone/jiosaavn"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := jiosaavn.NewClient(nil)
	assert.NotNil(t, c)
}

func TestSearchSongs(t *testing.T) {
	t.Run("with empty search query", func(t *testing.T) {
		c := jiosaavn.NewClient(nil)
		_, err := c.SearchSongs(context.Background(), "")
		assert.ErrorContains(t, err, "search query cannot be empty")
	})

	t.Run("with no search options", func(t *testing.T) {
		c := jiosaavn.NewClient(nil)
		res, err := c.SearchSongs(context.Background(), "Animals")
		assert.NoError(t, err)
		assert.Equal(t, 1, res.Page)
		assert.Equal(t, 10, res.Size)
		assert.NotEmpty(t, res.Songs)
		assert.True(t, res.HasNext)
	})

	t.Run("with valid limit search option", func(t *testing.T) {
		c := jiosaavn.NewClient(nil)
		limit := 40
		res, err := c.SearchSongs(context.Background(), "Animals", jiosaavn.WithLimit(limit))
		assert.NoError(t, err)
		assert.Equal(t, 1, res.Page)
		assert.Equal(t, limit, res.Size)
		assert.NotEmpty(t, res.Songs)
		assert.True(t, res.HasNext)
	})

	t.Run("with invalid limit search option", func(t *testing.T) {
		c := jiosaavn.NewClient(nil)
		limit := 50
		_, err := c.SearchSongs(context.Background(), "Animals", jiosaavn.WithLimit(limit))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "limit must be between 10 and 40")
	})

	t.Run("with page search option", func(t *testing.T) {
		c := jiosaavn.NewClient(nil)
		res1, err := c.SearchSongs(context.Background(), "Animals")
		assert.NoError(t, err)

		res2, err := c.SearchSongs(context.Background(), "Animals", jiosaavn.WithPage(2))
		assert.NoError(t, err)
		assert.NotEmpty(t, res2.Songs)
		assert.Equal(t, 2, res2.Page)
		assert.Equal(t, 10, res2.Size)

		assert.NotEqual(t, res1.Songs, res2.Songs)
	})

	t.Run("with page and limit options", func(t *testing.T) {
		c := jiosaavn.NewClient(nil)
		limit, page := 30, 2
		opts := []jiosaavn.SearchOption{
			jiosaavn.WithLimit(limit),
			jiosaavn.WithPage(page),
		}
		res, err := c.SearchSongs(context.Background(), "Animals", opts...)
		assert.NoError(t, err)
		assert.Equal(t, page, res.Page)
		assert.Equal(t, limit, res.Size)
		assert.NotEmpty(t, res.Songs)
	})

	t.Run("with next results", func(t *testing.T) {
		c := jiosaavn.NewClient(nil)
		res, err := c.SearchSongs(context.Background(), "Animals")
		assert.NoError(t, err)
		assert.Equal(t, 1, res.Page)
		assert.NotEmpty(t, res.Songs)
		assert.True(t, res.HasNext)

		resNext, err := res.Next(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 2, resNext.Page)
		assert.NotEmpty(t, resNext.Songs)
		assert.True(t, resNext.HasNext)
	})

	t.Run("with next and page results", func(t *testing.T) {
		t.Skip("tests pass on local and fail on workflow for some reason")
		c := jiosaavn.NewClient(nil)
		res, err := c.SearchSongs(context.Background(), "Animals")
		assert.NoError(t, err)
		assert.Equal(t, 1, res.Page)
		assert.True(t, res.HasNext)

		resNext, err := res.Next(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 2, resNext.Page)

		resPage, err := c.SearchSongs(context.Background(), "Animals", jiosaavn.WithPage(2))
		assert.NoError(t, err)
		assert.Equal(t, 2, resPage.Page)

		assert.ElementsMatch(t, resNext.Songs, resPage.Songs)
	})

	t.Run("with no search results", func(t *testing.T) {
		c := jiosaavn.NewClient(nil)
		res, err := c.SearchSongs(context.Background(), "qazwsxecrfvtgbyhnujmik")
		assert.NoError(t, err)
		assert.Equal(t, 1, res.Page)
		assert.Equal(t, 0, res.Size)
		assert.False(t, res.HasNext)
	})
}
