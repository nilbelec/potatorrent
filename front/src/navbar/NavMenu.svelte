<script>
  import Icon from "svelte-awesome";
  import { github } from "svelte-awesome/icons";
  import NavLink from "./NavLink.svelte";
</script>

<ul>
  <li class="hide">
    <NavLink to={"/"}>Inicio</NavLink>
  </li>
  <li class="hide">
    <NavLink to={"/searches"}>Búsquedas Programadas</NavLink>
  </li>
  <li class="hide">
    <NavLink to={"/configuration"}>Configuración</NavLink>
  </li>
  <li>
    <a
      title="Ir a la página del proyecto"
      href="https://github.com/nilbelec/potatorrent/"
      target="_blank"
    >
      <Icon data={github} scale="2" />
    </a>
  </li>
</ul>

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
  @media only screen and (max-width: 700px) {
    .hide {
      display: none;
    }
  }
</style>
