<script>
  import Select from "svelte-select";
  import {
    searchOptions,
    searchSubCategories,
    searchTorrents
  } from "../../api.js";
  import { onMount } from "svelte";

  export let searchParams = {
    category: undefined,
    subCategory: undefined,
    quality: undefined,
    text: ""
  };
  export let handleSubmit;
  export let disabled = false;

  let categories = [];
  let subCategories = [];
  let qualities = [];

  const load = async () => {
    try {
      const opts = await searchOptions();
      categories = opts.categories;
      qualities = opts.qualities;
    } catch {}
  };

  function onSelectCategory(e) {
    searchParams.subCategory = undefined;
    const selectedVal = e.detail.value;
    searchSubCategories(selectedVal).then(subs => {
      subCategories = subs;
    });
  }
  onMount(load);
</script>

<style>
  form {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    column-gap: 10px;
    row-gap: 10px;
    margin-bottom: 1rem;
  }
  input {
    grid-column-start: 1;
    grid-column-end: 4;
  }
  button {
    grid-column-start: 1;
    grid-column-end: 4;
    height: 42px;
  }
</style>

<form on:submit|preventDefault={handleSubmit}>
  <Select
    items={categories}
    bind:selectedValue={searchParams.category}
    on:select={onSelectCategory}
    on:clear={() => (searchParams.subCategory = undefined)}
    isDisabled={disabled || categories.length == 0}
    placeholder="Filtrar por categoría" />
  <Select
    items={subCategories}
    bind:selectedValue={searchParams.subCategory}
    isDisabled={disabled || searchParams.category == undefined || subCategories.length == 0}
    placeholder="Filtrar por subcategoría" />
  <Select
    items={qualities}
    bind:selectedValue={searchParams.quality}
    isDisabled={disabled || qualities.length == 0}
    placeholder="Filtrar por calidad o tipo" />
  <input
    {disabled}
    bind:value={searchParams.text}
    type="search"
    placeholder="Filtrar por palabras" />
  <button {disabled} type="submit">Buscar</button>
</form>
