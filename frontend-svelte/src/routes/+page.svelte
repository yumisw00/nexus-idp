<script>
let jobs=[{id:'1',type:'Parse',status:'Active',progress:45},{id:'2',type:'Extract',status:'Queued',progress:0}];
</script>
<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
<div class="bg-slate-900/40 backdrop-blur-lg border border-white/10 shadow-xl p-6 rounded-xl border-dashed hover:bg-slate-800/50 transition-colors flex items-center justify-center min-h-[200px] cursor-pointer"><span class="font-semibold text-lg text-slate-300">Drag & Drop Upload</span></div>
<div class="bg-slate-900/40 backdrop-blur-lg border border-white/10 shadow-xl p-6 rounded-xl flex flex-col justify-center"><h2 class="text-xl font-bold text-slate-200">Token Wallet</h2><p class="text-4xl mt-4 font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-purple-500">12,500</p></div>
<div class="bg-slate-900/40 backdrop-blur-lg border border-white/10 shadow-xl p-6 rounded-xl">
<h2 class="text-xl font-bold mb-4 text-slate-200">Active Jobs</h2>
{#each jobs as job}
<div class="mb-4"><div class="flex justify-between text-sm mb-2 text-slate-300"><span class="font-medium">{job.type} <span class="text-slate-500">#{job.id}</span></span><span>{job.status}</span></div><div class="w-full bg-slate-800 h-2 rounded-full overflow-hidden"><div class="bg-gradient-to-r from-blue-500 to-purple-500 h-2 rounded-full transition-all duration-500" style="width:{job.progress}%"></div></div></div>
{/each}
</div>
</div>
