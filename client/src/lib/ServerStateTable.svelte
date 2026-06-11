<script>
  import Row from "./ServerStateTableRow.svelte"
  import { store, patch } from "../rs/ServerState.res.mjs"

  const patchDeviceId = (devId) => {
    patch("device_id", devId, console.error, console.log)
  }
</script>

<table>
  <thead>
    <tr>
      <th colspan="2">
        {#each ($store.values?.devices?.split(";") || []) as devId}
          <button
            onclick={() => patchDeviceId(devId)}
          >
            {devId}
          </button>
        {/each}
      </th>
    </tr>
    <Row
      name="ACK"
      value={$store.status.ack}
    />
  </thead>
  <tbody>
    {#each $store.fields as field}
      {@const value = $store.values[field.key] || "-"}
      <Row {...field} {value} />
    {/each}
  </tbody>
</table>
