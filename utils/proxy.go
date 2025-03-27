package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Default timeout for HTTP requests
var httpTimeout = 15 * time.Second

// ProxyManager handles proxy rotation for HTTP requests
type ProxyManager struct {
	proxies   []*url.URL
	mutex     sync.RWMutex
	lastIndex int
}

var (
	globalProxyManager *ProxyManager
	proxyOnce          sync.Once
)

// GetProxyManager returns the singleton proxy manager
func GetProxyManager() *ProxyManager {
	proxyOnce.Do(func() {
		globalProxyManager = &ProxyManager{
			proxies:   make([]*url.URL, 0),
			lastIndex: -1,
		}

		// Load proxies from environment variable if available
		if proxyList := os.Getenv("HTTP_PROXIES"); proxyList != "" {
			proxyStrs := strings.Split(proxyList, ",")
			for _, p := range proxyStrs {
				if proxyURL, err := url.Parse(strings.TrimSpace(p)); err == nil {
					globalProxyManager.AddProxy(proxyURL)
				}
			}
		}
	})

	return globalProxyManager
}

// GetCookieJar returns a cookie jar for HTTP requests
func GetCookieJar() (http.CookieJar, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}
	return jar, nil
}

// AddProxy adds a proxy to the manager
func (pm *ProxyManager) AddProxy(proxy *url.URL) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.proxies = append(pm.proxies, proxy)
}

// GetTransport returns an http.Transport with the next proxy
func (pm *ProxyManager) GetTransport() *http.Transport {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	// If no proxies, return default transport
	if len(pm.proxies) == 0 {
		return &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     20,
			IdleConnTimeout:     20 * time.Second,
			DisableCompression:  false,
		}
	}

	// Use a random proxy
	proxyIndex := rand.Intn(len(pm.proxies))
	proxy := pm.proxies[proxyIndex]

	return &http.Transport{
		Proxy:               http.ProxyURL(proxy),
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     20,
		IdleConnTimeout:     20 * time.Second,
		DisableCompression:  false,
	}
}

// GetClient returns an http.Client with a proxy-enabled transport
func (pm *ProxyManager) GetClient() (*http.Client, error) {
	// Create a cookie jar for the client
	jar, err := GetCookieJar()
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	// Create a client with the proxy transport
	client := &http.Client{
		Transport: pm.GetTransport(),
		Timeout:   httpTimeout,
		Jar:       jar,
	}

	return client, nil
}

// HasProxies returns true if proxies are configured
func (pm *ProxyManager) HasProxies() bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return len(pm.proxies) > 0
}
