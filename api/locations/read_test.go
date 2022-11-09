package locations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationsQueryNoFilters(t *testing.T) {
	f := LocationFilter{}

	builder, err := ListLocationsQuery(&f)
	assert.NoError(t, err)

	sql, args, err := builder.ToSql()
	assert.NoError(t, err)
	assert.Empty(t, args) // args is empty for empty LocationFilter{}
	assert.NotEmpty(t, sql)
}

func TestLocationsQueryProviderFilter(t *testing.T) {
	p := "mvp"
	f := LocationFilter{Provider: &p}
	builder, err := ListLocationsQuery(&f)
	assert.NoError(t, err)

	sql, args, err := builder.ToSql()
	assert.NoError(t, err)
	assert.NotEmpty(t, args)
	assert.NotEmpty(t, sql)
}
