import { writable } from 'svelte/store';
import type {products} from "@prisma/client";

export const handlelapp = writable<products[]>([]);

handlelapp.subscribe((value) => {
    console.log('oppdatert');
    globalThis.localStorage?.setItem('handlelapp', JSON.stringify(value));
});