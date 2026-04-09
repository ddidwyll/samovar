<script>
  import { onMount } from "svelte"
  import { SvelteSet } from "svelte/reactivity"

  let connected = $state(false)
  let events = new SvelteSet()

  const addEvent = (name, data) => {
    console.log(name)
    events.add({name, data})
  }

  onMount(() => {
    const es = new EventSource("http://localhost:4000/feed")

    es.onopen = () => connected = true
    es.onerror = () => connected = false

    es.addEventListener("welcome", (e) => {
      addEvent("welcome", e.data)
    })

    es.addEventListener("time", (e) => {
      const data = JSON.parse(e.data)
      addEvent("time", e.data)
    })
  })
</script>

connected: {connected}

<br>

{#each events as {name, data}}
  {name}
  <pre>{data}</pre>
{:else}
  no events
{/each}
