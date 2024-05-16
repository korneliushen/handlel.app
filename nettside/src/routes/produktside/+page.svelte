<script>
    let antallProdukter = 0
</script>

<main class=" relative w-screen max-w-[1200px] rounded-lg flex min-h-fit flex-col px-5">
    <div class=" relative min-h-80 flex justify-center py-1 border-b-2 border-borderColor/50">
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
    <div class=" mt-3 h-40">
        <div>
            <p class=" font-bold">Om produktet</p>
        </div>
        <div>
            <p>
                Lorem ipsum dolor sit amet consectetur, adipisicing elit. 
                Quisquam, esse! Quis quod corporis molestias natus voluptatum adipisci! 
                Tempora, provident laudantium? 
                Ex minima rerum omnis nam a, quos dignissimos itaque! 
                Omnis.
                Lorem ipsum dolor sit amet consectetur, adipisicing elit. Amet id deleniti sed hic dicta magni rem, doloremque molestias? Unde voluptate officia repudiandae enim! Voluptatibus consequatur iure sapiente eius molestias totam.
                Lorem ipsum dolor sit amet consectetur adipisicing elit. Quos voluptas exercitationem tenetur ipsum asperiores cupiditate minus at, adipisci pariatur impedit molestiae nulla veritatis suscipit ad cum consectetur vero odit perferendis!
                Lorem ipsum dolor sit amet consectetur adipisicing elit. Earum, eos excepturi assumenda maiores nam accusantium voluptates! Quas consectetur, consequatur maxime ipsam qui cumque amet quod quis sit blanditiis accusamus similique.
            </p>
        </div>
    </div>
    <div class=" mb-24 h-20">
        <div>
            <p>Næringsinnhold per 100g</p>
            <img src="" alt="">
        </div>
        <div>

        </div>
    </div>
    <footer class=" bg-white fixed bottom-0 left-0 h-20 w-full border-t border-borderColor flex justify-center items-center">
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