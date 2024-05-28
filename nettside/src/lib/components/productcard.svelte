<script lang="ts">
    import {BookPlus} from 'lucide-svelte'
    import type { products } from '@prisma/client';
    export let product: products;

</script>

<a href="/produkt/{product.id}" title="product" class=" w-60 border border-borderColor rounded-lg h-72 flex flex-col justify-between hover:border-purple-500 shadow-md hover:scale-[1.01] transition">
    <div class=" relative h-40 flex justify-center py-1">
        <img loading="lazy" class=" h-full" src={product.imagelink} alt="produktbilde">
        <img loading="lazy" class=" absolute bottom-0 left-0 w-10 rounded-md ml-2" src="/{product.prices[0].Store}.svg" alt="nettside">
    </div>
    <div class=" pb-1">
        <div class=" space-y-0 mt-1 px-2">
            <p class=" font-extrabold">{product.title}</p>
            <p class=" text-gray-500/60 text-sm">{product.vendor}</p>
        </div>
        <div class=" flex mt-1 px-2 justify-between">
            <div class=" space-y-0">
                <p class=" font-extrabold">{product.prices[0].Price + " kr"}</p>
                <p class=" text-gray-500/60 text-sm border-t border-gray-500/60">{product.prices[0].UnitPrice + "/" + (product.unittype ? product.unittype : "N/A")}</p>
            </div>
            <button on:click={() => alert("added to favorites")} class=" h-10 w-10 bg-purple-500 rounded-md active:scale-[0.93] transition flex justify-center items-center">
                <BookPlus color="#fff"/>
            </button>
        </div>
    </div>
</a>