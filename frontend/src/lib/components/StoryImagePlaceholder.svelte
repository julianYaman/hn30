<script>
	import { mutedPreviewTone } from '$lib/utils.js';
	export let title = '';
	export let size = 'thumb';
	$: tone = mutedPreviewTone(title);
</script>

<div
	class="image-placeholder image-placeholder--title preview-tone-{tone}"
	class:image-placeholder--hero={size === 'hero'}
	aria-hidden="true"
>
	<p class="image-placeholder__title">{title}</p>
</div>
