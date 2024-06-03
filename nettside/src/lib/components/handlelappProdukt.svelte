<script lang="ts">
	import { Minus, Plus, X } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { handlelapp } from '$lib/stores/handlelapp';
	import type { ExtendedProduct } from '$lib/types/extendedPrisma';
	export let product: ExtendedProduct;
	export let id: number;

	let desktopView = false;
	let antallProdukt = 1;

	onMount(() => {
		if (window.innerWidth > 1024) {
			desktopView = true;
		}
	});

	const handleError = (e: any) => e.target.src = "/handlelapp.png"
</script>

<div class=" mt-6">
	<div class=" mt-5 flex h-20 items-center justify-between rounded-xl border border-gray-400">
		<div class=" ml-4 flex h-full items-center gap-4">
			<div class="flex w-16 items-center justify-center">
				<img class=" max-h-16 max-w-16" src={product.imagelinkxsmall} on:error={handleError} alt="produktbilde" />
			</div>
			<div>
				<p class=" text-xs md:text-sm font-medium">{product.title}</p>
				<p class=" text-xs text-gray-400">{product.brand || product.vendor}</p>
				{#if !desktopView}
					<p class=" text-xs md:text-sm font-medium">{product.prices[0].price} kr</p>
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
				class=" rounded-md border border-mainPurple {desktopView
					? 'h-10 w-32'
					: ' h-7 w-20'} flex items-center justify-between"
			>
				<button on:click={() => antallProdukt--} class=" flex w-10 items-center justify-center"
					><Minus size={desktopView ? '15px' : '10px'} /></button
				>
				<p class=" text-sm md:text-base">{antallProdukt}</p>
				<button on:click={() => antallProdukt++} class=" flex w-10 items-center justify-center"
					><Plus size={desktopView ? '15px' : '10px'} /></button
				>
			</div>
			<button
				on:click={() => ($handlelapp = [...$handlelapp.slice(0, id), ...$handlelapp.slice(id + 1)])}
				class=" ml-3 text-gray-400"><X size="15px" /></button
			>
		</div>
	</div>
</div>
