import { ImageResponse } from '@vercel/og';
import { PRIVATE_API_BASE_URL } from '$env/static/private';
import { getDomain } from '$lib/utils';

// Generate a deterministic color from a string
function stringToColor(str) {
	let hash = 0;
	for (let i = 0; i < str.length; i++) {
		hash = str.charCodeAt(i) + ((hash << 5) - hash);
	}
	let color = '';
	for (let i = 0; i < 3; i++) {
		const value = (hash >> (i * 8)) & 0xff;
		color += ('00' + value.toString(16)).slice(-2);
	}
	return color;
}

// Generate a gradient from title string
function generateGradient(str) {
	const color1 = stringToColor(str);
	const color2 = stringToColor(str.split('').reverse().join(''));
	return `linear-gradient(135deg, #${color1}, #${color2})`;
}

// One day in seconds for cache headers
const ONE_DAY_SECONDS = 86400;

export async function GET({ url, fetch }) {
	const idParam = url.searchParams.get('id');

	if (!idParam) {
		return new Response('Missing id parameter', { status: 400 });
	}

	const id = parseInt(idParam, 10);
	if (isNaN(id) || id <= 0) {
		return new Response('Invalid id parameter', { status: 400 });
	}

	// Fetch top stories from /api/top to check if story is in top 30
	let stories;
	try {
		const topRes = await fetch(`${PRIVATE_API_BASE_URL}/api/top`);
		if (!topRes.ok) {
			return new Response('Failed to fetch top stories', { status: 500 });
		}
		stories = await topRes.json();
	} catch (error) {
		console.error('Failed to fetch top stories:', error);
		return new Response('Failed to fetch top stories', { status: 500 });
	}

	// Find the story with matching ID
	const story = stories.find((s) => s.id === id);
	if (!story) {
		return new Response('Story not in top 30', { status: 404 });
	}

	// Load fonts via HTTP from static directory (works in both dev and production)
	const origin = url.origin;
	let fontExtraBoldData, fontRegularData;
	try {
		const [fontExtraBoldRes, fontRegularRes] = await Promise.all([
			fetch(`${origin}/Inter-ExtraBold.ttf`),
			fetch(`${origin}/Inter-Regular.ttf`)
		]);

		if (!fontExtraBoldRes.ok || !fontRegularRes.ok) {
			throw new Error('Font fetch failed');
		}

		fontExtraBoldData = await fontExtraBoldRes.arrayBuffer();
		fontRegularData = await fontRegularRes.arrayBuffer();
	} catch (error) {
		console.error('Failed to load fonts:', error);
		return new Response('Failed to load fonts', { status: 500 });
	}

	const gradient = generateGradient(story.title);
	const domain = getDomain(story.url);

	// Create the OG image using @vercel/og
	const imageResponse = new ImageResponse(
		{
			type: 'div',
			props: {
				style: {
					height: '100%',
					width: '100%',
					display: 'flex',
					flexDirection: 'column',
					justifyContent: 'center',
					alignItems: 'flex-start',
					fontFamily: 'Inter',
					backgroundImage: gradient,
					padding: '60px'
				},
				children: {
					type: 'div',
					props: {
						style: {
							display: 'flex',
							flexDirection: 'column',
							alignItems: 'flex-start',
							maxWidth: '85%'
						},
						children: [
							{
								type: 'div',
								props: {
									style: {
										fontSize: 80,
										color: '#ffffff',
										textAlign: 'left',
										fontWeight: 900,
										lineHeight: 1.1,
										textShadow: '0 2px 10px rgba(0,0,0,0.3)',
										marginBottom: '20px'
									},
									children: story.title
								}
							},
							domain && {
								type: 'div',
								props: {
									style: {
										fontSize: 28,
										color: 'rgba(255,255,255,0.85)',
										textShadow: '0 1px 4px rgba(0,0,0,0.2)'
									},
									children: domain
								}
							}
						].filter(Boolean)
					}
				}
			}
		},
		{
			width: 1200,
			height: 630,
			fonts: [
				{
					name: 'Inter',
					data: fontRegularData,
					style: 'normal',
					weight: 400
				},
				{
					name: 'Inter',
					data: fontExtraBoldData,
					style: 'normal',
					weight: 900
				}
			]
		}
	);

	// Get the response body and create a new response with cache headers
	const body = await imageResponse.arrayBuffer();

	return new Response(body, {
		status: 200,
		headers: {
			'Content-Type': 'image/png',
			'Cache-Control': `public, max-age=${ONE_DAY_SECONDS}, s-maxage=${ONE_DAY_SECONDS}, stale-while-revalidate=${ONE_DAY_SECONDS}`
		}
	});
}
