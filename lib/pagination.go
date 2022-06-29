package lib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Pagination struct {
	page    int
	count   int
	orderBy string
	sort    string
}

func NewPagination(page int, count int, orderBy string, sort string) Pagination {
	if len(sort) == 0 {
		sort = "asc"
	}

	return Pagination{
		page:    page,
		count:   count,
		orderBy: strings.ToLower(orderBy),
		sort:    strings.ToLower(sort),
	}
}

func (p *Pagination) Verify(list []string) error {
	countMap := map[int]int{
		10:  10,
		30:  30,
		50:  50,
		100: 100}

	numberRegexp := regexp.MustCompile(`\d`)
	if !numberRegexp.MatchString(strconv.Itoa(p.page)) {
		return fmt.Errorf("page format is ^[0-9]*$")
	}

	if countMap[p.count] == 0 {
		return fmt.Errorf("count not in range 10, 30, 50, 100")
	}

	if p.sort != "asc" && p.sort != "desc" {
		return fmt.Errorf("sort must be asc or desc")
	}

	for _, v := range list {
		if p.orderBy == v {
			return nil
		}
	}
	return fmt.Errorf("orderBy not in range %s", strings.Join(list, " "))
}

func (p *Pagination) GetOffset() int {
	return (p.page - 1) * p.count
}

func (p *Pagination) GetOrderString() string {
	return fmt.Sprintf("%s %s", p.orderBy, p.sort)
}

func (p *Pagination) GetCount() int {
	return p.count
}
