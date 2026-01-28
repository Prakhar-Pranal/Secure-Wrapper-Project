<script>
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  // Mock data for history
  let recentTransfers = [
    { name: "Project_Specs_v2.pdf", to: "Alice", status: "Sent", date: "2 mins ago" },
    { name: "Backend_Backup.zip", to: "Server", status: "Uploading", date: "Now" },
  ];
</script>

<div class="max-w-4xl mx-auto px-6 py-12 w-full">
  
  <div class="flex justify-between items-end mb-12">
    <div>
      <h2 class="text-3xl font-light text-gray-900">Dashboard</h2>
      <p class="text-gray-500 mt-2">Manage your secure transfers</p>
    </div>
    <button 
      on:click={() => dispatch('logout')}
      class="text-sm text-gray-500 hover:text-emerald-600 font-medium transition"
    >
      Sign Out
    </button>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
    
    <div class="bg-white p-8 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition cursor-pointer group">
      <div class="w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center text-emerald-600 mb-6 group-hover:scale-110 transition-transform">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
      </div>
      <h3 class="text-xl font-medium text-gray-900">Send New File</h3>
      <p class="text-gray-500 mt-2 text-sm">Wrap and encrypt a file for secure delivery.</p>
    </div>

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