<script>
    import {ChevronDown} from 'lucide-svelte'
    import autoAnimate from '@formkit/auto-animate';
    let antallProdukter = 0
    let dropdown = false;
    let lesmer = '10rem';
    let lesmerBool = false;

    function lesmerFunc() {
        lesmerBool = true
        if (lesmerBool) {
            lesmer = 'fit-content'
        } else {
            lesmer = '10rem'
        };
    };
</script>

<main class=" relative w-screen max-w-[400px] rounded-lg flex flex-col px-5 overflow-auto">
    <div class=" relative min-h-80 flex justify-center items-center py-1 border-b-2 border-borderColor/50">
        <img class=" h-full" src="/favicon.png" alt="Produktbilde">
        <img class=" absolute bottom-1 left-1 w-10" src="/favicon.png" alt="nettside">
    </div>
    <div class=" mt-5">
        <p class=" font-bold text-2xl">Sjokolade melk Q 1L</p>
        <div class=" flex justify-between mt-2">
            <p class=" text-lg text-gray-500/60">Q-meieriene</p>
            <div class=" text-end">
                <p class=" font-bold text-2xl text-purple-500">26,70kr</p>
                <p class=" text-lg text-gray-500/60">26,70kr/L</p>
            </div>
        </div>
    </div>
    <div class=" flex flex-col items-center">
        <p class=" font-bold text-xl">Andre butikker</p>
        <div class=" w-full flex justify-between border border-borderColor rounded-md p-2 my-1">
            <div class=" flex w-12 items-center">
                <img class=" w-full mx-1" src="/favicon.png" alt="Butikklogo">
                <p class=" font-bold">Meny</p>
            </div>
            <div class=" text-end">
                <p class=" font-bold text-lg">27,90kr</p>
                <p class=" text-gray-500/60 text-sm">27,90kr/L</p>
            </div>
        </div>
        <div class=" w-full flex justify-between border border-borderColor rounded-md p-2 my-1">
            <div class=" flex w-12 items-center">
                <img class=" w-full mx-1" src="/favicon.png" alt="Butikklogo">
                <p class=" font-bold">Meny</p>
            </div>
            <div class=" text-end">
                <p class=" font-bold text-lg">27,90kr</p>
                <p class=" text-gray-500/60 text-sm">27,90kr/L</p>
            </div>
        </div>
    </div>
    <div use:autoAnimate class=" my-3 relative overflow-hidden" style="max-height: {lesmer};">
        {#if !lesmerBool}
            <div class=" w-full h-full bg-gradient-to-b from-transparent to-85% to-white absolute z-10"/>
        {/if}
        <div>
            <p class=" font-bold text-xl py-1">Om produktet</p>
        </div>
        <div>
            <p>
                Lorem ipsum dolor sit amet consectetur, adipisicing elit. 
                Quisquam, esse! Quis quod corporis molestias natus voluptatum adipisci! 
                Tempora, provident laudantium? 
                Ex minima rerum omnis nam a, quos dignissimos itaque! 
                Omnis.
            </p>
            <div class=" flex justify-between my-4">
                <p class=" font-bold text-lg">Mengde</p>
                <p>1 Liter</p>
            </div>
            <div class=" flex my-4 w-full justify-between">
                <p class=" font-bold text-lg">Ingredienser</p>
                <p class=" w-2/3 text-end">lettmelk 4%, sukker, stivelse, aroma</p>
            </div>
        </div>
        {#if !lesmerBool}
            <button on:click={() => lesmerFunc()} class=" text-purple-500 font-bold bottom-1 left-1 absolute z-20">+ Les mer...</button>
        {/if}
    </div>
    <div use:autoAnimate class=" mb-24 mt-3 border-y border-borderColor">
        <button on:click={() => dropdown = !dropdown} class=" flex items-center justify-between w-full py-3 h-20">
            <p class=" font-bold text-xl py-1">Næringsinnhold per 100g</p>
            <ChevronDown class=" {dropdown ? "rotate-180" : "rotate-0"}" strokeWidth={3}/>
        </button>
        {#if dropdown}
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
            <div class=" flex w-full justify-between border-t border-borderColor py-2">
                <p class=" font-bold">Kalorier</p>
                <p class=" w-2/3 text-end">40kcal</p>
            </div>
        {/if}
    </div>
    <footer class=" bg-white fixed bottom-0 left-0 h-20 w-full border-t border-borderColor flex justify-center items-center z-50">
        <div class=" w-4/5 h-12 rounded-lg">
            {#if antallProdukter !== 0}
                <div class=" flex justify-between items-center h-full w-full border border-purple-500 rounded-lg">
                    <button on:click={() => antallProdukter--} class=" w-14 font-extrabold text-xl">-</button>
                    <p>{antallProdukter}</p>
                    <button on:click={() => antallProdukter++} class=" w-14 font-extrabold text-xl">+</button>
                </div>
            {:else}
                <button on:click={() => antallProdukter++} class=" w-full h-full rounded-lg flex items-center justify-center bg-purple-500 border border-purple-500">
                    <p class=" text-white font-bold">Legg til i handlelisten</p>
                    <img src="" alt="">
                </button>
            {/if}
            
        </div>
    </footer>
</main>