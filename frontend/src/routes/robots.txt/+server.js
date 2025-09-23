export const GET = () => {
    const body = `User-agent: *
Allow: /

Sitemap: https://hn.yamanlabs.com/sitemap.xml`;

    const headers = {
        'Content-Type': 'text/plain'
    };

    return new Response(body, { headers });
};
