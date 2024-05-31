<script lang="ts">
	import Filter from '$lib/components/filter.svelte';
	import HandlelappProdukt from '$lib/components/handlelappProdukt.svelte';
	import { filter } from '$lib/stores/filter';
	import { onMount } from 'svelte';
	import { ChevronDown, SlidersHorizontal, MapPin } from 'lucide-svelte';
	import { handlelapp } from '$lib/stores/handlelapp';

	let desktopView = false;
	let selected = false;
	onMount(() => {
		if (window.innerWidth > 1024) {
			$filter = true;
			desktopView = true;
		}
	});
</script>

<main class=" relative flex w-screen max-w-[400px] px-4 lg:max-w-[1200px]">
	<Filter />
	<div class=" w-full {desktopView ? 'border-l border-gray-200 pl-8' : ''}">
		<div class=" flex h-16 w-full">
			<button
				on:click={() => (selected = false)}
				class=" flex h-full flex-1 cursor-pointer flex-col px-1 {!selected && 'font-semibold'}"
			>
				<div class=" flex h-full w-full flex-col justify-between px-2">
					<div class=" flex justify-between">
						<p class=" text-black {!selected && '!font-bold !text-mainPurple'}">Billigste</p>
						<img class=" w-7" src="/favicon.png" alt="" />
					</div>
					<div class=" flex justify-between pb-1 text-sm">
						<p>10 min</p>
						<p>27.60 kr</p>
					</div>
				</div>
				<div class="flex h-1 w-full items-end">
					<div
						class=" h-[0.16rem] w-full rounded-t-2xl bg-gray-400 {!selected &&
							'!h-1 !bg-mainPurple'}"
					/>
				</div>
			</button>
			<button
				on:click={() => (selected = true)}
				class=" flex h-full flex-1 cursor-pointer flex-col px-1 {selected && 'font-semibold'}"
			>
				<div class=" flex h-full w-full flex-col justify-between px-2">
					<div class=" flex justify-between">
						<p class=" text-black {selected && '!font-bold !text-mainPurple'}">Raskeste</p>
						<img class=" w-7" src="/favicon.png" alt="" />
					</div>
					<div class=" flex justify-between pb-1 text-sm">
						<p>8 min</p>
						<p>29.90 kr</p>
					</div>
				</div>
				<div class="flex h-1 w-full items-end">
					<div
						class=" h-[0.16rem] w-full rounded-t-2xl bg-gray-400 {selected &&
							'!h-1 !bg-mainPurple'}"
					/>
				</div>
			</button>
		</div>
		<div class=" mt-3 flex w-full justify-between px-2">
			<div>
				<p class=" text-2xl">Handlelapp</p>
				<p class=" mt-2 flex text-sm">
					Sorter etter:<button class=" mx-1 flex items-center font-bold text-mainPurple"
						>Pris <ChevronDown size="20px" /></button
					>
				</p>
			</div>
			{#if !desktopView}
				<div class=" mt-2 flex">
					<button
						on:click={() => ($filter = true)}
						class=" mx-[0.13rem] flex h-9 w-9 items-center justify-center rounded-md bg-mainPurple"
						><SlidersHorizontal color="#ffffff" /></button
					>
					<button
						class=" mx-[0.13rem] flex h-9 w-9 items-center justify-center rounded-md bg-mainPurple"
						><MapPin color="#ffffff" /></button
					>
				</div>
			{/if}
		</div>
		{#each $handlelapp as product, id}
			<HandlelappProdukt {product} {id} />
		{/each}
	</div>
</main>
