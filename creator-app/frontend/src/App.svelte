<script>
  import Login from './pages/Login.svelte';
  import Dashboard from './pages/Dashboard.svelte';

  // Default to 'login', but we can change this via the toolbar
  let currentView = 'login'; 

  // Handlers for the real app flow
  function handleLoginSuccess() {
    currentView = 'dashboard';
  }

  function handleLogout() {
    currentView = 'login';
  }
</script>

<main class="min-h-screen w-full bg-gray-50 flex flex-col relative">
  
  {#if currentView === 'login'}
    <Login on:loginSuccess={handleLoginSuccess} />
  {:else if currentView === 'dashboard'}
    <Dashboard on:logout={handleLogout} />
  {/if}

  <div class="fixed bottom-4 right-4 bg-gray-800 text-white px-4 py-2 rounded-full shadow-lg opacity-75 hover:opacity-100 transition flex gap-4 text-xs z-50">
    <span class="font-bold text-gray-400 self-center">DEV MODE:</span>
    <button 
      class="hover:text-emerald-400 font-mono {currentView === 'login' ? 'text-emerald-400 underline' : ''}" 
      on:click={() => currentView = 'login'}
    >
      Login View
    </button>
    <div class="w-px bg-gray-600"></div>
    <button 
      class="hover:text-emerald-400 font-mono {currentView === 'dashboard' ? 'text-emerald-400 underline' : ''}" 
      on:click={() => currentView = 'dashboard'}
    >
      Dashboard View
    </button>
  </div>

</main>