import { PRIVATE_API_BASE_URL } from '$env/static/private';

export const GET = async ({ request }) => {
    const searchParams = new URL(request.url).searchParams;
    const imageUrl = searchParams.get('url');

    if (!imageUrl) {
        return new Response('Missing image URL', { status: 400 });
    }

    // Construct the full internal URL to the Go backend's image proxy
    const backendProxyUrl = `${PRIVATE_API_BASE_URL}/api/image-proxy?url=${encodeURIComponent(imageUrl)}`;

    try {
        // Fetch the image from the Go backend
        const res = await fetch(backendProxyUrl);

        // If the backend returned an error, we pass it along
        if (!res.ok) {
            return new Response(res.body, {
                status: res.status,
                statusText: res.statusText,
                headers: res.headers, // Pass through backend's error headers
            });
        }

        // The response body from the Go backend is a ReadableStream.
        // We can pass it directly into a new Response to efficiently stream
        // the image to the browser without buffering it all in memory.
        return new Response(res.body, {
            status: 200,
            headers: {
                // Forward the essential headers from the original image response
                'Content-Type': res.headers.get('Content-Type') || 'application/octet-stream',
                'Content-Length': res.headers.get('Content-Length') || '',
                'Cache-Control': res.headers.get('Cache-Control') || 'public, max-age=31536000, immutable',
            }
        });

    } catch (error) {
        console.error('Error proxying image:', error);
        return new Response('Failed to proxy image', { status: 500 });
    }
};
