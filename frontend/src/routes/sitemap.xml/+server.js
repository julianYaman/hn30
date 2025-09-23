export async function GET() {
    const baseUrl = 'https://hn.yamanlabs.com';
    const pages = ['/', '/privacy', '/bookmarks'];

    const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    ${pages.map(page => `
    <url>
        <loc>${baseUrl}${page}</loc>
        <changefreq>daily</changefreq>
        <priority>${page === '/' ? '1.0' : '0.7'}</priority>
    </url>
    `).join('')}
</urlset>`;

    const headers = {
        'Content-Type': 'application/xml',
        'Cache-Control': 'public, max-age=86400' // Cache for a day
    };

    return new Response(sitemap, { headers });
}
