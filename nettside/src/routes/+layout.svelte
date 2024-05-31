<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	import { handlelapp } from '$lib/stores/handlelapp';
	import { onMount } from 'svelte';
	import type { products } from '@prisma/client';
	import { browser } from '$app/environment';

	onMount(() => {
		$handlelapp = JSON.parse(localStorage.getItem('handlelapp') as string) as products[];
	});
</script>

<main class="flex min-h-screen flex-col items-center overflow-y-scroll font-sans">
	<Header />
	<slot />
</main>
