<script>
  import Select from "svelte-select";
  import {
    searchOptions,
    searchSubCategories,
    searchTorrents
  } from "../api.js";

  export let searchParams = {
    category: undefined,
    subCategory: undefined,
    quality: undefined,
    text: ""
  };
  export let handleSubmit;

  let categories = [];
  let subCategories = [];
  let qualities = [];

  searchOptions().then(opts => {
    categories = opts.categories;
    qualities = opts.qualities;
  });

  function onSelectCategory(e) {
    searchParams.subCategory = undefined;
    const selectedVal = e.detail.value;
    searchSubCategories(selectedVal).then(subs => {
      subCategories = subs;
    });
  }
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
    color: #3f4f5f;
    height: 42px;
    line-height: 42px;
    width: 100%;
    background: transparent;
    font-size: 14px;
    letter-spacing: -0.08px;
    border: 1px solid #d8dbdf;
    border-radius: 3px;
    box-sizing: border-box;
    padding: 0 16px;
    outline: none;
  }
  input::placeholder {
    color: #78848f;
  }
  input:focus {
    border-color: #006fe8;
    outline: none;
  }
  input:hover {
    border-color: #b2b8bf;
  }
  button {
    grid-column-start: 1;
    grid-column-end: 4;
    margin: 0;
    height: 42px;
    background-color: white;
    border: 1px solid #d8dbdf;
  }
</style>

<form on:submit|preventDefault={handleSubmit}>
  <Select
    items={categories}
    bind:selectedValue={searchParams.category}
    on:select={onSelectCategory}
    on:clear={() => (searchParams.subCategory = undefined)}
    isDisabled={categories.length == 0}
    placeholder="Filtrar por categoría" />
  <Select
    items={subCategories}
    bind:selectedValue={searchParams.subCategory}
    isDisabled={searchParams.category == undefined || subCategories.length == 0}
    placeholder="Filtrar por subcategoría" />
  <Select
    items={qualities}
    bind:selectedValue={searchParams.quality}
    isDisabled={qualities.length == 0}
    placeholder="Filtrar por calidad o tipo" />
  <input
    bind:value={searchParams.text}
    type="search"
    placeholder="Filtrar por palabras" />
  <button type="submit">Buscar</button>
</form>
