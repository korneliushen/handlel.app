<script lang="ts">
	import { BookPlus } from 'lucide-svelte';
	import type { products } from '@prisma/client';
	import { handlelapp } from '$lib/stores/handlelapp';
	export let product: products;
</script>

<div
	class=" border-borderColor/26 z-20 flex h-full w-44 flex-col justify-between rounded-lg border p-2 shadow-sm transition hover:scale-[1.01] hover:border-purple-500 hover:shadow-md md:w-60"
>
	<a
		href="/produkt/{product.id}"
		title="product"
		class=" relative flex aspect-square max-h-40 items-center justify-center py-1 lg:w-full"
	>
		<img class=" max-h-40" loading="lazy" src={product.imagelink} alt="produktbilde" />
		<img
			loading="lazy"
			class=" absolute bottom-0 left-0 ml-2 w-10 rounded-md"
			src="/{product.prices[0].Store}.svg"
			alt="nettside"
		/>
	</a>
	<div class=" pb-1">
		<a href="/produkt/{product.id}" title="product">
			<div class=" mt-1 space-y-0 px-2">
				<p class=" font-semibold">{product.title}</p>
				<p class=" text-sm text-gray-500/60">{product.vendor}</p>
			</div>
		</a>
		<div class=" mt-1 flex justify-between px-2">
			<a class=" w-full" href="/produkt/{product.id}" title="product">
				<div class=" w-fit space-y-0">
					<p class=" font-semibold">{product.prices[0].Price + ' kr'}</p>
					<p class=" border-t border-gray-500/60 text-sm text-gray-500/60">
						{product.prices[0].UnitPrice} kr/{product.unittype || 'N/A'}
					</p>
				</div>
			</a>
			<button
				on:click={() => [($handlelapp = [...$handlelapp, product])]}
				class=" z-50 flex aspect-square h-10 w-10 items-center justify-center rounded-md bg-purple-500 transition hover:brightness-110"
			>
				<BookPlus color="#fff" />
			</button>
		</div>
	</div>
</div>
