package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)

	i := 0
	for scanner.Scan() {
		var user User
		if err = easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return
		}
		result[i] = user
		i++
	}

	err = scanner.Err()
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	regDomain := regexp.MustCompile("\\." + domain)

	for _, user := range u {
		if regDomain.MatchString(user.Email) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
