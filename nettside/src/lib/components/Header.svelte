<script>
	import {Search, User, Home} from 'lucide-svelte'
	import { onMount } from 'svelte';

	let screensize = 0

	onMount(() => {
		screensize = window.innerWidth
	})
</script>

{#if (screensize > 640)}
	<header class="hidden items-center justify-between gap-4 p-3 mb-2 border-b w-full sm:flex">
		<div class=" flex-1">
			<img class=" h-12 ml-4" src="/handlelapp.png" alt="">
		</div>
		<search class=" flex-1 w-80 md:w-96">
			<form class=" flex items-center justify-between w-full rounded-md border border-borderColor/26 px-2 text-borderColor shadow-sm" action="/sok">
				<input
					type="text"
					placeholder="Søk etter produkter..."
					class="py-2 bg-transparent flex-1 outline-none text-black"
					name="search"
					autocomplete="off"
				/>
				<Search />
			</form>
		</search>
		<div class=" flex-1 flex justify-end items-center font-medium gap-4 ">
			<a href="/">Hjem</a>
			<a href="/">Produkter</a>
			<a href="/handlelapp">Handlelapp</a>
			<button class=" w-10 h-10 bg-mainPurple rounded-md flex justify-center items-center text-white active:scale-95 transition"><User /></button>
		</div>
	</header>
{:else}
	<header class=" flex sm:hidden fixed bottom-0 z-50 bg-white w-screen justify-center h-16 items-center border-t border-borderColor">
		<div class=" flex w-full justify-evenly items-center text-borderColor">
			<a class=" flex flex-col items-center" href="/">
				<Home size="25px" />
				<p class=" text-xs mt-1 text-black">Hjem</p>
			</a>
			<a class=" flex flex-col items-center" href="/sok">
				<Search size="25px" />
				<p class=" text-xs mt-1 text-black">Søk</p>
			</a>
			<a class=" flex flex-col items-center" href="/handlelapp">
				<img class=" h-[25px]" src="/handlelapp.png" alt="">
				<p class=" text-xs mt-1 text-black">Handlelapp</p>
			</a>
			<a class=" flex flex-col items-center" href="/">
				<User size="25px" />
				<p class=" text-xs mt-1 text-black">Konto</p>
			</a>
		</div>
	</header>
{/if}
