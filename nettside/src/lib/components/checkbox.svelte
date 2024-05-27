<script>
    export let kategori;
    export let tittel;
    let checked = false;
</script>
<div class=" flex justify-center items-center h-12 w-24 rounded-md {checked ? "bg-mainPurple text-white" : "bg-gray-200"}">
    <input on:change={() => checked = !checked} class="opacity-0 absolute" type="checkbox" id={tittel}>
    <label class=" h-full w-full flex justify-center items-center" for={tittel}>{tittel}</label>
</div>