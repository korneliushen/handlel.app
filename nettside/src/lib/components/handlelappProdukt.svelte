<script>
        import { Minus, Plus, X } from "lucide-svelte";
        import { onMount } from 'svelte';

        let desktopView = false
        let antallProdukt = 1

        onMount(() => {
            if (window.innerWidth > 1024) {
                desktopView = true
            }
        })
</script>

<div class=" mt-6">
    <div class=" h-20 flex items-center justify-between border border-gray-400 rounded-xl mt-5">
        <div class=" flex h-full items-center">
            <img class=" w-16" src="/favicon.png" alt="">
            <div>
                <p class=" font-medium text-sm">Sjokolade melk Q 1L</p>
                <p class=" text-xs text-gray-400">Q-meieriene</p>
                {#if !desktopView}
                    <p class=" font-medium text-sm">29.90 kr</p>
                {/if}
            </div>
        </div>
        <div class=" h-full flex items-center mr-4">
            {#if desktopView}
                <div class=" mr-3 text-end">
                    <p class=" font-medium">29.90 kr</p>
                    <p class=" font-medium text-xs text-gray-400">29.90 kr/L</p>
                </div>
            {/if}
            <div class=" border border-mainPurple rounded-md {desktopView ? "h-10 w-32" : " h-9 w-28"} flex justify-between items-center">
                <button on:click={() => antallProdukt--} class=" w-10 flex justify-center items-center"><Minus size="15px"/></button>
                <p>{antallProdukt}</p>
                <button on:click={() => antallProdukt++} class=" w-10 flex justify-center items-center"><Plus size="15px"/></button>
            </div>
            {#if desktopView}
                <button class=" ml-3 text-gray-400"><X size="15px"/></button>
            {/if}
        </div>
    </div>
</div>