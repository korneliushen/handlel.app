<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
</script>

<main class="flex min-h-screen flex-col items-center overflow-y-scroll font-sans">
	<Header />
	<slot />
</main>
