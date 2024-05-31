import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { ExtendedProduct } from '$lib/types/extendedPrisma';

export const handlelapp = writable<ExtendedProduct[]>(
	browser ? JSON.parse(window.localStorage.getItem('handlelapp') ?? '[]') : []
);

if (browser) {
	handlelapp.subscribe((value) => {
		window.localStorage.setItem('handlelapp', JSON.stringify(value));
	});
}
