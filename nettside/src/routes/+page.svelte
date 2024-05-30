<script lang="ts">
	import { ChevronRight } from 'lucide-svelte';
	import Productcard from '$lib/components/productcard.svelte';
	export let data: import('./$types').PageData;
	console.log(data.products);
</script>

<div class="my-6 flex items-center justify-center p-4 w-full">
	<img
		class="flex w-full items-center justify-center rounded-md h-[256px] md:w-[1200px] md:h-96"
		src="http://1.vgc.no/vgnett-prod/img/vgLogoSquare.png?28042014-1"
		alt="Banner"
	/>
</div>
<section class="flex flex-col gap-y-14">

	{#each Object.entries(data.products) as [category, products]}
		<div class="space-y-1">
			<a href="#" title="Gå til kategori" class="flex items-center">
				<h2 class="text-xl font-medium">{category}</h2>
				<ChevronRight color="#7A38D0" />
			</a>
			<div class="grid grid-cols-2 gap-[16px] md:grid-cols-3 lg:grid-cols-4">
				{#each products as product}
					<Productcard {product} />
				{/each}
			</div>
		</div>
	{/each}
</section>
