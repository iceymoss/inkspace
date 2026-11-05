package seo

// RenderRobots describes blog-origin crawl policy. sitemapURL must be generated
// from the validated canonical origin, never request headers or private records.
// Admin is a separate application/origin; its API is covered by /api/ here.
// These generic application paths are guidance, not authentication or secrecy.
func RenderRobots(sitemapURL string) []byte {
	return []byte("User-agent: *\n" +
		"Allow: /\n" +
		"Disallow: /api/\n" +
		"Disallow: /dashboard\n" +
		"Disallow: /favorites\n" +
		"Disallow: /profile\n" +
		"Disallow: /login\n" +
		"Disallow: /uploads/knowledge\n" +
		"\nSitemap: " + sitemapURL + "\n")
}
