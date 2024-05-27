<script>
    import {BookPlus} from 'lucide-svelte'
    export let produkt;
    export let produsent;
    export let pris;
    export let prisliter;
    export let nettside;
    export let produktbilde;

    
</script>

<a href="#" title="product" class=" w-60 border border-borderColor rounded-lg h-72 flex flex-col justify-between hover:border-purple-500 shadow-md hover:scale-[1.01] transition">
    <div class=" relative h-40 flex justify-center py-1">
        <img class=" h-full" src={produktbilde} alt="Produktbilde">
        <img class=" absolute bottom-0 left-0 w-10" src={nettside} alt="nettside">
    </div>
    <div class=" pb-1">
        <div class=" space-y-0 mt-1 px-2">
            <p class=" font-extrabold">{produkt}</p>
            <p class=" text-gray-500/60 text-sm">{produsent}</p>
        </div>
        <div class=" flex mt-1 px-2 justify-between">
            <div class=" space-y-0">
                <p class=" font-extrabold">{pris}</p>
                <p class=" text-gray-500/60 text-sm border-t border-gray-500/60">{prisliter}</p>
            </div>
            <button on:click={() => alert("added to favorites")} class=" h-10 w-10 bg-purple-500 rounded-md active:scale-[0.93] transition flex justify-center items-center">
                <BookPlus color="#fff"/>
            </button>
        </div>
    </div>
</a>