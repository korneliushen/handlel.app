 <script>
	import Checkbox from "./checkbox.svelte";
    import { Minus, Plus, X } from "lucide-svelte";
    import { filter } from "$lib/stores/filter";
    import { onMount } from 'svelte';

    let antallButikker = 1
    let desktopView = false

    onMount(() => {
        if (window.innerWidth > 1024) {
            desktopView = true
        }
    })

 </script>

 <div class=" w-full max-w-[400px] px-4 {desktopView ? "" : " absolute top-0"} bg-white {$filter ? "" : "hidden"}">
     <div class=" w-full mb-10">
        <div class=" flex w-full justify-between">
            <p class=" font-semibold text-2xl mb-3">Kjeder</p>
            {#if !desktopView}
                <button on:click={() => $filter = false}><X /></button>
            {/if}
        </div>
        <div class=" w-full grid grid-cols-3 mb-4">
            <Checkbox tittel="Meny" kategori="kjede"/>
            <Checkbox tittel="Joker" kategori="kjede"/>
            <Checkbox tittel="Spar" kategori="kjede"/>
        </div>
        <div class=" w-full grid grid-cols-3 mb-4">
            <Checkbox tittel="Coop Prix" kategori="kjede"/>
            <Checkbox tittel="Kiwi" kategori="kjede"/>
            <Checkbox tittel="Rema 1000" kategori="kjede"/>
        </div>
        <div class=" w-full grid grid-cols-3">
            <Checkbox tittel="Coop Mega" kategori="kjede"/>
            <Checkbox tittel="Extra" kategori="kjede"/>
            <Checkbox tittel="Oda" kategori="kjede"/>
        </div>
     </div>
     <div class=" w-full pt-6 mb-10 border-t border-gray-300">
        <p class=" font-semibold text-2xl mb-3">Antall butikker</p>
        <div class=" border border-mainPurple rounded-md h-10 flex w-40 justify-between items-center">
            <button on:click={() => antallButikker--} class=" w-10 flex justify-center items-center"><Minus size="15px"/></button>
            <p>{antallButikker}</p>
            <button on:click={() => antallButikker++} class=" w-10 flex justify-center items-center"><Plus size="15px"/></button>
        </div>
     </div>
     <div class=" w-full pt-6 mb-10 border-t border-gray-300">
        <p class=" font-semibold text-2xl mb-3">Lokasjoner</p>
        <div class=" w-full grid grid-cols-3">
            <Checkbox tittel="Online" kategori="lokasjon"/>
            <Checkbox tittel="Butikk" kategori="lokasjon"/>
        </div>
     </div>
     <div class=" w-full pt-6 mb-10 border-t border-gray-300">
        <button class=" border border-red-500 text-red-500 rounded-md px-4 py-2 active:scale-95 transition">Tøm handlelapp</button>
     </div>
 </div>