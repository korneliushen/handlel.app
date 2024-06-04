<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
</script>

<svelte:head>
	<title>Handlelapp</title>
	<meta
		name="description"
		content="Få full kontroll over matbudsjettet ditt med HandleLapp. Sammenlign priser fra alle de store norske kjedene, legg varer til handlelister og finn den billigste butikken for deg."
	/>
</svelte:head>
<main
	class="absolute inset-0 flex min-h-screen flex-col items-center overflow-x-hidden overflow-y-scroll font-sans"
>
	<Header />
	<div class=" mt-1 sm:mt-0">
		<slot />
	</div>
</main>
