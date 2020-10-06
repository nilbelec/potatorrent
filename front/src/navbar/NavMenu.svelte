<script>
  import { scale } from "svelte/transition";
  import { isNewVersion } from "../api.js";
  import { onMount } from "svelte";
  import Icon from "svelte-awesome";
  import { bell, github } from "svelte-awesome/icons";
  import NavLink from "./NavLink.svelte";

  let newVersion;
  const load = async () => {
    try {
      newVersion = await isNewVersion();
    } catch {}
  };
  onMount(load);
</script>

<style>
  ul {
    list-style: none;
    line-height: 1;
    display: flex;
    align-items: center;
  }

  li {
    display: inline;
    margin: 0 1rem;
  }
  .new-version {
    font-style: italic;
    font-size: 0.9rem;
    font-weight: bold;
    color: green;
  }
</style>

<ul>
  <li>
    <NavLink to={'/'}>Inicio</NavLink>
  </li>
  <li>
    <NavLink to={'/searches'}>Búsquedas Programadas</NavLink>
  </li>
  <li>
    <NavLink to={'/configuration'}>Configuración</NavLink>
  </li>
  {#if newVersion}
    <li class="new-version" in:scale|local={{ duration: 2000 }}>
      <NavLink to={'/update'}>¡Hay una nueva versión!</NavLink>
    </li>
  {/if}
  <li>
    <a href="https://github.com/nilbelec/potatorrent/" target="_blank">
      <Icon data={github} scale="2" />
    </a>
  </li>
</ul>
