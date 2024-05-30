<script lang="ts">
	import { BookPlus } from 'lucide-svelte';
	import type { products } from '@prisma/client';
	export let product: products;
</script>


<div
	class="border-borderColor/26 relative flex h-full w-44 flex-col justify-between rounded-lg border p-2 shadow-sm transition hover:scale-[1.01] hover:border-mainPurple hover:shadow-md md:w-60"
>
	<div title="product" class="relative flex aspect-square max-h-40 justify-center py-1 lg:w-full items-center">
		<img loading="lazy" class="max-h-40" src={product.imagelinksmall} alt="produktbilde" />
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
				<p class="font-semibold">{product.prices[0].price} kr</p>
				<p class="border-t border-gray-500/60 text-sm text-gray-500/60">
					{product.prices[0].unitprice || product.prices[0].price} kr/{product.unittype || 'stk'}
				</p>
			</div>
			<button
				on:click={() => alert('added to handlelapp')}
				class="z-10 flex aspect-square h-10 w-10 items-center justify-center rounded-md bg-mainPurple transition hover:brightness-110"
			>
				<BookPlus color="#fff" />
			</button>
		</div>
	</div>
	<a
		href="/produkt/{product.id}"
		class="absolute inset-0 z-0"
		title="Gå til produkt"
		target="_blank"
	></a>
</div>
