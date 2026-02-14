<script>
  import { createEventDispatcher } from 'svelte';
  import { WrapFile, SelectFile } from '../../wailsjs/go/main/App';
  import { registerAccessRule } from '../api';
  
  const dispatch = createEventDispatcher();

  // Modal state
  let showWrapModal = false;
  let isWrapping = false;
  let wrapMessage = "";
  let selectedFile = "";
  
  // Form data
  let filePassword = "";
  let fileIP = "";
  let fileMAC = "";

  // Activity history
  let recentTransfers = [
    { name: "Project_Specs_v2.pdf", to: "Alice", status: "Sent", date: "2 mins ago" },
    { name: "Backend_Backup.zip", to: "Server", status: "Uploading", date: "Now" },
  ];

  async function openFileSelection() {
    showWrapModal = true;
    wrapMessage = "";
  }

  async function selectFile() {
    try {
      const path = await SelectFile();
      if (path) {
        selectedFile = path;
        wrapMessage = "";
      } else {
        // user cancelled or no file
        selectedFile = "";
      }
    } catch (err) {
      console.error('SelectFile error:', err);
      wrapMessage = 'Could not open file dialog.';
    }
  }

  async function handleWrapFile() {
    if (!selectedFile) {
      wrapMessage = "Please select a file first";
      return;
    }

    if (!filePassword || !fileIP || !fileMAC) {
      wrapMessage = "Please fill in all fields";
      return;
    }

    isWrapping = true;
    wrapMessage = "Processing...";

    try {
      // Call Wails backend to wrap the file
      const result = await WrapFile(selectedFile, filePassword, fileIP, fileMAC);
      
      if (result.startsWith("Success")) {
        wrapMessage = result;
        // Add to history
        const fileName = selectedFile.split(/[\\\/]/).pop();
        recentTransfers = [
          { name: fileName, to: "Wrapped", status: "Sent", date: "now" },
          ...recentTransfers,
        ];
        
        // Close modal after 2 seconds
        setTimeout(() => {
          showWrapModal = false;
          filePassword = "";
          fileIP = "";
          fileMAC = "";
          selectedFile = "";
        }, 2000);
      } else {
        wrapMessage = "Error: " + result;
      }
    } catch (error) {
      console.error('Wrap error:', error);
      wrapMessage = "Error: " + error.message;
    } finally {
      isWrapping = false;
    }
  }

  function handleLogout() {
    dispatch('logout');
  }
</script>

<div class="max-w-4xl mx-auto px-6 py-12 w-full">
  
  <div class="flex justify-between items-end mb-12">
    <div>
      <h2 class="text-3xl font-light text-gray-900">Dashboard</h2>
      <p class="text-gray-500 mt-2">Manage your secure transfers</p>
    </div>
    <button 
      on:click={handleLogout}
      class="text-sm text-gray-500 hover:text-emerald-600 font-medium transition"
    >
      Sign Out
    </button>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
    
    <button 
      on:click={openFileSelection}
      class="bg-white p-8 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition cursor-pointer group text-left"
    >
      <div class="w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center text-emerald-600 mb-6 group-hover:scale-110 transition-transform">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      </div>
      <h3 class="text-xl font-medium text-gray-900">Send New File</h3>
      <p class="text-gray-500 mt-2 text-sm">Wrap and encrypt a file for secure delivery.</p>
    </button>

    <div class="bg-white p-8 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition cursor-pointer group">
      <div class="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center text-gray-600 mb-6 group-hover:scale-110 transition-transform">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
      </div>
      <h3 class="text-xl font-medium text-gray-900">Receive File</h3>
      <p class="text-gray-500 mt-2 text-sm">Download or accept a P2P transfer.</p>
    </div>
  </div>

  <div class="bg-white border border-gray-200 rounded-xl overflow-hidden">
    <div class="px-6 py-4 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
      <h3 class="font-medium text-gray-700">Recent Activity</h3>
      <a href="#" class="text-xs font-medium text-emerald-600 hover:text-emerald-700">View All</a>
    </div>
    <div class="divide-y divide-gray-100">
      {#each recentTransfers as file}
        <div class="px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition">
          <div class="flex items-center space-x-4">
            <div class="w-8 h-8 rounded-full bg-emerald-50 flex items-center justify-center text-emerald-600">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
            </div>
            <div>
              <p class="text-sm font-medium text-gray-900">{file.name}</p>
              <p class="text-xs text-gray-500">To: {file.to}</p>
            </div>
          </div>
          <div class="text-right">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium 
              {file.status === 'Sent' ? 'bg-emerald-100 text-emerald-800' : 'bg-blue-100 text-blue-800'}">
              {file.status}
            </span>
            <p class="text-xs text-gray-400 mt-1">{file.date}</p>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>

<!-- Modal for wrapping files -->
{#if showWrapModal}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-xl p-8 max-w-md w-full mx-4">
      <h3 class="text-2xl font-medium text-gray-900 mb-6">Wrap New File</h3>
      
      {#if wrapMessage}
        <div class="mb-4 p-3 rounded-lg {wrapMessage.startsWith('Success') ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-red-50 text-red-700 border border-red-200'} text-sm">
          {wrapMessage}
        </div>
      {/if}

      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Selected File</label>
          <p class="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-600 text-sm break-all">
            {selectedFile || "No file selected"}
          </p>
          <div class="mt-2 flex gap-2">
            <button
              on:click={selectFile}
              disabled={isWrapping}
              class="px-3 py-2 bg-white border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Choose File
            </button>
            <span class="text-xs text-gray-400 self-center">Use this to pick a file from your system</span>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Access Password</label>
          <input 
            type="password" 
            bind:value={filePassword}
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-200 focus:border-emerald-500 focus:outline-none"
            placeholder="Enter access password"
            disabled={isWrapping}
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">IP Address</label>
          <input 
            type="text" 
            bind:value={fileIP}
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-200 focus:border-emerald-500 focus:outline-none"
            placeholder="e.g., 192.168.1.100"
            disabled={isWrapping}
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">MAC Address</label>
          <input 
            type="text" 
            bind:value={fileMAC}
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-200 focus:border-emerald-500 focus:outline-none"
            placeholder="e.g., 00:1A:2B:3C:4D:5E"
            disabled={isWrapping}
          />
        </div>
      </div>

      <div class="flex gap-3 mt-8">
        <button 
          on:click={() => showWrapModal = false}
          disabled={isWrapping}
          class="flex-1 px-4 py-3 border border-gray-300 rounded-lg text-gray-700 font-medium hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
        >
          Cancel
        </button>
        <button 
          on:click={handleWrapFile}
          disabled={isWrapping}
          class="flex-1 px-4 py-3 bg-emerald-600 text-white rounded-lg font-medium hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed transition flex justify-center items-center"
        >
          {#if isWrapping}
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Wrapping...
          {:else}
            Wrap File
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}