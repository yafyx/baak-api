# Avoiding 403 Forbidden Errors in BAAK API

This guide provides information on how the API now bypasses 403 Forbidden errors from the BAAK website.

## The Solution

After extensive testing, we implemented a multi-layered approach to bypass the 403 restrictions:

1. **Direct IP Access**: The API can now connect directly to the server's IP address (103.23.40.57) while setting the proper host header
2. **TLS Adjustments**: Modified the TLS settings to bypass certain security checks
3. **Protocol Fallbacks**: Tries HTTP before HTTPS to establish initial sessions
4. **Progressive Session Establishment**: Uses multiple fallback methods in sequence

## Built-in Protections

The API now implements these advanced anti-403 measures:

1. **Simplified Session Establishment**: Less complex headers initially to avoid detection
2. **Direct IP Connection**: Bypasses domain-based blocking
3. **HTTP/1.1 Forced**: Using older protocol to avoid HTTP/2 fingerprinting
4. **Disabled Compression**: Prevents detection via compression behaviors
5. **Advanced TLS Configuration**: Customized TLS settings to appear more browser-like
6. **Increased Timeouts**: Better handling of slow responses
7. **Proper Cookie Handling**: Ensures session persistence

## Testing Your Setup

To verify the anti-403 measures are working correctly, you can run the test client:

```bash
go run cmd/test/main.go
```

You should see output similar to:

```bash
Testing session establishment...
[DEBUG] No cookies found, trying to establish session
[DEBUG] Session established using HTTP request
Testing document fetch...
Trying URL: https://baak.gunadarma.ac.id
Document fetched successfully from: https://baak.gunadarma.ac.id
Page title: BAAK Online
```

If you see "Session established using..." in the output, the bypass is working correctly.

## Troubleshooting

If you're still experiencing 403 errors, the system logs which session establishment method was attempted:

- HTTP Request: The simple HTTP request method worked
- HTTPS Request: The HTTPS request method worked
- Direct IP method: The direct IP access method worked

### Common Issues

1. **DNS Issues**: If you can't connect to the domain, the direct IP access should still work
2. **Network Restrictions**: Some networks (especially corporate) may block direct IP access
3. **Proxy Configuration**: If using proxies, ensure they support custom headers

## Using Proxy Servers (Optional)

For additional protection, you can still configure proxy servers:

```bash
export HTTP_PROXIES="http://proxy1.example.com:8080,http://username:password@proxy2.example.com:3128"
```

The system will try the direct methods first before falling back to proxies.

## Additional Configuration

### Debug Logging

In production, you may want to disable the debug logging by setting this environment variable:

```bash
export DISABLE_DEBUG_LOGS=true
```

### Rate Limit Adjustment

To avoid triggering rate limits, you can adjust how many requests per minute the API makes:

```bash
export RATE_LIMIT_PER_MIN=30  # Default is 60
```

## Best Practices

1. **Use Caching**: Implement client-side caching to reduce requests to the API
2. **Preload Common Data**: Load common data (like calendar) once and reuse it
3. **Backoff on Errors**: Add retry logic with increasing delays when errors occur

## Technical Details

The new bypass approach works by:

1. First trying standard HTTP and HTTPS requests with minimal headers
2. If those fail, connecting directly to the IP address (103.23.40.57) rather than the domain
3. Setting the "Host" header to "baak.gunadarma.ac.id" on the IP-based request
4. Using a custom TLS configuration with insecure verification skipping
5. Disabling HTTP/2 and compression to appear more like an older browser

This combination of techniques has proven effective in bypassing the 403 restrictions.

## Legal Disclaimer

This API is designed for educational purposes only. Be aware that scraping websites may violate terms of service. Use responsibly and in accordance with applicable laws and regulations.
