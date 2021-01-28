import { writable } from 'svelte/store';

export const searchParamsStore = writable({
    category: undefined,
    categoryLabel: undefined,
    subCategory: undefined,
    subCategoryLabel: undefined,
    quality: undefined,
    qualityLabel: undefined,
    text: ""
});

export const searchOptsStore = writable(undefined);