<script>
  import { createEventDispatcher } from 'svelte';
  import { login } from '../api';

  const dispatch = createEventDispatcher();

  let email = "";
  let password = "";
  let isLoading = false;
  let errorMessage = "";

  async function handleSubmit() {
    isLoading = true;
    errorMessage = "";

    try {
      // Call backend API
      const response = await login(email, password);
      
      // Check if backend returned success
      if (response && response.token) {
        dispatch('loginSuccess');
      } else {
        errorMessage = "Invalid credentials";
      }
    } catch (err) {
      console.error('Login error:', err);
      errorMessage = err.message || "Connection error. Please ensure the backend API is running on port 8080.";
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="max-w-md mx-auto px-6 py-16 w-full">
  <div class="text-center mb-8">
    <h2 class="text-3xl font-light text-gray-900">Welcome back</h2>
    <p class="text-gray-500 mt-2">Sign in to send secure files</p>
  </div>

  {#if errorMessage}
    <div class="mb-6 p-3 rounded-lg bg-red-50 border border-red-100 text-red-600 text-sm text-center">
      {errorMessage}
    </div>
  {/if}

  <form on:submit|preventDefault={handleSubmit} class="space-y-6">
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
      <input 
        type="email" 
        bind:value={email}
        class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-200 focus:border-emerald-500 focus:outline-none" 
        placeholder="you@example.com" 
        required 
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
      <input 
        type="password" 
        bind:value={password}
        class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-200 focus:border-emerald-500 focus:outline-none" 
        placeholder="••••••••" 
        required 
      />
    </div>

    <button 
      type="submit" 
      disabled={isLoading}
      class="w-full py-3 bg-emerald-600 text-white rounded-lg font-medium hover:bg-emerald-700 active:bg-emerald-800 disabled:opacity-70 disabled:cursor-not-allowed transition flex justify-center items-center"
    >
      {#if isLoading}
        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        Signing in...
      {:else}
        Sign In
      {/if}
    </button>

    <p class="text-center text-sm text-gray-500 mt-4">
      Don’t have an account? <a href="#" class="text-emerald-600 hover:underline">Create one</a>
    </p>
  </form>
</div>