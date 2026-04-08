<script>
  import heroImg from './assets/hero.png'
  import Counter from './lib/Counter.svelte'
</script>

<section id="center">
  <div class="hero">
    <img src={heroImg} class="base" width="170" height="179" alt="" />
  </div>
  <div>
    <h1>Get started</h1>
    <p>Edit <code>src/App.svelte</code></p>
  </div>
  <Counter />
</section>
