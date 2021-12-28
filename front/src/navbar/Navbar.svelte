<script>
  import Logo from "./Logo.svelte";
  import NavMenu from "./NavMenu.svelte";
  import Icon from "svelte-awesome";
  import { bars, times } from "svelte-awesome/icons";
  import { fly } from "svelte/transition";
  import NavLink from "./NavLink.svelte";

  let show = false;
</script>

<header>
  {#if show}
    <div
      transition:fly={{ y: -200, duration: 600, opacity: 1 }}
      class="dropdown"
    >
      <ul>
        <li class="hide">
          <NavLink to={"/"} onClick={() => (show = false)}>Inicio</NavLink>
        </li>
        <li class="hide">
          <NavLink to={"/searches"} onClick={() => (show = false)}>Búsquedas Programadas</NavLink>
        </li>
        <li class="hide">
          <NavLink to={"/configuration"} onClick={() => (show = false)}>Configuración</NavLink>
        </li>
      </ul>
    </div>
  {/if}
  <nav>
    <div class="hide">
      <button type="button" on:click|preventDefault={() => (show = !show)}>
        <Icon data={show ? times : bars} scale="2" />
      </button>
    </div>
    <Logo />
    <NavMenu />
  </nav>
</header>

<style>
  header {
    position: fixed;
    width: 100vw;
    height: 3.5rem;
    background-color: white;
    z-index: 100;
    box-shadow: 0 -0.4rem 0.9rem 0.2rem rgba(0, 0, 0, 0.5);
    opacity: 0.95;
    transition: opacity 0.3s;
  }

  header:hover {
    opacity: 1;
  }

  nav {
    height: 100%;
    align-items: center;
    justify-content: space-between;
    display: flex;
    padding: 0 0.5rem;
    background-color: white;
  }
  .hide {
    line-height: 1;
    display: flex;
    align-items: center;
    display: none;
  }
  .dropdown {
    position: absolute;
    top: 100%;
    width: 100%;
    background-color: white;
    box-shadow: 0 -0.4rem 0.9rem 0.2rem rgba(0, 0, 0, 0.5);
    z-index: -1;
    opacity: 0.9;
  }
  button {
    box-shadow: none;
    border: none;
    width: 40px;
  }

  button:active,
  button:hover:not([disabled]) {
    color: inherit;
    box-shadow: none;
  }
  ul {
    list-style: none;
    line-height: 1;
    padding: 0.5rem 0;
  }

  li {
    margin: 0;
    padding: 1rem 2rem;
  }

  @media only screen and (max-width: 700px) {
    header {
      height: 3rem;
    }
    .hide {
      display: block;
    }
  }
</style>
