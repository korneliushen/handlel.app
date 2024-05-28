<script lang="ts">
	import Productcard from "$lib/components/productcard.svelte";
  export let data: import("./$types").PageData
  console.log(data.products)
</script>

<section class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-[16px]">
  {#each data.products as product}
    <Productcard product={product}/>
  {/each}
</section>
