<script>
	import { notifications } from '$lib/stores/notifications.js';

	async function enableNotifications() {
		if (!$notifications.cookiesAccepted) return;
		if ($notifications.permission === 'denied') return;
		await notifications.setNotificationsEnabled(true);
	}

	$: disabled =
		$notifications.loading ||
		$notifications.permission === 'denied' ||
		!$notifications.cookiesAccepted;

	$: buttonLabel = $notifications.loading
		? $notifications.loadingMessage || 'Working…'
		: 'Enable notifications';
</script>

{#if !$notifications.wantsNotifications}
	<section class="sidebar-card sidebar-notifications">
		<h2>Push notifications</h2>
		<p>Get alerted when a story hits the top of Hacker News.</p>

		<button
			type="button"
			class="sidebar-notifications__button"
			on:click={enableNotifications}
			{disabled}
		>
			{buttonLabel}
		</button>

		{#if !$notifications.cookiesAccepted}
			<p class="sidebar-notifications__status">Accept cookies first to enable notifications.</p>
		{:else if $notifications.permission === 'denied'}
			<p class="sidebar-notifications__status">
				Notifications are blocked in your browser settings.
			</p>
		{:else if $notifications.error}
			<p class="sidebar-notifications__status is-error">{$notifications.error}</p>
		{:else if $notifications.message}
			<p
				class="sidebar-notifications__status"
				class:is-success={$notifications.messageType === 'success'}
			>
				{$notifications.message}
			</p>
		{/if}
	</section>
{/if}
