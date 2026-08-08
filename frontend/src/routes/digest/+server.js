import { PRIVATE_API_BASE_URL } from '$env/static/private';

/** Serves the email-style digest HTML at /digest (no site chrome). */
export const GET = async ({ fetch }) => {
	const res = await fetch(`${PRIVATE_API_BASE_URL}/api/newsletter/digest`);
	const body = await res.text();

	return new Response(body, {
		status: res.status,
		headers: {
			'Content-Type': res.headers.get('Content-Type') || 'text/html; charset=utf-8',
			'Cache-Control': res.headers.get('Cache-Control') || 'public, max-age=120'
		}
	});
};
