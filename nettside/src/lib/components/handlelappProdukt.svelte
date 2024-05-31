<script lang="ts">
	import { Minus, Plus, X } from 'lucide-svelte';
	import { onMount } from 'svelte';
	export let product: import('@prisma/client').products;
	export let id: number;

	let desktopView = false;
	let antallProdukt = 1;

	onMount(() => {
		if (window.innerWidth > 1024) {
			desktopView = true;
		}
	});
</script>

<div class=" mt-6">
	<div class=" mt-5 flex h-20 items-center justify-between rounded-xl border border-gray-400">
		<div class=" flex h-full items-center">
			<img class=" w-16" src={product.imagelinkxsmall} alt="" />
			<div>
				<p class=" text-sm font-medium">{product.title}</p>
				<p class=" text-xs text-gray-400">{product.brand || product.vendor}</p>
				{#if !desktopView}
					<p class=" text-sm font-medium">{product.prices[0].price} kr</p>
				{/if}
			</div>
		</div>
		<div class=" mr-4 flex h-full items-center">
			{#if desktopView}
				<div class=" mr-3 text-end">
					<p class=" font-medium">{product.prices[0].price} kr</p>
					<p class=" text-xs font-medium text-gray-400">
						{product.prices[0].unitprice || product.prices[0].price} kr/{product.unittype || 'stk'}
					</p>
				</div>
			{/if}
			<div
				class=" border-mainPurple rounded-md border {desktopView
					? 'h-10 w-32'
					: ' h-9 w-28'} flex items-center justify-between"
			>
				<button on:click={() => antallProdukt--} class=" flex w-10 items-center justify-center"
					><Minus size="15px" /></button
				>
				<p>{antallProdukt}</p>
				<button on:click={() => antallProdukt++} class=" flex w-10 items-center justify-center"
					><Plus size="15px" /></button
				>
			</div>
			<button on:click={() => removeProduct(product, id)} class=" ml-3 text-gray-400"
				><X size="15px" /></button
			>
		</div>
	</div>
</div>
