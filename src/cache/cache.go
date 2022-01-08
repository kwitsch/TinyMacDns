package cache

import (
	"fmt"
	"strings"
	"sync"
)

type Cache struct {
	dns  map[string]string
	rdns map[string]string
	lock sync.RWMutex
}

func New() *Cache {
	return &Cache{
		dns:  make(map[string]string),
		rdns: make(map[string]string),
	}
}

func (c *Cache) Update(hostname, ip string) {
	ihn := strings.ToLower(hostname)
	c.Delete(ihn)

	if revIp, revErr := reverseIP(ip); revErr == nil {
		c.lock.Lock()
		c.dns[ihn] = ip
		c.rdns[revIp] = ihn
		c.lock.Unlock()
	}
}

func (c *Cache) Delete(hostname string) {
	ihn := strings.ToLower(hostname)
	if oip, dok := c.dns[ihn]; dok {
		c.lock.Lock()
		delete(c.dns, ihn)
		revIp, _ := reverseIP(oip)
		if _, rok := c.rdns[revIp]; rok {
			delete(c.rdns, revIp)
		}
		c.lock.Unlock()
	}
}

func (c *Cache) GetIp(hostname string) (string, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	ip, ok := c.dns[hostname]
	return ip, ok
}

func (c *Cache) GetHostname(reverseIP string) (string, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	hostname, ok := c.rdns[reverseIP]
	return hostname, ok
}

func reverseIP(ip string) (string, error) {
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return fmt.Sprintf("%s.%s.%s.%s", parts[3], parts[2], parts[1], parts[0]), nil
	}
	return "", fmt.Errorf("invalid ip")
}
