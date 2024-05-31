<script>
	import Filter from "$lib/components/filter.svelte";
	import HandlelappProdukt from "$lib/components/handlelappProdukt.svelte";
    import { filter } from "$lib/stores/filter";
    import { onMount } from 'svelte';
    import { ChevronDown, SlidersHorizontal, MapPin } from "lucide-svelte";

    let desktopView = false
    let selected = false;

    onMount(() => {
        if (window.innerWidth > 1024) {
            $filter = true
            desktopView = true
        }
    })
</script>

<main class=" relative flex w-screen px-4 max-w-[400px] lg:max-w-[1200px] py-4 md:py-0">
    <Filter/>
    <div class=" w-full {desktopView ? "border-l border-gray-200 pl-8" : ""}">
        <div class=" w-full flex h-16">
            <button on:click={() => selected = false} class="flex flex-col flex-1 px-1 h-full cursor-pointer {!selected && "font-semibold"}">
                <div class=" flex flex-col justify-between h-full w-full px-2">
                    <div class=" flex justify-between">
                        <p class=" text-black {!selected && "!text-mainPurple !font-bold"}">Billigste</p>
                        <img class=" w-7 " src="/favicon.png" alt="">
                    </div>
                    <div class=" flex justify-between text-sm pb-1">
                        <p>10 min</p>
                        <p>27.60 kr</p>
                    </div>
                </div>
                <div class="h-1 w-full flex items-end">
                    <div class=" w-full h-[0.16rem] rounded-t-2xl bg-gray-400 {!selected && "!bg-mainPurple !h-1"}"/>
                </div>
            </button>
            <button on:click={() => selected = true} class=" flex flex-col flex-1 px-1 h-full cursor-pointer {selected && "font-semibold"}">
                <div class=" flex flex-col justify-between h-full w-full px-2">
                    <div class=" flex justify-between">
                        <p class=" text-black {selected && "!text-mainPurple !font-bold"}">Raskeste</p>
                        <img class=" w-7 " src="/favicon.png" alt="">
                    </div>
                    <div class=" flex justify-between text-sm pb-1">
                        <p>8 min</p>
                        <p>29.90 kr</p>
                    </div>
                </div>
                <div class="h-1 w-full flex items-end">
                    <div class=" w-full h-[0.16rem] rounded-t-2xl bg-gray-400 {selected && "!bg-mainPurple !h-1"}"/>
                </div>
            </button>
        </div>
        <div class=" mt-3 flex justify-between w-full px-2">
            <div>
                <p class=" text-2xl">Handlelapp</p>
                <p class=" text-sm flex mt-2">Sorter etter:<button class=" text-mainPurple font-bold flex mx-1 items-center">Pris <ChevronDown size="20px"/></button></p>
            </div>
            {#if !desktopView}
                <div class=" flex mt-2">
                    <button on:click={() => $filter = true} class=" w-9 h-9 bg-mainPurple rounded-md flex justify-center items-center mx-[0.13rem]"><SlidersHorizontal color="#ffffff" /></button>
                    <button class=" w-9 h-9 bg-mainPurple rounded-md flex justify-center items-center mx-[0.13rem]"><MapPin color="#ffffff"/></button>
                </div>
            {/if}
        </div>
        <HandlelappProdukt/>
        <HandlelappProdukt/>
        <HandlelappProdukt/>
    </div>
</main>
