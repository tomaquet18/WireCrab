<script>
  import SelectDevice from './routes/SelectDevice.svelte'
  import { writable } from 'svelte/store'

  const selectedDevice = writable(null)
  let currentView = 'select'

  function handleNext(event) {
    selectedDevice.set(event.detail.device)
    currentView = 'capture'
  }
</script>

{#if currentView === 'select'}
  <SelectDevice on:next={handleNext} />
{:else if currentView === 'capture'}
  <!-- TODO: CaptureView -->
{/if}