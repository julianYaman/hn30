<script>
	import { onMount } from 'svelte';

	const MAILJET_FORM_SRC =
		'https://1orpm.mjt.lu/wgt/1orpm/08l9/form?c=58f63092';
	const MAILJET_EMBED_SCRIPT = 'https://app.mailjet.com/pas-nc-embedded-v2.js';

	let formFrame;

	onMount(() => {
		const initializeForm = () => {
			if (formFrame && window.iFrameResize && !formFrame.iFrameResizer) {
				window.iFrameResize({ checkOrigin: false }, formFrame);
			}
		};

		if (window.iFrameResize) {
			initializeForm();
			return;
		}

		const existingScript = document.querySelector(`script[src="${MAILJET_EMBED_SCRIPT}"]`);
		if (existingScript) {
			existingScript.addEventListener('load', initializeForm, { once: true });
			return () => existingScript.removeEventListener('load', initializeForm);
		}

		const script = document.createElement('script');
		script.src = MAILJET_EMBED_SCRIPT;
		script.async = true;
		script.addEventListener('load', initializeForm, { once: true });
		document.body.appendChild(script);

		return () => {
			script.removeEventListener('load', initializeForm);
		};
	});
</script>

<svelte:head>
	<title>Subscribe · Daily dispatch · hn30</title>
	<meta
		name="description"
		content="Subscribe to the hn30 Daily dispatch — top 30 Hacker News stories in your inbox each morning."
	/>
	<link rel="canonical" href="https://hn30.eu/subscribe" />
</svelte:head>

<main class="subscribe-page">
	<header class="subscribe-page__header">
		<p class="subscribe-page__kicker">Daily dispatch</p>
		<h1>Subscribe</h1>
		<p class="subscribe-page__lede">
			The top thirty Hacker News stories, once a day. After you sign up, check your inbox to confirm
			your subscription.
		</p>
		<p class="subscribe-page__back"><a href="/">← Back to hn30</a></p>
	</header>

	<div class="subscribe-page__form">
		<iframe
			bind:this={formFrame}
			data-w-type="embedded"
			title="hn30 Daily dispatch subscription form"
			frameborder="0"
			scrolling="no"
			marginheight="0"
			marginwidth="0"
			src={MAILJET_FORM_SRC}
			width="100%"
			style="height: 0;"
		></iframe>
	</div>

	<p class="subscribe-page__privacy">
		By subscribing you agree to our <a href="/privacy">Privacy Policy</a>. You can unsubscribe at any
		time from any email.
	</p>
</main>

<style>
	.subscribe-page {
		max-width: 32rem;
		margin: 0 auto;
		padding: 2rem 1.25rem 4rem;
	}

	.subscribe-page__header {
		text-align: center;
		margin-bottom: 1.75rem;
		padding-bottom: 1.5rem;
		border-bottom: 2px solid var(--rule, #232321);
	}

	.subscribe-page__kicker {
		margin: 0;
		font: 500 12px var(--font-mono, monospace);
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--subtle, #99958d);
	}

	.subscribe-page__header h1 {
		margin: 0.4rem 0 0.65rem;
		font-family: var(--font-serif, Georgia, serif);
		font-size: clamp(1.75rem, 4vw, 2.25rem);
		font-weight: 700;
		color: var(--ink, #1d1d1b);
	}

	.subscribe-page__lede,
	.subscribe-page__back,
	.subscribe-page__privacy {
		margin: 0;
		font-family: var(--font-sans, system-ui, sans-serif);
		font-size: 0.95rem;
		line-height: 1.55;
		color: var(--muted, #66645f);
	}

	.subscribe-page__back {
		margin-top: 0.85rem;
	}

	.subscribe-page__back a,
	.subscribe-page__privacy a {
		color: var(--orange, #ff6600);
		text-decoration: none;
	}

	.subscribe-page__back a:hover,
	.subscribe-page__privacy a:hover {
		text-decoration: underline;
	}

	.subscribe-page__form {
		min-height: 12rem;
	}

	.subscribe-page__privacy {
		margin-top: 1.5rem;
		font-size: 0.85rem;
		text-align: center;
	}
</style>
