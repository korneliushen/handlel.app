<script>
    let antallProdukter = 0
</script>

<main class=" relative w-screen max-w-[1200px] rounded-lg h-full flex flex-col px-5 overflow-y-scroll">
    <div class=" relative h-80 flex justify-center py-1 border-b-2 border-borderColor/50">
        <img class=" h-full" src="/favicon.png" alt="Produktbilde">
        <img class=" absolute bottom-1 left-1 w-10" src="/favicon.png" alt="nettside">
    </div>
    <div>
        <p>Sjokolade melk Q 1L</p>
        <div class=" flex justify-between">
            <p>Q-meieriene</p>
            <div>
                <p>26,70kr</p>
                <p>26,70/L</p>
            </div>
        </div>
    </div>
    <div class=" flex flex-col items-center">
        <p>Andre butikker</p>
        <div class=" w-full flex justify-between border border-borderColor rounded-md p-2 my-1">
            <div class=" flex w-12 items-center">
                <img class=" w-full" src="/favicon.png" alt="Butikklogo">
                <p>Meny</p>
            </div>
            <div>
                <p>27,90kr</p>
                <p>27,90/L</p>
            </div>
        </div>
        <div class=" w-full flex justify-between border border-borderColor rounded-md p-2 my-1">
            <div class=" flex w-12 items-center">
                <img class=" w-full" src="/favicon.png" alt="Butikklogo">
                <p>Meny</p>
            </div>
            <div>
                <p>27,90kr</p>
                <p>27,90/L</p>
            </div>
        </div>
    </div>
    <div>
        <div>
            <p>Om produktet</p>
        </div>
        <div>

        </div>
    </div>
    <div>
        <div>
            <p>Næringsinnhold per 100g</p>
            <img src="" alt="">
        </div>
        <div>

        </div>
    </div>
    <footer class=" fixed bottom-0 left-0 h-20 w-full border-t border-borderColor flex justify-center items-center">
        <div class=" w-4/5 h-12 rounded-lg">
            {#if antallProdukter !== 0}
                <div class=" flex justify-between items-center h-full w-full border border-purple-500 rounded-lg">
                    <button on:click={() => antallProdukter--} class=" w-14 font-extrabold text-xl">-</button>
                    <p>{antallProdukter}</p>
                    <button on:click={() => antallProdukter++} class=" w-14 font-extrabold text-xl">+</button>
                </div>
            {:else}
                <button on:click={() => antallProdukter++} class=" w-full h-full rounded-lg flex items-center justify-center bg-purple-500 border border-purple-500 active:scale-95 transition">
                    <p class=" text-white font-bold">Legg til i handlelisten</p>
                    <img src="" alt="">
                </button>
            {/if}
            
        </div>
    </footer>
</main>