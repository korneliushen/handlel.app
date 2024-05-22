<script>
    let selected = false;
</script>

<main class=" w-screen px-4">
    <div class=" w-full flex h-16">
        <button on:click={() => selected = false} class=" flex flex-col flex-1 px-1 h-full cursor-pointer {!selected && "font-semibold"}">
            <div class=" flex flex-col justify-between h-full w-full px-2">
                <div class=" flex justify-between">
                    <p class=" text-black {!selected && "!text-purple-500 !font-bold"}">Billigste</p>
                    <img class=" w-7 " src="/favicon.png" alt="">
                </div>
                <div class=" flex justify-between text-sm pb-1">
                    <p>10 min</p>
                    <p>27.60 kr</p>
                </div>
            </div>
            <div class="h-1 w-full flex items-end">
                <div class=" w-full h-[0.16rem] rounded-t-2xl bg-gray-400 {!selected && "!bg-purple-500 !h-1"}"/>
            </div>
        </button>
        <button on:click={() => selected = true} class=" flex flex-col flex-1 px-1 h-full cursor-pointer {selected && "font-semibold"}">
            <div class=" flex flex-col justify-between h-full w-full px-2">
                <div class=" flex justify-between">
                    <p class=" text-black {selected && "!text-purple-500 !font-bold"}">Raskeste</p>
                    <img class=" w-7 " src="/favicon.png" alt="">
                </div>
                <div class=" flex justify-between text-sm pb-1">
                    <p>8 min</p>
                    <p>29.90 kr</p>
                </div>
            </div>
            <div class="h-1 w-full flex items-end">
                <div class=" w-full h-[0.16rem] rounded-t-2xl bg-gray-400 {selected && "!bg-purple-500 !h-1"}"/>
            </div>
        </button>
    </div>
    <div class=" mt-3 flex justify-between w-full px-5">
        <div>
            <p class=" text-2xl">Handlelapp</p>
            <p class=" text-sm flex">Sorter etter:<button class=" text-purple-500 font-bold flex mx-1">Pris <img src="" alt="pil ned"></button></p>
        </div>
        <div class=" flex mt-2">
            <button class=" w-9 h-9 bg-purple-500 rounded-md flex justify-center items-center mx-[0.13rem]"><img class=" w-5/6" src="/favicon.png" alt=""></button>
            <button class=" w-9 h-9 bg-purple-500 rounded-md flex justify-center items-center mx-[0.13rem]"><img class=" w-5/6" src="/favicon.png" alt=""></button>
        </div>
    </div>
    <div class=" mt-3 px-3">
        <div class=" relative h-20 flex items-center justify-between border border-gray-400 rounded-xl mt-5">
            <img class=" w-7 absolute right-[-10px] top-[-15px]" src="/favicon.png" alt="">
            <div class=" flex h-full items-center">
                <img class=" w-16" src="/favicon.png" alt="">
                <div>
                    <p class=" font-medium">Sjokolade melk Q 1L</p>
                    <p class=" text-xs text-gray-400">Q-meieriene</p>
                    <p class=" font-medium">29.90 kr</p>
                </div>
            </div>
            <div class=" h-full flex items-center mr-4">
                <div class=" border border-purple-500 rounded-md h-9 flex w-32 justify-between items-center">
                    <button class=" w-10 flex justify-center items-center">-</button>
                    <p>1</p>
                    <button class=" w-10 flex justify-center items-center">+</button>
                </div>
            </div>
        </div>
        <div class=" relative h-20 flex items-center justify-between border border-gray-400 rounded-xl mt-5">
            <img class=" w-7 absolute right-[-10px] top-[-15px]" src="/favicon.png" alt="">
            <div class=" flex h-full items-center">
                <img class=" w-16" src="/favicon.png" alt="">
                <div>
                    <p class=" font-medium">Sjokolade melk Q 1L</p>
                    <p class=" text-xs text-gray-400">Q-meieriene</p>
                    <p class=" font-medium">29.90 kr</p>
                </div>
            </div>
            <div class=" h-full flex items-center mr-4">
                <div class=" border border-purple-500 rounded-md h-9 flex w-32 justify-between items-center">
                    <button class=" w-10 flex justify-center items-center">-</button>
                    <p>1</p>
                    <button class=" w-10 flex justify-center items-center">+</button>
                </div>
            </div>
        </div>
        <div class=" relative h-20 flex items-center justify-between border border-gray-400 rounded-xl mt-5">
            <img class=" w-7 absolute right-[-10px] top-[-15px]" src="/favicon.png" alt="">
            <div class=" flex h-full items-center">
                <img class=" w-16" src="/favicon.png" alt="">
                <div>
                    <p class=" font-medium">Sjokolade melk Q 1L</p>
                    <p class=" text-xs text-gray-400">Q-meieriene</p>
                    <p class=" font-medium">29.90 kr</p>
                </div>
            </div>
            <div class=" h-full flex items-center mr-4">
                <div class=" border border-purple-500 rounded-md h-9 flex w-32 justify-between items-center">
                    <button class=" w-10 flex justify-center items-center">-</button>
                    <p>1</p>
                    <button class=" w-10 flex justify-center items-center">+</button>
                </div>
            </div>
        </div>
    </div>
</main>