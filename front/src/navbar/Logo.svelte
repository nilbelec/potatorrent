<script>
  import { Link } from "svelte-routing";
  import { fly } from "svelte/transition";
  import { getVersion } from "../api.js";
  import { onMount } from "svelte";

  let version;
  const loadVersion = async () => {
    version = await getVersion();
  };

  onMount(loadVersion);
</script>

<div>
  <img src="/public/img/logo.png" alt="Potatorrent logo" />
  potatorrent
  {#if version}
    <span in:fly={{ y: -100 }} class="version">{version.version}</span>
  {/if}
</div>

<style>
  div {
    display: flex;
    align-items: center;
    text-transform: uppercase;
    font-size: 1.5rem;
    font-family: "Anton", sans-serif;
    pointer-events: none;
  }

  img {
    height: 3rem;
    margin-right: 0.5rem;
  }
  .version {
    transform: rotate(-6deg);
    font-size: 0.8rem;
    margin-bottom: -28px;
    margin-left: -16px;
    color: #fff;
    padding: 2px 6px;
    background-color: #dc3545;
  }

  @media only screen and (max-width: 700px) {
    div {
      font-size: 1.3rem;
    }

    img {
      height: 2.4rem;
      margin-right: 0.1rem;
    }
    .version {
      font-size: 0.6rem;
      margin-bottom: -28px;
      margin-left: -16px;
      color: #fff;
      padding: 2px 6px;
      background-color: #dc3545;
    }
  }
</style>
