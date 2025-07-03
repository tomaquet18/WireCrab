<script>
  import { createEventDispatcher, onMount } from 'svelte'
  import { GetAvailableDevices } from '../../wailsjs/go/devices/DeviceService'
  import logo from '../assets/images/logo-universal.png'

  const dispatch = createEventDispatcher()
  let devices = []

  onMount(async () => {
    devices = await GetAvailableDevices()
  })

  function selectDevice(dev) {
    if (dev.type !== 'wifi') return
    dispatch('next', { device: dev }) // trigger transition to capture view
  }

  $: sortedDevices = [...devices].sort((a, b) => {
    if (a.type === 'wifi' && b.type !== 'wifi') return -1
    if (a.type !== 'wifi' && b.type === 'wifi') return 1
    return 0
  })
</script>

<main>
  <img class="logo" src="{logo}" alt="WireCrab logo" />
  <div class="title">Welcome to WireCrab</div>
  <div class="subtitle">Select a Wi-Fi interface to begin</div>

  <div class="device-list">
    {#each sortedDevices as dev}
      <div
        class="device-item {dev.type !== 'wifi' ? 'disabled' : ''}"
        on:click={() => selectDevice(dev)}
      >
        <div>
          <strong>{dev.description || dev.name}</strong>
          <div style="font-size: 0.8rem; color: #cbd5e1;">
            {dev.type === 'wifi' ? 'Wireless interface' : dev.type}
          </div>
          {#if dev.type !== 'wifi'}
            <div class="hint">Only wireless interfaces are supported at the moment</div>
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <footer>v0.1.0 — Only Wi-Fi interfaces are available for now</footer>
</main>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;600&display=swap');

  main {
    font-family: 'Poppins', sans-serif;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
    min-height: 100vh;
    background-color: #1b2636;
    color: white;
  }

  .title {
    font-size: 2rem;
    font-weight: 600;
    margin-top: 2rem;
  }

  .subtitle {
    margin-top: 0.5rem;
    margin-bottom: 2rem;
    color: #cbd5e1;
  }

  .device-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    width: 100%;
    max-width: 500px;
  }

  .device-item {
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 0.1rem;
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    background-color: #2c3e50;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .device-item:hover {
    background-color: #3c4f63;
  }

  .device-item.disabled {
    opacity: 0.4;
    cursor: not-allowed;
    background-color: #1f2c3a;
    border: 1px solid #3a475a;
  }

  .device-item .hint {
    font-size: 0.75rem;
    color: #94a3b8;
    margin-top: 0.25rem;
  }

  footer {
    margin-top: auto;
    padding-top: 2rem;
    font-size: 0.875rem;
    color: #94a3b8;
  }

  img.logo {
    width: 120px;
    margin-top: 2rem;
    margin-bottom: 1rem;
  }
</style>