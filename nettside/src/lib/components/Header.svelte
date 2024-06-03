<script>
	import { Search, User, Home } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { handlelapp } from '$lib/stores/handlelapp';

	let screensize = 0;

	onMount(() => {
		screensize = window.innerWidth;
	});
</script>

{#if screensize > 640}
	<header class="mb-2 hidden w-full items-center justify-between gap-4 border-b p-3 sm:flex">
		<div class="flex-1">
			<div class="w-max">
				<a href="/">
					<img class="ml-4 h-12" src="/handlelapp.png" alt="logo" />
				</a>
			</div>
		</div>
		<search class="w-80 flex-1 md:w-96">
			<form
				class="border-borderColor/26 flex w-full items-center justify-between rounded-md border px-2 text-borderColor shadow-sm"
				action="/sok"
			>
				<input
					type="text"
					placeholder="Søk etter produkter..."
					class="flex-1 bg-transparent py-2 text-black outline-none"
					name="search"
					autocomplete="off"
				/>
				<Search />
			</form>
		</search>
		<div class=" flex flex-1 items-center justify-end gap-4 font-medium">
			<a href="/">Hjem</a>
			<a href="/">Produkter</a>
			<a class=" relative" href="/handlelapp">
				Handlelapp
				<p class=" absolute top-[-17px] right-[-10px] border border-red-600 rounded-full text-sm w-6 h-6 place-items-center grid text-red-600 font-bold">{$handlelapp.length}</p>
			</a>
			<button
				class=" flex h-10 w-10 items-center justify-center rounded-md bg-mainPurple text-white transition active:scale-95"
				><User /></button
			>
		</div>
	</header>
{:else}
	<header
		class=" fixed bottom-0 z-50 flex h-16 w-screen items-center justify-center border-t border-borderColor bg-white sm:hidden"
	>
		<div class=" flex w-full items-center justify-evenly text-borderColor">
			<a class=" flex flex-col items-center" href="/">
				<Home size="25px" />
				<p class=" mt-1 text-xs text-black">Hjem</p>
			</a>
			<a class=" flex flex-col items-center" href="/sok">
				<Search size="25px" />
				<p class=" mt-1 text-xs text-black">Søk</p>
			</a>
			<a class=" flex flex-col items-center relative" href="/handlelapp">
				<img class=" h-[25px]" src="/handlelapp.png" alt="" />
				<p class=" mt-1 text-xs text-black">Handlelapp</p>
				<p class=" absolute top-[-5px] right-[-5px] rounded-full text-sm w-5 h-5 place-items-center grid text-red-600 font-bold">{$handlelapp.length}</p>
			</a>
			<a class=" flex flex-col items-center" href="/">
				<User size="25px" />
				<p class=" mt-1 text-xs text-black">Konto</p>
			</a>
		</div>
	</header>
{/if}
