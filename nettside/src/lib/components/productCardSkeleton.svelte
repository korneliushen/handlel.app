<script>
  import { BookPlus } from "lucide-svelte";
</script>

<div class="border border-gray-300 relative f<lex h-full w-44 flex-col justify-between rounded-lg border p-2 shadow-sm transition hover:scale-[1.01] hover:border-purple-500 hover:shadow-md md:w-60">
    <div title="product" class="relative flex aspect-square max-h-40 items-center justify-center py-1 lg:w-full">
        <div class="h-36 w-52 rounded-md bg-gray-200 animate-pulse"></div>
        <div class="absolute bottom-0 left-0 ml-2 w-10 rounded-md bg-gray-200 animate-pulse"></div>
    </div>
    <div class="pb-1">
        <div class="mt-1 space-y-0 px-2">
            <div class="font-semibold bg-gray-200 h-4 w-24 animate-pulse"></div>
            <div class="text-sm text-gray-500/60 bg-gray-200 h-3 w-20 animate-pulse"></div>
        </div>
        <div class="mt-1 flex justify-between px-2">
            <div class="w-fit space-y-0">
                <div class="font-semibold bg-gray-200 h-4 w-12 animate-pulse"></div>
                <div class=" text-sm text-gray-500/60 bg-gray-200 h-3 w-24 animate-pulse"></div>
            </div>
            <button class="z-10 flex aspect-square h-10 w-10 items-center justify-center rounded-md bg-purple-500 transition hover:brightness-110">
				<BookPlus color="#fff" />
            </button>
        </div>
    </div>
</div>


