<script lang="ts">
	import { ChevronRight } from 'lucide-svelte';
	import Productcard from '$lib/components/productcard.svelte';
	export let data: import('./$types').PageData;
</script>

<div class="my-6 flex w-full items-center justify-center p-4">
	<div
		class="flex bg-gray-300 h-[256px] w-full items-center justify-center rounded-md md:h-96 md:w-[1200px]"
	/>
</div>
<section class="flex flex-col items-center gap-y-14">
	{#each Object.entries(data.products) as [category, products]}
		<div class="space-y-1 p-1 md:p-4">
			<a href="#" title="Gå til kategori" class="flex items-center">
				<h2 class="text-xl font-medium">{category}</h2>
				<ChevronRight color="#7A38D0" />
			</a>
			<div class="grid grid-cols-2 gap-x-4 gap-y-4 md:grid-cols-3 md:gap-x-8 lg:grid-cols-4">
				{#each products as product}
					<Productcard {product} />
				{/each}
			</div>
		</div>
	{/each}
</section>
