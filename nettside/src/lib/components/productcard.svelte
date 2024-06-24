<script lang="ts">
	import { BookPlus, Trash2 } from 'lucide-svelte';
	import { handlelapp } from '$lib/stores/handlelapp';
	import type { ExtendedProduct } from '$lib/types/extendedPrisma';
	import { onMount } from 'svelte';
	export let product: ExtendedProduct;

	let imageError = false;
	let productIHandlelapp = false;
	let productIHandlelappId = 0;

	onMount(() => {
		for (let i = 0; i < $handlelapp.length; i++) {
			if ($handlelapp[i].id === product.id) {
				productIHandlelapp = true;
				productIHandlelappId = i;
			}
		}
	});

	function leggTilIHandlelapp() {
		for (let i = 0; i < $handlelapp.length; i++) {
			if ($handlelapp[i].id === product.id) {
				productIHandlelapp = true;
				productIHandlelappId = i;
			}
		}
		if (productIHandlelapp) {
			$handlelapp = [
				...$handlelapp.slice(0, productIHandlelappId),
				...$handlelapp.slice(productIHandlelappId + 1)
			];
			productIHandlelapp = false;
		} else {
			$handlelapp = [...$handlelapp, product];
			productIHandlelapp = true;
		}
	}
</script>

<div
	class="border-borderColor/26 relative flex h-full w-44 flex-col justify-between rounded-lg border p-2 shadow-sm transition hover:scale-[1.01] hover:border-mainPurple hover:shadow-md md:w-60"
>
	<div
		title="product"
		class="relative flex aspect-square max-h-40 items-center justify-center py-1 lg:w-full"
	>
		{#if imageError}
			<p class=" text-center text-xl text-gray-500">Det finnes ikke bilde for dette produktet</p>
		{:else}
			<img
				loading="lazy"
				class="max-h-40"
				src={product.images.small}
				on:error={() => (imageError = true)}
				alt="produktbilde"
			/>
		{/if}
		<img
			loading="lazy"
			class="absolute bottom-0 left-0 ml-2 w-10 rounded-md"
			src="/{product.prices[0].store}.svg"
			alt="nettside"
		/>
	</div>
	<div class="pb-1">
		<div class="mt-1 space-y-0 px-2">
			<p class="font-semibold">{product.title}</p>
			<p class="text-sm text-gray-500/60">{product.vendor}</p>
		</div>
		<div class="mt-1 flex justify-between px-2">
			<div class="w-fit space-y-0">
				<p class="font-semibold">{product.prices[0].price.toFixed(2)} kr</p>
				<p class="border-t border-gray-500/60 text-sm text-gray-500/60">
					{product.prices[0].unitprice.toFixed(2) || product.prices[0].price.toFixed(2)} kr/{product.unittype ||
						'stk'}
				</p>
			</div>
			{#if productIHandlelapp}
				<button
					on:click={() => leggTilIHandlelapp()}
					class="z-10 flex aspect-square h-10 w-10 items-center justify-center rounded-md bg-red-600 transition hover:brightness-110"
				>
					<Trash2 color="#fff" />
				</button>
			{:else}
				<button
					on:click={() => leggTilIHandlelapp()}
					class="z-10 flex aspect-square h-10 w-10 items-center justify-center rounded-md bg-mainPurple transition hover:brightness-110"
				>
					<BookPlus color="#fff" />
				</button>
			{/if}
		</div>
	</div>
	<a href="/produkt/{product.id}" class="absolute inset-0 z-0" title="Gå til produkt"></a>
</div>
