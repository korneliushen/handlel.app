<script lang="ts">
	import { ChevronRight, LoaderCircle } from 'lucide-svelte';
	import Productcard from '$lib/components/productcard.svelte';
	import ProductCardSkeleton from '$lib/components/productCardSkeleton.svelte';

	export let data: import('./$types').PageData;
</script>

<div class="my-6 flex w-full items-center justify-center p-4">
	<div
		class="flex h-[256px] w-full items-center justify-center rounded-md bg-gray-300 md:h-96 md:w-[1200px]"
	/>
</div>
<section class="flex flex-col items-center gap-y-3">
	{#await data.streamed.products}
		<div class="space-y-1 p-1 md:p-4">
			<div title="Gå til kategori" class="flex items-center gap-2">
				<h2 class="text-xl font-medium">Laster inn produkter</h2>
				<LoaderCircle color="#7A38D0" class="animate-spin" />
			</div>
			<div class="grid grid-cols-2 gap-x-4 gap-y-4 md:grid-cols-3 md:gap-x-8 lg:grid-cols-4">
				<ProductCardSkeleton />
				<ProductCardSkeleton />
				<ProductCardSkeleton />
				<ProductCardSkeleton />
			</div>
		</div>
	{:then loadedProducts}
		{#each loadedProducts as { category, products }}
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
	{/await}
</section>
