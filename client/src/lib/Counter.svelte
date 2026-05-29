<script>
  import { store } from "../rs/ServerState.res.mjs"
</script>

{#each $store.fields as {key, name, unit}}
  {@const value = $store.values[key] || "?"}
  {name}: {value} {unit}
  <br>
{:else}
  empty
{/each}
<!--
<pre>{JSON.stringify($store.fields, null, 2)}</pre>
-->
