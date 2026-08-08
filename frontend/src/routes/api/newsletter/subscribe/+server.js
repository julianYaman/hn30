import { PRIVATE_API_BASE_URL } from '$env/static/private';
import { json } from '@sveltejs/kit';

export const POST = async ({ request }) => {
	let body;
	try {
		body = await request.json();
	} catch {
		return json({ error: 'Invalid JSON body' }, { status: 400 });
	}

	const email = typeof body?.email === 'string' ? body.email.trim() : '';
	if (!email) {
		return json({ error: 'Email is required' }, { status: 400 });
	}

	const res = await fetch(`${PRIVATE_API_BASE_URL}/api/newsletter/subscribe`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email })
	});

	const text = await res.text();
	if (!res.ok) {
		return new Response(text || 'Subscribe failed', {
			status: res.status,
			headers: { 'Content-Type': 'text/plain' }
		});
	}

	return new Response(text || JSON.stringify({ status: 'ok' }), {
		status: 200,
		headers: { 'Content-Type': 'application/json' }
	});
};
