package player

import (
	"fmt"
	"strings"
)

func (repo *PlayerRepository) buildWhereClause(fields string, s *PlayerSearch) (string, []any) {
	q := `select ` + fields + ` from player where true `
	args := make([]any, 0)
	i := 1

	if s.Name != nil {
		q += fmt.Sprintf(" and name = $%d", i)
		args = append(args, *s.Name)
		i++
	}

	if s.Namelike != nil {
		q += fmt.Sprintf(" and name like $%d", i)
		args = append(args, *s.Namelike)
		i++
	}

	if s.OrderColumn != nil && orderColumns[*s.OrderColumn] {
		order := Ascending
		if s.OrderDirection != nil && orderDirections[*s.OrderDirection] {
			order = *s.OrderDirection
		}

		q += fmt.Sprintf(" order by %s %s", *s.OrderColumn, order)
	}

	// limit result length to 1000 per default
	limit := 1000
	if s.Limit != nil {
		limit = *s.Limit
	}
	q += fmt.Sprintf(" limit %d", limit)

	return q, args
}

func (repo *PlayerRepository) Search(s *PlayerSearch) ([]*Player, error) {
	q, args := repo.buildWhereClause(strings.Join(getColumns(repo.dbtype), ","), s)
	rows, err := repo.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	list := make([]*Player, 0)
	for rows.Next() {
		p, err := scanPlayer(rows.Scan)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (repo *PlayerRepository) Count(s *PlayerSearch) (int, error) {
	q, args := repo.buildWhereClause("count(*)", s)
	row := repo.db.QueryRow(q, args...)
	count := 0
	err := row.Scan(&count)
	return count, err
}
