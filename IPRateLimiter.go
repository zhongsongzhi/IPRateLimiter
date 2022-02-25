package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu *sync.RWMutex
	r rate.Limit
	b int
}

func NewIPRateLimiter(r rate.Limit,b int) *IPRateLimiter  {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu: &sync.RWMutex{},
		r: r,
		b: b,
	}
	return i
}

//添加ip地址
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()
	fmt.Println("添加ip地址",ip) //注释
	limiter := rate.NewLimiter(i.r,i.b)//r时间b个令牌

	i.ips[ip] = limiter
	return limiter
}

//如果ip地址存在则返回其limiter
func (i IPRateLimiter) GetLimiter(ip string) *rate.Limiter  {
	i.mu.Lock()
	limiter, exists := i.ips[ip]
	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}
	i.mu.Unlock()
	return limiter
}

func GetUserIP(r *http.Request)	string  {
	forwarded := r.Header.Get("X-FORWARDED-FOR");
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}



