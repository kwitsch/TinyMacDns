package cache

import (
	"fmt"
	"strings"
)

type Cache struct {
	dns  map[string]string
	rdns map[string]string
}

func New() *Cache {
	return &Cache{
		dns:  make(map[string]string),
		rdns: make(map[string]string),
	}
}

func (c *Cache) Update(hostname, ip string) {
	c.Delete(hostname)
	if revIp, revErr := reverseIP(ip); revErr == nil {
		c.dns[hostname] = ip
		c.rdns[revIp] = hostname
	}
}

func (c *Cache) Delete(hostname string) {
	if oip, dok := c.dns[hostname]; dok {
		delete(c.dns, hostname)
		revIp, _ := reverseIP(oip)
		if _, rok := c.rdns[revIp]; rok {
			delete(c.rdns, revIp)
		}
	}
}

func (c *Cache) GetIp(hostname string) (string, bool) {
	ip, ok := c.dns[hostname]
	return ip, ok
}

func (c *Cache) GetHostname(reverseIP string) (string, bool) {
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
