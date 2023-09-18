package hw10programoptimization

import (
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
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
	decoder := jsoniter.ConfigFastest.NewDecoder(r)
	i := 0
	for decoder.More() {
		var user User
		if err := decoder.Decode(&user); err != nil {
			return result, err
		}
		result[i] = user
		i++
	}
	return result, err
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if user.Email == "" {
			continue
		}
		userHost := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		if userHost == "" {
			continue
		}
		userDomain := strings.ToLower(strings.SplitN(userHost, ".", 2)[1])

		if userDomain == strings.ToLower(domain) {
			result[userHost]++
		}
	}
	return result, nil
}
