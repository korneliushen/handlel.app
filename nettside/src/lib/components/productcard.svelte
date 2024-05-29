<script lang="ts">
    import {BookPlus} from 'lucide-svelte'
    import type { products } from '@prisma/client';
    export let product: products;

</script>

<div class=" z-20 w-44 md:w-60 h-full border border-borderColor/26 rounded-lg flex flex-col justify-between hover:border-purple-500 shadow-sm hover:shadow-md hover:scale-[1.01] transition p-2">
    <a href="/produkt/{product.id}" title="product" class=" relative max-h-40 flex justify-center py-1 aspect-square">
        <img loading="lazy" src={product.imagelink} alt="produktbilde">
        <img loading="lazy" class=" absolute bottom-0 left-0 w-10 rounded-md ml-2" src="/{product.prices[0].Store}.svg" alt="nettside">
    </a>
    <div class=" pb-1">
        <a href="/produkt/{product.id}" title="product">
            <div class=" space-y-0 mt-1 px-2">
                <p class=" font-semibold">{product.title}</p>
                <p class=" text-gray-500/60 text-sm">{product.vendor}</p>
            </div>
        </a>
        <div class=" flex mt-1 px-2 justify-between">
            <a class=" w-full" href="/produkt/{product.id}" title="product">
                <div class=" space-y-0 w-fit">
                    <p class=" font-semibold">{product.prices[0].Price + " kr"}</p>
                    <p class=" text-gray-500/60 text-sm border-t border-gray-500/60">{product.prices[0].UnitPrice} kr/{product.unittype || "N/A"}</p>
                </div>
            </a>
            <button on:click={() => alert("added to handlelapp")} class=" h-10 w-10 aspect-square bg-purple-500 rounded-md flex justify-center items-center hover:brightness-110 transition z-50">
                <BookPlus color="#fff"/>
            </button>
        </div>
    </div>
</div>