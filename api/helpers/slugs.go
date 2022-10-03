package helpers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Slugify(inString string) string {

	// Slugify string
	return slug.Make(inString)

}

func GetUniqueSlug(name string, slugMap map[string]bool) (string, error) {

	slugIsTaken := func(str string) bool {
		if _, ok := slugMap[str]; ok {
			return true
		}
		return false
	}

	// slugify the input name
	slug := Slugify(name)

	// if slug is unique without appending an integer, return it
	if !(slugIsTaken(slug)) {
		return slug, nil
	}

	// max 1000 iterations trying to get unique slug
	// if we reach the end of 1000 iterations, it means there are more than 1000 things with the same
	// name in the database table.
	i := 1
	for i < 1000 {
		slug := fmt.Sprintf("%s-%d", slug, i)
		if !(slugIsTaken(slug)) {
			return slug, nil
		}
		i++
	}
	return "", errors.New("reached max iteration %i without finding a unique slug")
}

// NextUniqueSlug returns the next available slug given a table
// contextField is a column name in table; contextString is a value in contextField;
// If contextField="" or contextString="", returned string will be table unique
// If contextField is provided, the returned string will be unique among rows having contextField = contextValue
func NextUniqueSlug(db *pgxpool.Pool, table, field, inString, contextField, contextValue string) (string, error) {

	// SQL Query Builder logic
	sql := func() string {
		// Find inString or any variants suffixed by xxx-1, xxx-2, etc.
		q := fmt.Sprintf(
			`SELECT %s FROM %s WHERE %s ~ ('^'||$1||'(?:-[0-9]+)?$')`,
			field, table, field,
		)
		if contextField != "" && contextValue != "" {
			q += fmt.Sprintf(" AND %s = %s", contextField, contextValue)
		}
		q += fmt.Sprintf(" ORDER BY %s DESC", field)
		return q
	}

	// Slugify string; this is the first choice for a slug if it's not already taken
	slugTry := Slugify(inString)

	ss := make([]string, 0)
	if err := pgxscan.Select(context.Background(), db, &ss, sql(), slugTry); err != nil {
		return "", err
	}

	// If the last entry based on sort DESC is not slugTry
	// we know the slug is free. Regex gets dicey when a slug like
	// "test-watershed-1" exists, "test-watershed" does not, and you
	// run this method on inString for "test-watershed".
	if len(ss) == 0 || ss[len(ss)-1] != slugTry {
		return slugTry, nil
	}
	// If slugTry already taken; add "-1" to the end
	if len(ss) == 1 {
		return fmt.Sprintf("%s-1", slugTry), nil
	}
	// If there are many of these slugs already (i.e. myslug, myslug-1, myslug-2, etc.)
	// iterate to find the next available integer to tag on the end
	largest := 1
	pLargest := &largest
	for idx := range ss[:len(ss)-1] {
		parts := strings.Split(ss[idx], "-")
		i, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			return "", err
		}
		if i > *pLargest {
			*pLargest = i
		}
	}
	return fmt.Sprintf("%s-%d", slugTry, *pLargest+1), nil
}
