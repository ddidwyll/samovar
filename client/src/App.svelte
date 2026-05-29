<script>
  import Counter from "./lib/Counter.svelte"
  import { store } from "./rs/ServerState.res.mjs"
</script>

<section id="center">
  ACK: {$store.status.ack}
  <br>
  <Counter />
</section>
