<script lang="ts">
	import Productcard from '$lib/components/productcard.svelte';
	export let data: import('./$types').PageData;
</script>

{#if data.products}
	<div class="space-y-3 p-4">
		<h2 class="text-start text-xl font-medium">
			{data.products?.length}
			{data.products?.length != 1 ? 'resultater' : 'resultat'} for "{data.param}"
		</h2>
		<section class="grid grid-cols-2 gap-[16px] md:grid-cols-3 lg:grid-cols-4">
			{#each data.products as product}
				<Productcard {product} />
			{/each}
		</section>
	</div>
{:else}
  <p class="p-4 text-start text-xl font-medium">Søk etter et produkt</p>
{/if}
