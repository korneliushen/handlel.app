<script lang="ts">
  import { ChevronRight } from 'lucide-svelte';
	import Productcard from '$lib/components/productcard.svelte';
	export let data: import('./$types').PageData;
	console.log(data.products);
</script>

<section class="flex flex-col gap-y-14">
  <img class="w-full flex items-center justify-center h-96 rounded-md my-6" src="http://1.vgc.no/vgnett-prod/img/vgLogoSquare.png?28042014-1" alt="Banner"/>
	{#each Object.entries(data.products) as [category, products]}
		<div class="space-y-1">
    <a href="#" title="Gå til kategori" class="flex items-center">
			<h2 class="text-xl font-medium">{category}</h2>
      <ChevronRight color="#7A38D0" />
    </a>
			<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-[16px]">
				{#each products as product}
          <Productcard product={product} />
				{/each}
			</div>
		</div>
	{/each}
</section>
