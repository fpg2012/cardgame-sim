package main

import (
	"math/rand"
	"time"
)

type TokenManager struct {
	tokenToName     map[uint64]string
	nameToToken     map[string]uint64
	tokenExpireTime map[uint64]time.Time
}

func NewTokenManager() *TokenManager {
	token := new(TokenManager)
	token.nameToToken = make(map[string]uint64)
	token.tokenToName = make(map[uint64]string)
	token.tokenExpireTime = make(map[uint64]time.Time)
	return token
}

func (tm *TokenManager) AllocateToken(username string) (token uint64, err error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for {
		token = uint64(r.Int63())
		_, ok := tm.IsTokenUsed(token)
		if !ok {
			break
		}
	}
	tm.tokenToName[token] = username
	tm.nameToToken[username] = token
	// expire after 12 hours
	tm.tokenExpireTime[token] = time.Now().Add(time.Hour * 12)
	err = nil
	return
}

func (tm *TokenManager) IsTokenUsed(token uint64) (string, bool) {
	// check whether token exists
	username, ok := tm.tokenToName[token]
	// check whether token expired
	expiretime, ok2 := tm.tokenExpireTime[token]
	ok_final := ok && ok2
	if ok_final && time.Now().After(expiretime) {
		tm.ReleaseToken(token)
		ok_final = false
	}
	return username, ok
}

// check whether token and username matches
func (tm *TokenManager) ValidateToken(token uint64, username string) bool {
	_username, ok1 := tm.IsTokenUsed(token)
	return ok1 && _username == username
}

func (tm *TokenManager) ReleaseToken(token uint64) {
	username, ok := tm.IsTokenUsed(token)
	if ok {
		delete(tm.tokenToName, token)
		delete(tm.nameToToken, username)
	}
}
