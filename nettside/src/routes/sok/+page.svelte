<script lang="ts">
	import Productcard from '$lib/components/productcard.svelte';
	import { Search } from 'lucide-svelte';
	export let data: import('./$types').PageData;
</script>

<div class="visible mx-4 mb-2 mt-4 flex w-full items-center justify-center md:invisible">
	<search class=" w-80 flex-1">
		<form
			class="border-borderColor/26 flex w-full items-center justify-between rounded-md border px-2 text-borderColor shadow-sm"
			action="/sok"
		>
			<input
				type="text"
				placeholder="Søk etter produkter..."
				class="flex-1 bg-transparent py-2 text-black outline-none"
				name="search"
				autocomplete="off"
			/>
			<Search />
		</form>
	</search>
</div>
{#if data.products}
	<div class="space-y-3 p-4">
		<h2 class="text-start text-xl font-medium">
			{data.products?.length}
			{data.products?.length != 1 ? 'resultater' : 'resultat'} for "{data.param}"
		</h2>
		<section class="grid grid-cols-2 gap-[16px] pb-14 md:grid-cols-3 lg:grid-cols-4">
			{#each data.products as product}
				<Productcard {product} />
			{/each}
		</section>
	</div>
{:else}
	<p class="p-4 text-start text-xl font-medium">Søk etter et produkt</p>
{/if}
