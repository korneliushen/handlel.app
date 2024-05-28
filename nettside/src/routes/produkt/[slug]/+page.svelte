<script lang="ts">
    import {ChevronDown, Trash2, Minus, Plus} from 'lucide-svelte'
    import autoAnimate from '@formkit/auto-animate';
	import { onMount } from 'svelte';
    let antallProdukter = 0
    let dropdown = false;
    let lesmer = '10rem';
    let lesmerBool = false;

    onMount(() => {
        if (window.innerWidth > 1024) {
            dropdown = true
        }
    })

    function lesmerFunc() {
        lesmerBool = true
        if (lesmerBool) {
            lesmer = 'fit-content'
        } else {
            lesmer = '10rem'
        };
    };
    export let data: import('./$types').PageData
</script>

<main class=" relative w-screen max-w-[600px] rounded-lg flex flex-col px-5 overflow-hidden lg:max-w-[1200px] lg:px-20">
    <div class=" block lg:flex w-full justify-between">
        <div class=" relative min-h-80 lg:w-1/2 flex justify-center items-center py-1 border-b border-borderColor">
            <img class=" w-10/12 " src={data.product.imagelink} alt="Produktbilde">
            <img class=" absolute bottom-1 left-1 w-10 rounded-md" src="/{data.product.prices[0].Store}.svg" alt="nettside">
        </div>
        <div class=" lg:w-[400px] flex flex-col justify-between">
            <div>
                <div class=" mt-5">
                    <p class=" font-bold text-2xl">{data.product?.title}</p>
                    <div class=" flex justify-between mt-2">
                        <p class=" text-lg text-gray-500/60">{data.product?.vendor}</p>
                        <div class=" text-end">
                            <p class=" font-bold text-2xl text-mainPurple">{data.product?.prices[0].Price} kr</p>
                            <p class=" text-lg text-gray-500/60">{data.product?.prices[0].UnitPrice}/{data.product?.unittype}</p>
                        </div>
                    </div>
                </div>
                <div class=" flex flex-col items-center">
                    <p class=" font-bold text-xl">Andre butikker</p>
                    {#each data.product?.prices.slice(1) as price}
                        <div class=" w-full flex justify-between border border-borderColor rounded-md p-2 my-1">
                            <div class=" flex w-12 items-center">
                                <img class=" w-full mr-3 rounded-md" src="/{price.Store}.svg" alt="Butikklogo">
                                <p class=" font-bold">{price.Store[0].toUpperCase() + price.Store.substring(1)}</p>
                            </div>
                            <div class=" text-end">
                                <p class=" font-bold text-lg">{price.Price} kr</p>
                                <p class=" text-gray-500/60 text-sm">{price.UnitPrice}/{data.product.unittype}</p>
                            </div>
                        </div>
                    {/each}
                    
                </div>
            </div>
            <div class=" bg-white fixed bottom-0 left-0 h-20 w-full border-t border-borderColor flex justify-center items-center z-50 lg:relative">
                <div class=" w-4/5 h-12 rounded-lg lg:w-full">
                    {#if antallProdukter !== 0}
                        <div class=" flex justify-between items-center h-full w-full border border-mainPurple rounded-lg">
                            {#if antallProdukter !== 1}
                                <button on:click={() => antallProdukter--} class=" w-14 font-extrabold text-xl flex justify-center"><Minus /></button>
                                {:else}
                                <button on:click={() => antallProdukter--} class=" w-14 font-extrabold text-xl flex justify-center"><Trash2 /></button>
                            {/if}
                            <p>{antallProdukter}</p>
                            <button on:click={() => antallProdukter++} class=" w-14 font-extrabold text-xl flex justify-center"><Plus /></button>
                        </div>
                    {:else}
                        <button on:click={() => antallProdukter++} class=" w-full h-full rounded-lg flex items-center justify-center bg-mainPurple border border-mainPurple">
                            <p class=" text-white font-bold">Legg til i handlelisten</p>
                            <img src="" alt="">
                        </button>
                    {/if}
                </div>
            </div>
        </div>
    </div>
    <div class=" lg:flex justify-between">
        <div use:autoAnimate class=" my-3 relative overflow-hidden lg:w-1/2" style="max-height: {lesmer};">
            {#if !lesmerBool}
                <div class=" w-full h-full bg-gradient-to-b from-transparent to-85% to-white absolute z-10"/>
            {/if}
            <div>
                <p class=" font-bold text-xl py-1">Om produktet</p>
            </div>
            <div>
                <p>
                    {data.product.description || "Dette produktet har ingen beskrivelse"}
                </p>
                <div class=" flex justify-between my-4">
                    <p class=" font-bold text-lg">Mengde</p>
                    <p>{data.product.weight || "N/A"}</p>
                </div>
                <div class=" flex my-4 w-full justify-between">
                    <p class=" font-bold text-lg">Ingredienser</p>
                    <p class=" w-2/3 text-end">{data.product.ingredients || "N/A"}</p>
                </div>
            </div>
            {#if !lesmerBool}
                <button on:click={() => lesmerFunc()} class=" text-mainPurple font-bold bottom-1 left-1 absolute z-20">+ Les mer...</button>
            {/if}
        </div>
        <div use:autoAnimate class=" mb-24 mt-3 border-y border-borderColor lg:w-[400px] lg:mt-0 lg:h-fit">
            <button on:click={() => dropdown = !dropdown} class=" flex items-center justify-between w-full py-3 h-20">
                <p class=" font-bold text-xl py-1">Næringsinnhold per 100g</p>
                <ChevronDown class=" {dropdown ? "rotate-180" : "rotate-0"}" strokeWidth={3}/>
            </button>
            {#if dropdown}
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Kalorier</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.Kalorier || "N/A"}</p>
                </div>
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Energi</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.Energi || "N/A"}</p>
                </div>
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Fett</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.Fett || "N/A"}</p>
                </div>
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Karbohydrater</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.Karbohydrater || "N/A"}</p>
                </div>
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Kostfiber</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.Kostfiber || "N/A"}</p>
                </div>
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Mettet fett</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.MettetFett || "N/A"}</p>
                </div>
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Protein</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.Protein || "N/A"}</p>
                </div>
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Salt</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.Salt || "N/A"}</p>
                </div>
                <div class=" flex w-full justify-between border-t border-borderColor p-2">
                    <p>Sukkerarter</p>
                    <p class=" w-2/3 text-end font-bold">{data.product.nutritionalcontent.Sukkerarter || "N/A"}</p>
                </div>
            {/if}
        </div>
    </div>
</main>