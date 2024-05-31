<script lang="ts">
	import { ChevronDown, Trash2, Minus, Plus } from 'lucide-svelte';
	import autoAnimate from '@formkit/auto-animate';
	import { onMount } from 'svelte';
	let antallProdukter = 0;
	let dropdown = false;
	let lesmer = '13rem';
	let lesmerBool = false;

	onMount(() => {
		if (window.innerWidth > 1024) {
			dropdown = true;
		}
	});

	function lesmerFunc() {
		lesmerBool = true;
		if (lesmerBool) {
			lesmer = 'fit-content';
		} else {
			lesmer = '10rem';
		}
	}
	export let data: import('./$types').PageData;
</script>

<main
	class=" relative grid w-screen max-w-[600px] grid-cols-1 gap-x-8 gap-y-4 overflow-hidden rounded-lg px-5 lg:max-w-[1200px] lg:grid-cols-2 lg:px-20"
>
	<div class=" relative flex aspect-square min-h-80 items-center justify-center py-1 lg:w-full">
		<img src={data.product.imagelinkmedium} alt="Produktbilde" />
	</div>
	<div class=" flex flex-col lg:relative">
		<div class=" mt-5">
			<a target="_blank" href={data.product?.prices[0].url} class=" text-2xl font-bold"
				>{data.product?.title}</a
			>
			<div class=" mt-2 flex justify-between">
				<p class=" text-lg text-gray-500/60">{data.product?.vendor}</p>
				<div class=" flex items-center">
					<a target="_blank" href={data.product?.prices[0].url}
						><img
							class=" mr-4 h-12 rounded-md"
							src="/{data.product.prices[0].store}.svg"
							alt="nettside"
						/></a
					>
					<div class=" text-end">
						<p class=" text-2xl font-bold text-mainPurple">{data.product?.prices[0].price} kr</p>
						<p class=" text-lg text-gray-500/60">
							{data.product?.prices[0].unitprice || data.product?.prices[0].price} kr/{data.product
								?.unittype || 'stk'}
						</p>
					</div>
				</div>
			</div>
		</div>
		<div class=" mt-3 flex flex-col items-center">
			<p class=" text-xl font-bold">Andre butikker</p>
			{#each data.product?.prices.slice(1) as price}
				<a
					href={price.url}
					target="_blank"
					class=" my-1 flex w-full justify-between rounded-md border border-borderColor p-2"
				>
					<div class=" flex w-12 items-center">
						<img class=" mr-3 w-full rounded-md" src="/{price.store}.svg" alt="Butikklogo" />
						<p class=" font-bold">{price.store[0].toUpperCase() + price.store.substring(1)}</p>
					</div>
					<div class=" text-end">
						<p class=" text-lg font-bold">{price.price} kr</p>
						<p class=" text-sm text-gray-500/60">
							{price.unitprice || price.prices[0].price} kr/{data.product.unittype || 'stk'}
						</p>
					</div>
				</a>
			{/each}
		</div>
		<div
			class=" fixed bottom-16 left-0 z-50 flex h-20 w-full items-center justify-center border-t border-borderColor bg-white lg:absolute lg:bottom-0 lg:border-none"
		>
			<div class=" h-12 w-4/5 rounded-lg lg:w-full">
				{#if antallProdukter !== 0}
					<div
						class=" flex h-full w-full items-center justify-between rounded-lg border border-mainPurple"
					>
						{#if antallProdukter !== 1}
							<button
								on:click={() => antallProdukter--}
								class=" flex w-14 justify-center text-xl font-extrabold"><Minus /></button
							>
						{:else}
							<button
								on:click={() => antallProdukter--}
								class=" flex w-14 justify-center text-xl font-extrabold"><Trash2 /></button
							>
						{/if}
						<p>{antallProdukter}</p>
						<button
							on:click={() => antallProdukter++}
							class=" flex w-14 justify-center text-xl font-extrabold"><Plus /></button
						>
					</div>
				{:else}
					<button
						on:click={() => antallProdukter++}
						class=" flex h-full w-full items-center justify-center rounded-lg border border-mainPurple bg-mainPurple"
					>
						<p class=" font-bold text-white">Legg til i handlelisten</p>
						<img src="" alt="" />
					</button>
				{/if}
			</div>
		</div>
	</div>
	<div
		use:autoAnimate
		class=" relative overflow-hidden border-t border-borderColor lg:w-full"
		style="max-height: {lesmer};"
	>
		{#if !lesmerBool}
			<div class=" absolute z-10 h-full w-full bg-gradient-to-b from-transparent to-white to-85%" />
		{/if}
		<div>
			<p class=" py-1 pt-6 text-xl font-bold">Om produktet</p>
		</div>
		<div>
			<p class=" mb-7">
				{data.product.description || 'Dette produktet har ingen beskrivelse'}
			</p>
			<div class=" my-4 flex justify-between">
				<p class=" text-lg font-bold">Mengde</p>
				<p>{data.product.weight || 'N/A'}</p>
			</div>
			<div class=" my-4 flex w-full justify-between">
				<p class=" text-lg font-bold">Ingredienser</p>
				<p class=" w-2/3 text-end">{data.product.ingredients || 'N/A'}</p>
			</div>
			<div class=" my-4 flex w-full justify-between">
				<p class=" text-lg font-bold">Allergener</p>
				<p class=" w-2/3 text-end">{data.product.allergens || 'N/A'}</p>
			</div>
			<div class=" my-4 flex w-full justify-between">
				<p class=" text-lg font-bold">Opprinnelsesland</p>
				<p class=" w-2/3 text-end">{data.product.origincountry || 'N/A'}</p>
			</div>
			<div class=" my-4 flex w-full justify-between">
				<p class=" text-lg font-bold">Produsent</p>
				<p class=" w-2/3 text-end">{data.product.vendor || 'N/A'}</p>
			</div>
		</div>
		{#if !lesmerBool}
			<button
				on:click={() => lesmerFunc()}
				class=" absolute bottom-1 left-1 z-20 font-bold text-mainPurple">+ Les mer...</button
			>
		{/if}
	</div>
	<div use:autoAnimate class=" mb-24 mt-3 border-y border-borderColor lg:mt-0 lg:h-fit lg:w-full">
		<button
			on:click={() => (dropdown = !dropdown)}
			class=" flex h-20 w-full items-center justify-between py-3"
		>
			<p class=" py-1 text-xl font-bold">Næringsinnhold per 100g</p>
			<ChevronDown class=" {dropdown ? 'rotate-180' : 'rotate-0'}" strokeWidth={3} />
		</button>
		{#if dropdown}
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Kalorier</p>
				<p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.calories || 'N/A'}</p>
			</div>
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Energi</p>
				<p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.energy || 'N/A'}</p>
			</div>
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Fett</p>
				<p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.fat || 'N/A'}</p>
			</div>
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Karbohydrater</p>
				<p class=" w-2/3 text-end font-bold">
					{data.product.nutritionalcontent.carbohydrates || 'N/A'}
				</p>
			</div>
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Kostfiber</p>
				<p class=" w-2/3 text-end font-bold">
					{data.product.nutritionalcontent.dietaryfiber || 'N/A'}
				</p>
			</div>
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Mettet fett</p>
				<p class=" w-2/3 text-end font-bold">
					{data.product.nutritionalcontent.saturatedfat || 'N/A'}
				</p>
			</div>
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Protein</p>
				<p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.protein || 'N/A'}</p>
			</div>
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Salt</p>
				<p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.salt || 'N/A'}</p>
			</div>
			<div class=" flex w-full justify-between border-t border-borderColor p-2">
				<p>Sukkerarter</p>
				<p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.sugars || 'N/A'}</p>
			</div>
		{/if}
	</div>
</main>
