<script>
  import Row from "./ServerStateTableRow.svelte"
  import { store } from "../rs/ServerState.res.mjs"
</script>

<table>
  <thead>
    <Row
      name="ACK"
      value={$store.status.ack}
    />
  </thead>
  <tbody>
    {#each $store.fields as field}
      {@const value = $store.values[field.key] || "?"}
      <Row {...field} {value} />
    {/each}
  </tbody>
</table>
